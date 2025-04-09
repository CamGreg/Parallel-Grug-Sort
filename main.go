package main

import (
	"fmt"
	"maps"
	"math/rand"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/jfcg/sorty/v2"
)

func CustomSortInit(input []int, compare func(int, int) int) []int {
	output := make([]int, len(input))
	CustomSort(input, compare, 16, output)
	return output
}

// todo: add final sort algorithm as input parameter
func CustomSort(input []int, compare func(int, int) int, chunks int, sorted []int) {
	inputLength := len(input)
	if inputLength <= 1000 {
		slices.SortStableFunc(input, compare)
		copy(sorted, input)
		// for i, val := range parallelGrugSort(input, compare) {
		// 	sorted[i] = val
		// }
		return
	}
	if chunks > inputLength {
		chunks = inputLength
	}
	// if chunks > inputLength {
	// 	for i, val := range parallelGrugSort(input, compare) {
	// 		sorted[i] = val
	// 	}
	// 	return
	// }

	subSorted := make([][]int, chunks+1)
	subSortedLastIndex := len(subSorted) - 1
	matchingValues := make([][]int, chunks)

	// TODO should find uniques based on passed compare function as map uniqueness is not guaranteed to be the same. For now works for out integer test.
	var uniquePivots = make(map[int]bool)
	for i := range chunks {
		uniquePivots[input[i*inputLength/chunks]] = true
	}

	pivotValues := slices.Collect(maps.Keys(uniquePivots))
	slices.Sort(pivotValues)
	// pivotValues := parallelGrugSort(slices.Collect(maps.Keys(uniquePivots)), compare)
	lastIndex := len(pivotValues) - 1

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done() // Decrement the counter when the goroutine finishes
		for _, val := range input {
			if compare(val, pivotValues[0]) < 0 {
				subSorted[0] = append(subSorted[0], val)
			}
		}
	}()

	if len(uniquePivots) > 1 {
		wg.Add(1)
		go func() {
			defer wg.Done() // Decrement the counter when the goroutine finishes
			for _, val := range input {
				if compare(val, pivotValues[lastIndex]) > 0 {
					subSorted[subSortedLastIndex] = append(subSorted[subSortedLastIndex], val)
				} else if compare(val, pivotValues[lastIndex]) == 0 {
					matchingValues[subSortedLastIndex-1] = append(matchingValues[subSortedLastIndex-1], val)
				}
			}
		}()
	} else if len(uniquePivots) > 2 {
		wg.Add(lastIndex - 1)
		for i := 1; i <= lastIndex; i++ {
			go func() {
				defer wg.Done() // Decrement the counter when the goroutine finishes
				for _, val := range input {
					if compare(val, pivotValues[i]) < 0 && compare(val, pivotValues[i-1]) > 0 {
						subSorted[i] = append(subSorted[i], val)
					} else if compare(val, pivotValues[i-1]) == 0 {
						matchingValues[i-1] = append(matchingValues[i-1], val)
					}
				}
			}()
		}
	}

	wg.Wait()

	// wg.Add(len(subSorted))
	for i, val := range subSorted {
		startIndex := 0
		if i > 0 {
			for j, list := range subSorted[:i] {
				startIndex += len(list)
				if j > 0 {
					startIndex += len(matchingValues[j-1])
				}
			}
		}

		// go func() {
		// 	defer wg.Done()

		if len(matchingValues) < i {
			duplicateIndex := startIndex + len(matchingValues[i])
			for j, list := range matchingValues[i] {
				sorted[duplicateIndex+j] = list
			}
		}

		if len(val) == 0 {
			return
		}
		CustomSort(val, compare, chunks, sorted[startIndex:])
		// }()
	}
	// wg.Wait()
}

func GolangSort(input []int, compareFunc func(int, int) int) []int {
	slices.SortStableFunc(input, compareFunc)
	return input
}

func sortySort(input []int, compareFunc func(int, int) int) []int {
	sorty.SortSlice(input)
	return input
}

func compareInts(a, b int) int {
	return a - b
}

func benchmark(input []int, sortFunc func([]int, func(int, int) int) []int, funcName string) {
	iterations := 100

	outputs := make([]time.Duration, iterations)
	totalDuration := time.Duration(0)
	longestDuration := time.Duration(0)
	shortestDuration := time.Duration(time.Second)

	for i := range iterations {
		inputCopy := make([]int, len(input))
		copy(inputCopy, input)
		debug.SetGCPercent(-1)
		runtime.GC()
		start := time.Now()
		sortFunc(input, compareInts)

		if time.Since(start) > longestDuration {
			longestDuration = time.Since(start)
		}
		if time.Since(start) < shortestDuration {
			shortestDuration = time.Since(start)
		}

		outputs[i] = time.Since(start)
		totalDuration += outputs[i]
		runtime.GC()
		debug.SetGCPercent(100)
	}

	slices.Sort(outputs)

	if len(funcName) < 10 {
		padding := strings.Repeat(" ", 10-len(funcName))
		funcName = funcName + padding
	}

	fmt.Printf("%s Mean: %s, Median %s, Max: %s, Min: %s\n", funcName[:10], totalDuration/time.Duration(iterations), outputs[iterations/2], longestDuration, shortestDuration)
}

func validate() {
	array := make([]int, 1000)
	for i := range len(array) {
		array[i] = rand.Intn(10000)
	}

	Grug := CustomSortInit(array, compareInts)
	slices.Sort(array)

	for i := range Grug {
		if Grug[i] != array[i] {
			fmt.Println("Mismatch")
		}
	}
}

func main() {
	arraySizes := []int{ /*10, 100, 1000, 10000, 100000,*/ 10000000}

	dataDistributions := map[string]func(int) []int{
		"random": func(size int) []int {
			array := make([]int, size)
			for i := range size {
				array[i] = rand.Intn(10000000000)
			}
			return array
		},
		// "sorted": func(size int) []int {
		// 	array := make([]int, size)
		// 	for i := range size {
		// 		array[i] = i
		// 	}
		// 	return array
		// },
		// "reverse_sorted": func(size int) []int {
		// 	array := make([]int, size)
		// 	for i := range size {
		// 		array[i] = size - 1 - i
		// 	}
		// 	return array
		// },
	}

	for _, size := range arraySizes {
		fmt.Printf("\nArray Size: %d\n", size)
		for distributionName, dataGenerator := range dataDistributions {
			sorty.MaxGor = 20
			inputArray := dataGenerator(size)
			fmt.Printf("  Distribution: %s\n", distributionName)
			benchmark(inputArray, CustomSortInit, "Custom")
			benchmark(inputArray, GolangSort, "Golang")
			benchmark(inputArray, sortySort, "sorty")
		}
	}
	validate()
}

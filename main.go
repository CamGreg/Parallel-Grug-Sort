package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func parallelGrugSort(input []int, compare func(int, int) int) []int {
	n := len(input)
	sorted := make([]int, n)
	var wg sync.WaitGroup

	processElement := func(i int) {
		defer wg.Done()
		value := input[i]
		sortedIndex := 0
		offset := 0

		for j := range n {
			comparisonResult := compare(input[j], value)
			if comparisonResult < 0 {
				sortedIndex++
			} else if comparisonResult == 0 && j < i {
				offset++
			}
		}
		sorted[sortedIndex+offset] = value
	}

	wg.Add(n)
	for i := range n {
		go processElement(i)
	}
	wg.Wait()

	return sorted
}

func LimitedParallelGrugSortInit(input []int, compare func(int, int) int) []int {
	output := make([]int, len(input))
	LimitedParallelGrugSort(input, compare, 32, output, 0)
	return output
}

func LimitedParallelGrugSort(input []int, compare func(int, int) int, chunks int, sorted []int, preInputSize int) {
	inputLength := len(input)
	if chunks > inputLength || preInputSize == inputLength {
		for i, val := range parallelGrugSort(input, compare) {
			sorted[i] = val
		}
		return
	}

	subSorted := make([][]int, chunks+1)
	subSortedLastIndex := len(subSorted) - 1

	var uniquePivots = make(map[int]bool)
	for i := range chunks {
		uniquePivots[input[i*inputLength/chunks]] = true
	}

	if len(uniquePivots) <= 2 {
		for i, val := range parallelGrugSort(input, compare) {
			sorted[i] = val
		}
		return
	}

	pivotValues := make([]int, len(uniquePivots))
	i := 0
	for pivot := range uniquePivots {
		pivotValues[i] = pivot
		i++
	}
	lastIndex := len(pivotValues) - 1

	pivotValues = parallelGrugSort(pivotValues, compare)

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done() // Decrement the counter when the goroutine finishes
		for _, val := range input {
			if compare(val, pivotValues[lastIndex]) >= 0 {
				subSorted[subSortedLastIndex] = append(subSorted[subSortedLastIndex], val)
			}
		}
	}()
	go func() {
		defer wg.Done() // Decrement the counter when the goroutine finishes
		for _, val := range input {
			if compare(val, pivotValues[0]) < 0 {
				subSorted[0] = append(subSorted[0], val)
			}
		}
	}()

	for i := 1; i <= lastIndex; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done() // Decrement the counter when the goroutine finishes
			for _, val := range input {
				if compare(val, pivotValues[i]) < 0 && compare(val, pivotValues[i-1]) >= 0 {
					subSorted[i] = append(subSorted[i], val)
				}
			}
		}()
	}

	wg.Wait()

	wg.Add(len(subSorted))

	for i, val := range subSorted {
		startIndex := 0
		if i > 0 {
			for _, list := range subSorted[:i] {
				startIndex += len(list)
			}
		}
		go func() {
			defer wg.Done()
			LimitedParallelGrugSort(val, compare, chunks, sorted[startIndex:], inputLength)
		}()
	}
	wg.Wait()
}

func LimitedparallelMergeSortInit(input []int, compare func(int, int) int) []int {
	output := make([]int, len(input))
	LimitedparallelMergeSort(input, compare, 32, output, 0)
	return output
}

func LimitedparallelMergeSort(input []int, compare func(int, int) int, chunks int, sorted []int, preInputSize int) {
	inputLength := len(input)
	if chunks > inputLength || preInputSize == inputLength {
		for i, val := range parallelMergeSort(input) {
			sorted[i] = val
		}
		return
	}

	subSorted := make([][]int, chunks+1)
	subSortedLastIndex := len(subSorted) - 1

	var uniquePivots = make(map[int]bool)
	for i := range chunks {
		uniquePivots[input[i*inputLength/chunks]] = true
	}

	if len(uniquePivots) <= 2 {
		for i, val := range parallelMergeSort(input) {
			sorted[i] = val
		}
		return
	}

	pivotValues := make([]int, len(uniquePivots))
	i := 0
	for pivot := range uniquePivots {
		pivotValues[i] = pivot
		i++
	}
	lastIndex := len(pivotValues) - 1

	pivotValues = parallelMergeSort(pivotValues)

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done() // Decrement the counter when the goroutine finishes
		for _, val := range input {
			if compare(val, pivotValues[lastIndex]) >= 0 {
				subSorted[subSortedLastIndex] = append(subSorted[subSortedLastIndex], val)
			}
		}
	}()
	go func() {
		defer wg.Done() // Decrement the counter when the goroutine finishes
		for _, val := range input {
			if compare(val, pivotValues[0]) < 0 {
				subSorted[0] = append(subSorted[0], val)
			}
		}
	}()

	for i := 1; i <= lastIndex; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done() // Decrement the counter when the goroutine finishes
			for _, val := range input {
				if compare(val, pivotValues[i]) < 0 && compare(val, pivotValues[i-1]) >= 0 {
					subSorted[i] = append(subSorted[i], val)
				}
			}
		}()
	}

	wg.Wait()

	wg.Add(len(subSorted))

	for i, val := range subSorted {
		startIndex := 0
		if i > 0 {
			for _, list := range subSorted[:i] {
				startIndex += len(list)
			}
		}
		go func() {
			defer wg.Done()
			LimitedParallelGrugSort(val, compare, chunks, sorted[startIndex:], inputLength)
		}()
	}
	wg.Wait()
}

func parallelMergeSort(input []int) []int {
	n := len(input)
	if n <= 1 {
		return input
	}

	mid := n / 2
	var wg sync.WaitGroup
	wg.Add(2)

	var left, right []int
	go func() {
		defer wg.Done()
		left = parallelMergeSort(input[:mid])
	}()
	go func() {
		defer wg.Done()
		right = parallelMergeSort(input[mid:])
	}()
	wg.Wait()

	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, len(left)+len(right))
	i, j, k := 0, 0, 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result[k] = left[i]
			i++
		} else {
			result[k] = right[j]
			j++
		}
		k++
	}

	for i < len(left) {
		result[k] = left[i]
		i++
		k++
	}

	for j < len(right) {
		result[k] = right[j]
		j++
		k++
	}

	return result
}

func parallelQuickSort(input []int) []int {
	if len(input) <= 1 {
		return input
	}

	pivot := input[0]
	left, right := []int{}, []int{}

	for _, v := range input[1:] {
		if v <= pivot {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var leftSorted, rightSorted []int
	go func() {
		defer wg.Done()
		leftSorted = parallelQuickSort(left)
	}()
	go func() {
		defer wg.Done()
		rightSorted = parallelQuickSort(right)
	}()
	wg.Wait()

	result := append(leftSorted, pivot)
	result = append(result, rightSorted...)
	return result
}

func parallelCountingSort(input []int) []int {
	if len(input) == 0 {
		return []int{}
	}

	max := input[0]
	for _, v := range input {
		if v > max {
			max = v
		}
	}

	count := make([]int, max+1)
	for _, v := range input {
		count[v]++
	}

	for i := 1; i < len(count); i++ {
		count[i] += count[i-1]
	}

	sorted := make([]int, len(input))
	for i := len(input) - 1; i >= 0; i-- {
		sorted[count[input[i]]-1] = input[i]
		count[input[i]]--
	}

	return sorted
}

func compareInts(a, b int) int {
	return a - b
}

func benchmark(input []int, sortFunc func([]int, func(int, int) int) []int, funcName string, n int) {
	start := time.Now()
	iterations := 100
	// var output []int
	for range iterations {
		// output = sortFunc(input, compareInts)
		sortFunc(input, compareInts)
	}
	duration := time.Since(start)
	fmt.Printf("%s %s\n", funcName, duration)
	// fmt.Printf("%s: n %d us\n", funcName, int(duration.Microseconds())/n)
	// fmt.Println(output)
}

func validate() {
	array := make([]int, 100000)
	for i := range 100000 {
		array[i] = rand.Intn(100000)
	}

	merge := parallelMergeSort(array)
	gruf := parallelGrugSort(array, compareInts)

	for i := range merge {
		if merge[i] != gruf[i] {
			fmt.Println("Mismatch")
		}
	}
}

func main() {
	arraySizes := []int{10, 100, 1000, 10000, 100000}

	dataDistributions := map[string]func(int) []int{
		"random": func(size int) []int {
			array := make([]int, size)
			for i := range size {
				array[i] = rand.Intn(100)
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
			inputArray := dataGenerator(size)
			fmt.Printf("  Distribution: %s\n", distributionName)
			// benchmark(inputArray, parallelGrugSort, "Grug Sort             ", size)
			benchmark(inputArray, LimitedParallelGrugSortInit, "LimitedParallelGrugSort ", size)
			benchmark(inputArray, LimitedparallelMergeSortInit, "LimitedParallelMergeSort", size)

			benchmark(inputArray, func(input []int, compare func(int, int) int) []int {
				return parallelMergeSort(input)
			}, "Parallel Merge Sort   ", size)
			// benchmark(inputArray, func(input []int, compare func(int, int) int) []int {
			// 	return parallelQuickSort(input)
			// }, "Parallel Quick Sort   ", size)
			benchmark(inputArray, func(input []int, compare func(int, int) int) []int {
				return parallelCountingSort(input)
			}, "Parallel Counting Sort", size)
		}
	}
	validate()
}

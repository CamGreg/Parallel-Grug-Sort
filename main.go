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

		for j := 0; j < n; j++ {
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
	for i := 0; i < n; i++ {
		go processElement(i)
	}
	wg.Wait()

	return sorted
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

func benchmark(input []int, sortFunc func([]int, func(int, int) int) []int, funcName string) {
	start := time.Now()
	iterations := 1000
	for i := 0; i < iterations; i++ {
		sortFunc(input, compareInts)
	}
	duration := time.Since(start)
	fmt.Printf("%s: %s\n", funcName, duration)
}

func main() {
	arraySizes := []int{10, 100, 300, 500, 1000}
	dataDistributions := map[string]func(int) []int{
		"random": func(size int) []int {
			array := make([]int, size)
			for i := 0; i < size; i++ {
				array[i] = rand.Intn(100000)
			}
			return array
		},
		"sorted": func(size int) []int {
			array := make([]int, size)
			for i := 0; i < size; i++ {
				array[i] = i
			}
			return array
		},
		"reverse_sorted": func(size int) []int {
			array := make([]int, size)
			for i := 0; i < size; i++ {
				array[i] = size - 1 - i
			}
			return array
		},
	}

	for _, size := range arraySizes {
		fmt.Printf("\nArray Size: %d\n", size)
		for distributionName, dataGenerator := range dataDistributions {
			inputArray := dataGenerator(size)
			fmt.Printf("  Distribution: %s\n", distributionName)
			benchmark(inputArray, parallelGrugSort, "Grug Sort")
			benchmark(inputArray, func(input []int, compare func(int, int) int) []int {
				return parallelMergeSort(input)
			}, "Parallel Merge Sort")
			benchmark(inputArray, func(input []int, compare func(int, int) int) []int {
				return parallelQuickSort(input)
			}, "Parallel Quick Sort")
			benchmark(inputArray, func(input []int, compare func(int, int) int) []int {
				return parallelCountingSort(input)
			}, "Parallel Counting Sort")
		}
	}

}

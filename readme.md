Toy "benchmark" for sorting algorithms.

# Exploring extreme parrallelism sorting

Grug sort (probabaly the naive sorting implementation) is O(n^2), but can be parallelized infinitely to in practice execute in O(n) time. Obviously quickly hit the limits of parallelism on my CPU, instead of trying to run on a GPU, i created a divide and conquer ingress for it. Somehow that ended up seeming to execute in < O(n) time, but has a lot of overhead compared to golangs sort, so you need to get to ~10M elements before it starts winning.

Sorty's sort was added for a parrallel sort comparison. It behaves similar to Golang's sort, slightly faster in the high n tests but still slower than divided Grug sort. Once the MaxGor parameter is increaded from its default of 3, it's the fastest of the three.

Substituting Golang's sort instead of Grug sort, in the divider function, i didn't see significant differences in performance. However adding a decent cutoff size (~1-10k) and using Golang's sort at that point mitigates the performance differences at lower n while still maintaining the high n performance.

Testing uncontrolled Grug sort vs Golang's sort, i found that Golang's sort was faster, even when able to parallelize every element

## Typical Results

Array Size: 10
  Distribution: random
LimitedParallelGrugSort  19.2µs
GolangSort 5µs
sortySort 10µs
  Distribution: sorted
LimitedParallelGrugSort  8.9µs
GolangSort 4.6µs
sortySort 7.9µs
  Distribution: reverse_sorted
LimitedParallelGrugSort  9.6µs
GolangSort 4.9µs
sortySort 7.7µs

Array Size: 100
  Distribution: random
LimitedParallelGrugSort  82.6µs
GolangSort 62.6µs
sortySort 40.5µs
  Distribution: sorted
LimitedParallelGrugSort  77.2µs
GolangSort 63.9µs
sortySort 32.6µs
  Distribution: reverse_sorted
LimitedParallelGrugSort  208µs
GolangSort 68µs
sortySort 32.7µs

Array Size: 1000
  Distribution: random
LimitedParallelGrugSort  774µs
GolangSort 674.4µs
sortySort 678.299µs
  Distribution: sorted
LimitedParallelGrugSort  641.4µs
GolangSort 541.3µs
sortySort 672µs
  Distribution: reverse_sorted
LimitedParallelGrugSort  824.5µs
GolangSort 697.7µs
sortySort 971.3µs

Array Size: 10000
  Distribution: sorted
LimitedParallelGrugSort  14.801999ms
GolangSort 5.6706ms
sortySort 13.074799ms
  Distribution: reverse_sorted
LimitedParallelGrugSort  16.225099ms
GolangSort 5.7274ms
sortySort 12.928899ms
  Distribution: random
LimitedParallelGrugSort  15.586999ms
GolangSort 7.792499ms
sortySort 12.2757ms

Array Size: 100000
  Distribution: random
LimitedParallelGrugSort  89.127397ms
GolangSort 86.898597ms
sortySort 71.867397ms
  Distribution: sorted
LimitedParallelGrugSort  83.225197ms
GolangSort 55.010798ms
sortySort 74.314598ms
  Distribution: reverse_sorted
LimitedParallelGrugSort  91.813796ms
GolangSort 71.601198ms
sortySort 75.358997ms

Array Size: 10000000
  Distribution: random
LimitedParallelGrugSort  7.552250934s
GolangSort 12.102853273s
sortySort 8.033167177s
  Distribution: sorted
LimitedParallelGrugSort  6.572895944s
GolangSort 7.693557274s
sortySort 7.688783135s
  Distribution: reverse_sorted
LimitedParallelGrugSort  6.51090596s
GolangSort 8.010926036s
sortySort 7.66803099s

Benchmarks on a 2.3 GHz Intel Core i7 running OSX:

```
BenchmarkTokenizeRabin-8            	   20000	     64687 ns/op
BenchmarkTokenizeFnv-8              	   30000	     46575 ns/op
BenchmarkTokenizeSpooky-8           	   20000	     62806 ns/op
BenchmarkConvertToShinglesRabin-8   	  100000	     21595 ns/op
BenchmarkConvertToShinglesFnv-8     	   30000	     43278 ns/op
BenchmarkPermutationFnv-8           	  100000	     19496 ns/op
BenchmarkPermutationLinear-8        	   20000	     65613 ns/op
```

Supposing that the text in this benchmark is representative, we have about 20 microseconds per permutation, which means about 2 milliseconds per document (assuming the document signature is calculated using 100 permutations). This means we can calculate minhashes for one million pages in about 40 minutes.

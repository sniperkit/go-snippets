# FastSS
A derivation of Fast Similarity Search with some differences to keep search time consistent with the average case.
(This is something I built quickly. I didn't test thoroughly)
## Issues
* Does not build in an appropriate timeline
* Should not include the index of a word multiple times for a single match

## TODOs
* Benchmark
* More unit tests (especially for different sizes of d)
* Review and benchmark the insertion time and memory footprint

## Example
```go
fss = NewFss(2)
filename := "/usr/share/dict/words"
f, err := os.Open(filename)
if err != nil {
	t.Error("Unable to open file %s", filename)
}

scanner := bufio.NewScanner(f)
for scanner.Scan() {
	fss.Insert(scanner.Text())
}

results := fss.Search("avast")
fmt.Println(results)
```

# Logger

Another log library based on `go-kit/log`

## Why another log library?

Because I need a simpler and tinier log library

## But we usually use logrus, what is the benefit?

Well, you have most of the features of `Logrus`

- Independent logger object
- Level
- WithFields
- Write to another `*os.File`

## But this is only a wrapper, is it fast?

Yeah this is basically only a wrapper of `go-kit/log`.

You can see the benchmark by yourself:
```
BenchmarkSimpleLogger-4                   200000              5848 ns/op            1424 B/op         36 allocs/op
BenchmarkLoggerWithFields-4               200000              7961 ns/op            2161 B/op         49 allocs/op
BenchmarkLoggerWithLongFields-4           100000             16851 ns/op            3876 B/op         95 allocs/op
BenchmarkErrors-4                         200000              7346 ns/op            1760 B/op         44 allocs/op
BenchmarkLongErrorsFields-4               100000             14248 ns/op            2793 B/op         73 allocs/op
BenchmarkErrorsWithFields-4               200000              8522 ns/op            2296 B/op         52 allocs/op
BenchmarkGokitLog-4                       500000              2307 ns/op             696 B/op         15 allocs/op
PASS
ok      github.com/albert-widi/go_common/logger/go-kit  10.914s
```

The benchmark is quite good, most of the time all type of log should be written below 0.01ms. `time.Now()` is actually killing most of the performance.

You can try to run the benchmark on your machine.

## Errors Type

Errors if faster than using `WithFields`, use `errors.Fields` in `errors` package to add more context to your `error`.

But you can also combone `WithFields` with `errors.Fields`. It is much slower, so it is recommended to just use `errors.Fields`

## Tags

Now error can have tags, but this is experimental and the implementation is kinda rough.

```go
errors.AddTags("tag1", "tag2")
```

Tags will be printed everytime log is called.



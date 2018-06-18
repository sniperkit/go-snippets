if you're doing whitebox testing (code in same package with test file) and found the undefined error even tho the vars/functions are in same package, you may have run the test for single file (ie. `cd my_code && go test my_code_test.go`). try run the test as a _whole_ (ie. `cd my_code && go test .`). 

if you don't want to run the whole file, specify the test with `-run` option. it will only run tests that match the regular expression passed to `-run` option. ie. `go test . -run 2` will check test files inside the directory but only run `hello_test.Test2` because `Test2` is the only test that has `2` in the directory. 

> Running the whole test in same directory won't hurt too because go cache the test result. if the code doesn't change, it will use the cached result.
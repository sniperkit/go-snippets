# newman
It's a test runner for postman.

https://github.com/postmanlabs/newman

you can download it with `docker pull gianarb/newman:3.0.0`

If you like you can use this alias in your laptop, in order to make your
experience with newman container better.
```
 alias newman='docker run --rm -it -v $PWD:/opt -w /opt gianarb/newman:3.0.0'
```

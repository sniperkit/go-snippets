
# storage

This go package wraps storage options used by Caltech Library tooling. It supports standard 
CRUD operations for local file system, AWS S3 and Google Cloud Storage (GS).

If you have your local AWS environment setup then you'll probably only need to set two
environment variables for this to work just like CRUD to the local file system.

## AWS Setup

```shell
    export AWS_SDK_LOAD_CONFIG=1
    export AWS_BUCKET="bucket.example.edu"
```

## Google Cloud Setup

```shell
    export GOOGLE_BUCKET="bucket.example.edu"
```

## Testing the package

By default the test run for local disc only. There are options for testing with
S3 and GS individually and an option -all for running all the tests. For S3 or
GS tests to succeed the the buckets need to exist and you need to setup your
authorization before hand. This can usually be done in a shell script and sourced
into your local environment for your tests.

Tests can be run with the Go test option in the repository directory.

```shell
    go test
    go test -s3
    go test -gs
    go test -all
```

## Package installation

_storage_ package is Go get-able

```
    go get -u github.com/caltechlibrary/storage
```


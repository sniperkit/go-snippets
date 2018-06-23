# apifs: The API file system

apifs enables users to interact with REST APIs (and any HTTP resource) as a
file system.

# Installation

`go get github.com/gophergala/api-fs/cmd/apifs`

# Usage

To start the apifs service:

`apifs -mountpoint=(mountpoint)`

Once the apifs service is running, make a path that corresponds to a valid URL.
For example:

 apifs -mountpoint=/mnt/apifs &
 mkdir -p /mnt/apifs/golang.org/pkg

Each subdirectory in apifs contains one file, `clone`. To initiate a new 
connection, read from `clone`. The resulting ID `n` will be the ID of the next
connection, with the side effect of creating two files, `n.ctl` and `n.body`.

To read the request response, read the `n.body` file in the same directory. If
the request was successful, `n.body` will contain the full body of the
response. Reads from `n.body` will block until `n.ctl` has been closed and the
request has finished.

# Control file

The control file `n.ctl` contains a newline-delimited list of arguments for the
request. Valid forms are as follows:

* method method_type: HTTP request method
* query key [value]: HTTP query parameter with optional value
* header key [value]: HTTP header with optional value

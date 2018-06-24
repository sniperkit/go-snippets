# ScrapeStore

ScrapeStore is a distributed storage system for storing web­scraped data, and analyzing the data for
obtaining Page Ranks of URLs using the stored URL mappings. ScrapeStore offers high reliability and
fault tolerance by using redundancy, much like the philosophy behind popular distributed storage systems
like GFS and HDFS. At the heart of ScrapeStore is a collection of metadata servers, which, together, act
as a replicated metadata store. The actual storage sites are managed in the form of many more slave
nodes, usually much much higher in number than the master nodes. The master/paxos servers synchronize
all their metadata using the Paxos algorithm. The metadata is mainly of two types: 1) mappings of which
key’s data is present on which collection of slaves, and 2) the amount of data stored on each slave node.
The data stored on the slaves is also N­way replicated, where the N is tunable based on the requirements.

This system was built as part of the 15­640 Distributed Systems course at Carnegie Mellon University in
Fall 2015 by Haibin Lin (haibinl) and Abhishek Joshi (abj1). This report delineates the design choices and
implementation details of ScrapeStore.

For more information, see [documentation](https://github.com/eric-haibin-lin/ScrapeStore/blob/master/cmu440-F15/report/application.pdf). 

## System Components

<img src="cmu440-F15/report/archi.png" width="540px" height="300px" />

## Starter Code

The starter code for this project is organized roughly as follows:

```
bin/                               Student-compiled binaries

src/github.com/cmu440-F15/        
  paxos/                           implement Paxos

  tests/                           Source code for official tests
    paxostest/                     Tests your paxos implementation
  
  rpc/
    paxosrpc/                      Paxos RPC helpers/constants
    
tests/                             Shell scripts to run the tests
    paxostest.sh                   Script for checkpoint tests
```

## Instructions

### Compiling your code

To and compile your code, execute one or more of the following commands (the
resulting binaries will be located in the `$GOPATH/bin` directory):

```bash
go install github.com/cmu440/runners/prunner
```

To simply check that your code compiles (i.e. without creating the binaries),
you can use the `go build` subcommand to compile an individual package as shown below:

```bash
# Build/compile the "paxos" package.
go build path/to/paxos

# A different way to build/compile the "paxos" package.
go build github.com/cmu440-F15/paxos
```

##### How to Write Go Code

If at any point you have any trouble with building, installing, or testing your code, the article
titled [How to Write Go Code](http://golang.org/doc/code.html) is a great resource for understanding
how Go workspaces are built and organized. You might also find the documentation for the
[`go` command](http://golang.org/cmd/go/) to be helpful. As always, feel free to post your questions
on Piazza.

##### The `prunner` program

The `prunner` program creates and runs an instance of your
`PaxosNode` implementation. Some example usage is provided below:

```bash
# Start a ring of three paxos nodes, where node 0 has prot 9009, node 1 has port 9010, and so on.
./prunner -ports=9009,9010,9011 -N=3 -id=0 -retries=5
./prunner -ports=9009,9010,9011 -N=3 -id=1 -retries=5
./prunner -ports=9009,9010,9011 -N=3 -id=2 -retries=5
```

You can read further descriptions of these flags by running 
```bash
./prunner -h
```

### Executing the official tests

#### 1. Checkpoint
The tests for checkpoint are provided as bash shell scripts in the `p3/tests` directory.
The scripts may be run from anywhere on your system (assuming your `GOPATH` has been set and
they are being executed on a 64-bit Mac OS X or Linux machine). Simply execute the following:

```bash
$GOPATH/tests/paxostest.sh
```

#### 2. Full test

The tests for the whole project are not provided for this project. You should develop your own
ways of testing your code. After the checkpoint is due, we will open up submissions for the final
version and run the full suite of tests.

If you have other questions about the testing policy please don't hesitate to ask us a question on Piazza!

### Submitting to Autolab

To submit your code to Autolab, create a `paxosapp.tar` file containing your implementation as follows:

```sh
cd $GOPATH/src/github.com/cmu440-F15/
tar -cvf paxosapp.tar paxosapp/
```

In order for us to test your Paxos implementation, you should maintain the same directory structure and
API for RPC calls as in the handout. However, you can add additional packages as part of your application
component.

## Miscellaneous

### Using Go on AFS

For those students who wish to write their Go code on AFS (either in a cluster or remotely), you will
need to set the `GOROOT` environment variable as follows (this is required because Go is installed
in a custom location on AFS machines):

```bash
export GOROOT=/usr/local/depot/go
```

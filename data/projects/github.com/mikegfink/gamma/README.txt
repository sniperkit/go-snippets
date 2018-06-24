To run the code, you need a text file that contains graph data that resembles the Wiki-Vote.txt file in the sampleData directory. Each row must contain a start vertex and an end vertex, and only those two things. Each row will contain a new edge.

To start the system, start the server with:

$GOPATH/bin/server [client connection address] [worker connection address]

Then start a number of workers with:

$GOPATH/bin/workerApp [server connection address] [worker address] [worker id]

And finally start a client to run a job:

$GOPATH/bin/client [server address] [client id] [path to file with graph data] [initial value for PageRank]

An example to fun everything locally in one command is the following:  

alias server='$GOPATH/bin/server 127.0.0.1:8001 127.0.0.1:9000'
alias worker1='$GOPATH/bin/workerApp 127.0.0.1:9000 127.0.0.1:9005 1'
alias worker2='$GOPATH/bin/workerApp 127.0.0.1:9000 127.0.0.1:9010 2'
alias client1='$GOPATH/bin/client 127.0.0.1:8001 1 ./data/minigraph-PageRanks.txt 0'
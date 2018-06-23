# mapreduce
A simple Map Reduce framework for processing data for analytics.

Based on
- Marcio Castilho's [Cheap MapReduce in Go](http://marcio.io/2015/07/cheap-mapreduce-in-go/)

Maintained by
- @sumner-mccarty
- The [RightScale Team](https://www.rightscale.com)


## Setup instructions
Assuming you have a working [Go](https://golang.org) setup:

- `go get github.com/sumner-mccarty/go-libs/mapreduce`

## Testing
- `go test -v`

## Test Coverage
- `go test -coverprofile=c.out                  # output test coverage`
- `sed -i -e "s#.*/\(.*\.go\)#\./\\1#" c.out    # fix paths in c.out`
- `go tool cover -html=c.out                    # open as html`

## Usage
- `import "github.com/sumner-mccarty/go-libs/mapreduce"`
- Create 3 methods: sendInput, mapper, reducer, 
  - sendInput method: makes a chan {} and in a goroutine sends input to it
  - mapper method: takes input and formats it for the reducer and sends it to the reducer channel
  - reducer method: takes input and 'reduces' it, i.e. "count by key", "group by key", etc.
- Decide/Tune the number of workers: `const ( MaxWorkers = 20 )`
- Putting it all together 
  - `inputChan := sendInput()`
  - `result := Process(mapper, reducer, inputChan, MaxWorkers)`
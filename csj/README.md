# csj
Comma-Separated JSON reader and writer
(Allows better streaming than JSON and more flexible than pure CSV)

Based on concepts in 
- Kirit Sælensminde's [Comma Separated JSON](http://www.kirit.com/Comma%20Separated%20JSON)

Maintained by
- @sumner-mccarty
- The [RightScale Team](https://www.rightscale.com)


## Setup instructions
Assuming you have a working [Go](https://golang.org) setup:

- `go get github.com/sumner-mccarty/go-libs/csj`

## Testing
- `go test -v`

## Test Coverage
- `go test -coverprofile=c.out                  # output test coverage`
- `sed -i -e "s#.*/\(.*\.go\)#\./\\1#" c.out    # fix paths in c.out`
- `go tool cover -html=c.out                    # open as html`

## Usage
- `import "github.com/sumner-mccarty/go-libs/csj"`

#### Write CSJ
- Create the CSJ Writer passing in any *io.Writer option
  - `var buffer bytes.Buffer`
  - `csjWriter := NewWriter(&buffer)`
- Make an array of interface{} and push any/different types of objects into it
  - `lineArray := make([]interface{}, 2)`
  - `lineArray[0] = "key1"`
  - `lineArray[1] = myStructArray`
- Write the line
  - `csjWriter.Write(lineArray)`

#### Read CSJ
- Create the CSJ Reader passing in any *io.Reader option
  - `strReader := strings.NewReader(validCsjString)`
  - `csjReader := NewReader(strReader)`
- The Reader json.Unmarshal's each element it extracts via csv.Read into an empty object
  - BUT it needs to know what type of object it is / make an empty instance of the object
  - csj.Reader handles this by accepting an `[]MakeEmptyObject` function pointers
    - `type MakeEmptyObject func() interface{}`
- Setup the MakeEmptyObject methods and array
  - Create methods that return an empty object of the type that was marshalled
    - `func NewString() interface{} { var s string; return &s }`
    - `func NewMyStructArray() interface{} { var myStructArray = make([]MyStruct, 0); return &myStructArray }`
  - Make the []MakeEmptyObject 
    - `rowTypesArray := []MakeEmptyObject{NewString, NewMyStructArray}`
- Read this line (note: the next line could have a different format if needed)
  - `valueArray, err = csjReader.Read(rowTypesArray)`
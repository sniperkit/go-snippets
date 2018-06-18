// Generate an implementation for an interface which can act as a decorator.
// Calls to the decorator object wrap the original interface with new behavior

package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/zpatrick/go-parser"
	"go/format"
	"log"
	"os"
	"strings"
)

// Sample is the set of test cases for this program.  To test:
// > go run decorater.go -type Sample decorater.go > sample_decorator.go
// > go build
type Sample interface {
	Do()
	Maybe() error
	Repeat(s, t, v string)
	Arg1(t string, f func(*int)) (e error)
	Out1([]byte, map[string]func() int) (string, error)
	Complex(**OtherString) (*[]SampleStruct, error)
	Remote(os.File) ([]strings.Reader, *os.File, error)
	Interface(anything interface{}) string
	Range(format string, args ...interface{})
	Prefix(base64.Encoding) error
}

type OtherString string

type SampleStruct struct {
	Val  string
	File *os.File
}

// importList implements a flag Set/String interface
// this allows the flag to be specified multiple times.
type importList []string

func (i *importList) Set(value string) error {
	*i = append(*i, fmt.Sprintf("\"%s\"", strings.Trim(value, "\"")))
	return nil
}

func (i *importList) String() string {
	return strings.Join([]string(*i), ",")
}

var (
	typeName   = flag.String("type", "", "interface name, required")
	importFlag = importList{}
)

func init() {
	flag.Var(&importFlag, "import", "additional path to import (can be repeated)")
}

func Usage() {
	app := os.Args[0]
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", app)
	fmt.Fprintf(os.Stderr, "\t%s -type T [file]\n", app)
	fmt.Fprintf(os.Stderr, "For example:\n")
	fmt.Fprintf(os.Stderr, "\t%s -type Sample decorater.go\n", app)
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetPrefix("deco: ")

	flag.Parse()
	if len(*typeName) == 0 {
		Usage()
		os.Exit(2)
	}

	path := getInputPath()
	log.Printf("Type Name: %s\n", *typeName)
	log.Printf("Searching %s\n", path)
	log.Printf("Additional paths `%s`\n", importFlag.String())

	types := parsePackage(path)
	match := filterTypes(types, *typeName)

	output := writeDecorator(match, types.Imports)
	log.Printf("len: %d", len(output))

	fmt.Println(output)
}

func filterTypes(types *parser.GoFile, name string) *parser.GoInterface {
	for _, goInterface := range types.Interfaces {
		if goInterface.Name == name {
			return goInterface
		}
	}

	log.Fatalf("Interface %s not found", name)
	return nil
}

func getInputPath() string {
	args := flag.Args()
	if len(args) == 0 {
		return ""
	}

	path := args[0]

	return path
}

func parsePackage(name string) *parser.GoFile {
	if !strings.HasSuffix(name, ".go") {
		return nil
	}

	goFile, err := parser.ParseFile(name)
	if err != nil {
		log.Fatal("Failed to parse %s, %v", name, err)
	}

	// printInterfaces(goFile)
	return goFile
}

func printInterfaces(goFile *parser.GoFile) {
	log.Println("-----Interfaces-----")
	for _, goInterface := range goFile.Interfaces {
		log.Println("Interface: " + goInterface.Name)
		for _, goMethod := range goInterface.Methods {
			log.Printf("%s (%d) -> %d\n", goMethod.Name,
				len(goMethod.Params), len(goMethod.Results))
		}
	}
}

func selectTypes(goType *parser.GoType) []string {
	result := []string{}
	if len(goType.Inner) == 0 {
		result = append(result, goType.Type)
	} else {
		for _, goType := range goType.Inner {
			result = append(result, selectTypes(goType)...)
		}
	}

	return result
}

func selectPrefixes(typeNames []string) []string {
	result := []string{}
	for _, name := range typeNames {
		// remove pointer prefix
		name = strings.Trim(name, "*")
		if index := strings.Index(name, "."); index > -1 {
			result = append(result, name[:index])
		}
	}
	return result
}

func findPackages(goInterface *parser.GoInterface) []string {
	types := []string{}
	for _, method := range goInterface.Methods {
		for _, field := range method.Params {
			types = append(types, selectTypes(field)...)
		}
		for _, field := range method.Results {
			types = append(types, selectTypes(field)...)
		}
	}

	prefixes := selectPrefixes(types)

	return prefixes
}

func selectImports(allImports []*parser.GoImport, selectedNames []string) []string {
	output := []string{}
	for _, name := range selectedNames {
		for _, goImport := range allImports {
			if name == goImport.Prefix() {
				formattedPrefix := fmt.Sprintf("%s %s", goImport.Name, goImport.Path)
				output = append(output, formattedPrefix)
			}
		}
	}
	return output
}

func lookupImports(goInterface *parser.GoInterface, allImports []*parser.GoImport) []string {
	imports := []string{}
	imports = append(imports, []string(importFlag)...)

	// loop through all types and find prefixes
	packageNames := findPackages(goInterface)
	// match prefixes to allImports
	selected := selectImports(allImports, packageNames)
	imports = append(imports, selected...)
	return imports
}

func writeDecorator(goInterface *parser.GoInterface, allImports []*parser.GoImport) string {
	b := new(bytes.Buffer)
	// warning
	fmt.Fprintf(b, "// Generated by go-decorator, DO NOT EDIT\n")
	// package
	packageName := goInterface.File.Package
	fmt.Fprintf(b, "package %s\n", packageName)

	// imports
	imports := lookupImports(goInterface, allImports)
	fmt.Fprintf(b, "import (")
	for _, n := range imports {
		fmt.Fprintf(b, "\t%s\n", n)
	}
	fmt.Fprintf(b, ")\n")

	// Struct
	name := goInterface.Name
	fmt.Fprintf(b, structFormat, name)

	// methods
	for _, method := range goInterface.Methods {
		methodBody := writeMethod(method, name)
		fmt.Fprintf(b, methodBody)
	}

	return string(formatSource(b))
}

func writeMethod(method *parser.GoMethod, name string) string {
	var template string
	if methodReturnsError(method.Results) {
		template = methodFormatWithErr
	} else {
		template = methodFormatNoErr
	}

	// give static names to params and results
	modifyNames(method.Params, nameParam)
	modifyNames(method.Results, nameResult)

	methodName := method.Name
	params := formatNameAndType(method.Params)
	returns := formatNameAndType(method.Results)
	passArgs := formatNames(method.Params, true)
	returnArgs := formatNames(method.Results, false)

	// if method does not return anything, don't emit an equals sign
	equals := ""
	if len(method.Results) > 0 {
		equals = "="
	}

	return fmt.Sprintf(template, name, methodName, params, returns, returnArgs, passArgs, equals)
}

func modifyNames(params []*parser.GoType, nameMethod func(int, string, int) string) {
	length := len(params)
	for i, p := range params {
		p.Name = nameMethod(i, p.Type, length)
	}
}

func methodReturnsError(params []*parser.GoType) bool {
	if len(params) == 0 {
		return false
	}
	return params[len(params)-1].Type == "error"
}

func nameParam(i int, typeName string, length int) string {
	return fmt.Sprintf("p%d", i)
}

func nameResult(i int, typeName string, length int) string {
	if i == length-1 && typeName == "error" {
		return "err"
	}
	return fmt.Sprintf("v%d", i)
}

func formatNameAndType(params []*parser.GoType) string {
	names := []string{}
	for _, t := range params {
		names = append(names, fmt.Sprintf("%s %s", t.Name, t.Type))
	}

	return strings.Join(names, ", ")

}

func formatNames(params []*parser.GoType, expandEllipsis bool) string {
	names := []string{}
	for _, t := range params {
		name := t.Name

		if expandEllipsis && strings.HasPrefix(t.Type, "...") {
			name = name + "..."
		}

		names = append(names, name)
	}

	return strings.Join(names, ", ")
}

func formatSource(buf *bytes.Buffer) []byte {
	// use go/format to propertly indent the code and sort imports.
	src, err := format.Source(buf.Bytes())
	if err != nil {
		// Our generated code should not be invalid go, but
		// but this error does happen while developing this project.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return buf.Bytes()
	}
	return src
}

// [1] - interface name `Sample`
const structFormat = `type %[1]sDecorator struct {
	Inner %[1]s
	Decorator func(name string, call func() error) error
}
`

// [1] - interface name `Sample`
// [2] - method name `Call`
// [3] - method params `s1 string, s2 func(int)`
// [4] - method returns `(string, error)`
// [5] - assign temp variables
// [6] - pass method arguments to inner method
// [7] - optional equals sign
const methodFormatWithErr = `func (this * %[1]sDecorator) %[2]s(%[3]s) (%[4]s) {
	call := func() error {
		var err error
		%[5]s %[7]s this.Inner.%[2]s(%[6]s)
		return err
	}
	err = this.Decorator("%[2]s", call)
	return %[5]s
}
`

// [1] - interface name `Sample`
// [2] - method name `Call`
// [3] - method params `s1 string, s2 func(int)`
// [4] - method returns `(string, error)`
// [5] - assign temp variables
// [6] - pass method arguments to inner method
// [7] - optional equals sign
const methodFormatNoErr = `func (this * %[1]sDecorator) %[2]s(%[3]s) (%[4]s) {
	call := func() error {
		%[5]s %[7]s this.Inner.%[2]s(%[6]s)
		return nil
	}
	_ = this.Decorator("%[2]s", call)
	return %[5]s
}
`

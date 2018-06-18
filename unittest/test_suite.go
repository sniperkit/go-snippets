package unittest

import (
	"fmt"
	om "github.com/jacobsa/oglematchers"
	"github.com/mgutz/ansi"
	"os"
	"reflect"
	"testing"
)

// Test suite.
type TestSuite struct {
	t *testing.T
}

// Initialize the test suite.
func (s *TestSuite) Initialize(t *testing.T) {
	s.t = t
}

// Run a test suite.
func RunTestSuite(suite interface{}, t *testing.T) {
	// Determine the base value of the test case.
	sValue := reflect.ValueOf(suite)
	sType := sValue.Type()

	sIndValue := reflect.Indirect(sValue)
	sIndType := sIndValue.Type()

	// Map methods.
	var setUpSuiteFunc, tearDownSuiteFunc, setUpFunc, tearDownFunc, initializeFunc reflect.Value
	testMethods := make([]reflect.Method, 0)

	for i := 0; i < sType.NumMethod(); i++ {
		meth := sType.Method(i)
		fun := meth.Func

		if meth.Name == "SetUpSuite" {
			if !funcTakesSelfReturns0(fun) {
				t.Fatalf("Test suite set-up method must have the following signature: SetUpSuite()")
			}

			setUpSuiteFunc = fun
		} else if meth.Name == "TearDownSuite" {
			if !funcTakesSelfReturns0(fun) {
				t.Fatalf("Test suite tear-down method must have the following signature: TearDownSuite()")
			}

			tearDownSuiteFunc = fun
		} else if meth.Name == "SetUp" {
			if !funcTakesSelfReturns0(fun) {
				t.Fatalf("Test case set-up method must have the following signature: SetUp()")
			}

			setUpFunc = fun
		} else if meth.Name == "TearDown" {
			if !funcTakesSelfReturns0(fun) {
				t.Fatalf("Test case tear-down method must have the following signature: TearDown()")
			}

			tearDownFunc = fun
		} else if meth.Name == "Initialize" {
			initializeFunc = fun
		} else if len(meth.Name) > 4 && meth.Name[:4] == "Test" {
			if !funcTakesSelfReturns0(fun) {
				t.Logf("Ignoring test as it does not match the test method signature: %s", meth.Name)
			}

			testMethods = append(testMethods, meth)
		}
	}

	// Initialize the test case.
	initializeFunc.Call([]reflect.Value{sValue, reflect.ValueOf(t)})

	// Set up the test suite.
	if setUpSuiteFunc.IsValid() {
		setUpSuiteFunc.Call([]reflect.Value{sValue})
	}

	// Write the head of the test suite run.
	fmt.Println(ansi.Color(sIndType.Name(), "blue"))

	// Run each test method.
	for _, testMethod := range testMethods {
		fmt.Printf("    %s...", testMethod.Name)
		os.Stdout.Sync()

		// Set up a hijacker and hijack stdout/stderr.
		hijacker, hErr := newStdHijacker()
		if hErr != nil {
			panic(hErr)
		}

		hijacker.Hijack()
//		hijacker.Release()

		if setUpFunc.IsValid() {
			setUpFunc.Call([]reflect.Value{sValue})
		}

		fun := testMethod.Func
		fun.Call([]reflect.Value{sValue})

		if tearDownFunc.IsValid() {
			tearDownFunc.Call([]reflect.Value{sValue})
		}

		hijacker.Release()

		// Produce the output.
		output := string(hijacker.Bytes())

		formattedResult := ansi.Color("OK", "green")

		if len(output) > 0 {
			fmt.Print("\n")
			fmt.Println(ansi.Color("---------------------------> captured stdout/stderr <--------------------------", "cyan"))
			if output[len(output)-1] != '\n' {
				fmt.Println(output)
			} else {
				fmt.Print(output)
			}
			fmt.Print("\x1b[0m")
			fmt.Println(ansi.Color("---------------------------< captured stdout/stderr >--------------------------", "cyan"))
			fmt.Printf("    ... %s\n", formattedResult)
		} else {
			fmt.Printf(" %s\n", formattedResult)
		}
	}

	// Tear down the test suite.
	if tearDownSuiteFunc.IsValid() {
		tearDownSuiteFunc.Call([]reflect.Value{sValue})
	}
}

func (s *TestSuite) Error(args ...interface{}) {
	s.t.Error(args...)
}

func (s *TestSuite) Errorf(format string, args ...interface{}) {
	s.t.Errorf(format, args...)
}

func (s *TestSuite) Fatal(args ...interface{}) {
	s.t.Fatal(args...)
}

func (s *TestSuite) Fatalf(format string, args ...interface{}) {
	s.t.Fatalf(format, args...)
}

func (s *TestSuite) AssertEqual(expected, actual interface{}) {
	if err := om.Equals(expected).Matches(actual); err != nil {
		s.Fatalf("%v is not equal to %v", actual, expected)
	}
}

func (s *TestSuite) AssertIsNil(actual interface{}) {
	if !isNil(actual) {
		s.Fatalf("%v is not nil", actual)
	}
}

func (s *TestSuite) AssertIsNotNil(actual interface{}) {
	if isNil(actual) {
		s.Fatalf("%v is nil", actual)
	}
}

func (s *TestSuite) Logf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (s *TestSuite) T() *testing.T {
	return s.t
}

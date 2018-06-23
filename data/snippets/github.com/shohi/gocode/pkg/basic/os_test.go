package basic

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var osClear map[string]func() //create a map for storing clear funcs

func init() {
	osClear = make(map[string]func()) //Initialize it
	osClear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	osClear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") // Windows example it is untested, but I think its working
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	osClear["darwin"] = osClear["linux"]
}

func TestHostname(t *testing.T) {
	log.Println(os.Hostname())
}

func TestExit(t *testing.T) {
	log.Println("going to exit")
	os.Exit(-21)
}

func TestDefer(t *testing.T) {
	defer func() {
		log.Println("exit")
	}()
	os.Exit(-1)
}

func TestGetenv(t *testing.T) {
	log.Println(os.Getenv("GOPATH"))
}

func TestOSRemove(t *testing.T) {
	fp := "test/test/test.txt"
	_, err := os.Stat(fp)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(fp), 0777) // for mock, use 0777 for now
	}

	err = ioutil.WriteFile(fp, []byte("Test test"), 0777)
	rootDir := strings.Split(fp, string(filepath.Separator))[0]
	err = os.RemoveAll(rootDir)
	if err != nil {
		log.Println(err)
	}

}

// ref https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go
// not work in test environment
func clearTerminalScreen() {
	value, ok := osClear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                            //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func TestOSInGolang(t *testing.T) {
	log.Println(runtime.GOOS)
	log.Println(runtime.GOARCH)
}

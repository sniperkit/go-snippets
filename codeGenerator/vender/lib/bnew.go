package lib

import (
	"flag"
	"os"
	"os/exec"
)

var URL string = ""

// CreateApp ...
// @desc create project
func CreateApp() bool {
	//Add runtime preconditions

	flag.StringVar(&URL, "gitUrl", "", " git repository url ")

	cmd := exec.Command("go", "get", "github.com/astaxie/beego")
	cmd = exec.Command("go", "get", "github.com/astaxie/beego/orm")
	cmd = exec.Command("go", "get", "github.com/go-sql-driver/mysql")
	cmd = exec.Command("go", "get", "github.com/beego/bee")
	cmd = exec.Command("bee", "new", HeadToUpper(GetDBName()))
	cmd.Start()
	cmd.Wait()
	return true
}

func Src() string {
	return os.Getenv("GOPATH") + "/src/" + URL + "/"
}

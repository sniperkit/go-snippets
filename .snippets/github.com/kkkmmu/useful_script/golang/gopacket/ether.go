package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))

	for {
		buf := make([]byte, 1024)
		numRead, err := f.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("% X\n", buf[:numRead])
	}
}

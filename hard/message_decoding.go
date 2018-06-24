package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func genkey() func() string {
	sl, c := uint(1), -1
	return func() string {
		for {
			if c++; c < (1<<sl)-1 {
				return fmt.Sprintf("%0"+fmt.Sprint(sl)+"b", c)
			}
			sl++
			c = -1
		}
	}
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		if scanner.Text() == "\n" {
			continue
		}
		nextKey, trans := genkey(), make(map[string]string)
		trans[nextKey()] = scanner.Text()
		for scanner.Scan() {
			if scanner.Text() == "\n" {
				continue
			} else if scanner.Text() == "0" || scanner.Text() == "1" {
				break
			}
			trans[nextKey()] = scanner.Text()
		}
		head, curr, l := scanner.Text(), "", 0
		for scanner.Scan() {
			if scanner.Text() == "\n" {
				continue
			}
			if len(head) < 3 {
				head += scanner.Text()
				if head == "000" {
					break
				}
				continue
			} else if l == 0 {
				fmt.Sscanf(head, "%b", &l)
			}
			curr += scanner.Text()
			if len(curr) == l {
				if curr == strings.Repeat("1", l) {
					head, curr, l = "", "", 0
				} else {
					fmt.Print(trans[curr])
					curr = ""
				}
			}
		}
		fmt.Println()
	}
}

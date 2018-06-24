package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	reader := bufio.NewReader(data)
lines:
	for {
		var (
			s        []byte
			isPrefix bool
			curr     struct {
				Menu struct {
					Header string
					Items  []struct {
						ID    int
						Label string
					}
				}
			}
			sumLabels int
		)
		for {
			var sc []byte
			sc, isPrefix, err = reader.ReadLine()
			if err != nil {
				break lines
			}
			s = append(s, sc...)
			if !isPrefix {
				break
			}
		}
		if len(s) < 2 {
			continue
		}
		if err := json.Unmarshal(s, &curr); err != nil {
			print("Warning deserializing: ", fmt.Sprint(err))
		}
		for i := range curr.Menu.Items {
			var c int
			fmt.Sscanf(curr.Menu.Items[i].Label, "Label %d", &c)
			sumLabels += c
		}
		fmt.Println(sumLabels)
	}
}

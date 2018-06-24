package main

import (
	"bufio"
	"fmt"
	"github.com/rsmohamad/comp4321/database"
	"os"
	"strconv"
)

func main() {
	p, _ := database.LoadPrinter("index.db")
	reader := bufio.NewReader(os.Stdin)
	defer p.Close()

	for {
		fmt.Println("Print: 1)Words, 2)Pages, 3)AdjList, 4)PageRank, 5)FwdIndex 6)FwdIndexTitle")
		fmt.Print("Enter option (q to quit): ")
		opt, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		num, _ := strconv.Atoi(string(opt[0]))

		if opt == "q\n" {
			break
		}

		switch num {
		case 1:
			p.PrintAllWords()
			break
		case 2:
			p.PrintAllPages()
			break
		case 3:
			p.PrintAdjList()
			break
		case 4:
			p.PrintPageRank()
			break
		case 5:
			p.PrintForwardIndex(false)
			break
		case 6:
			p.PrintForwardIndex(true)
			break
		default:
			fmt.Println("Invalid option")
			break
		}

	}
}

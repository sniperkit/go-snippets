package main

import (
	_ "database/sql"
	_ "encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
	"github.com/tealeg/xlsx"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

//@TODO  Learn better way to handle value, error returns when you just want to return error.

var logger = log.New(os.Stderr, "DB ", 0)

func main() {
	if strings.Contains(os.Args[1], "://") {
		ViewServerFile(os.Args[1])
	} else {
		ViewLocalFile(os.Args[1])
	}
}

func ViewLocalFile(filepath string) {
	xlFile, err := xlsx.OpenFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	for _, sheet := range xlFile.Sheets {
		table := tablewriter.NewWriter(os.Stdout)
		headered := false
		for _, row := range sheet.Rows {
			vals := []string{}
			for _, cell := range row.Cells {
				value, err := cell.String()
				vals = append(vals, value)
				checkError(err)
			}
			if !headered {
				table.SetHeader(vals)
				headered = true
			} else {
				table.Append(vals)
			}
		}
		table.Render() // Send output
	}

}

func ViewServerFile(urlpath string) {
	tempfilepath, err := downloadFromUrl(urlpath)
	checkError(err)
	ViewLocalFile(tempfilepath)
}

func downloadFromUrl(url string) (string, error) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return "", err
	}
	defer output.Close()

	response, err := http.Get(url)

	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}

	fmt.Println(n, "bytes downloaded.")
	return fileName, nil
}

func checkError(err error) {
	if err != nil {
		logger.Println(err)
	}

}

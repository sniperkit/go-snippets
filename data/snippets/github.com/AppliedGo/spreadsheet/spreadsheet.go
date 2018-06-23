/*
<!--
Copyright (c) 2017 Christoph Berger. Some rights reserved.

Use of the text in this file is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

Use of the code in this file is governed by a BSD 3-clause license that can be found
in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "Processing spreadsheet data in Go"
description = "Read a CSV file, process the data, and write a CSV file, in Go"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2017-12-04"
draft = "false"
domains = ["Automation"]
tags = ["csv", "office", "spreadsheet"]
categories = ["Tutorial"]
+++

Your managers, all through the hierarchy, love circulating spreadsheets via email. (They simply don't know better.) How to extract and analyze the relevant data from the daily mess? Go can help.

<!--more-->

- - -

*This article is also available as a video on the [Applied Go YouTube channel](http://www.youtube.com/c/AppliedGo):*

{{< youtube fGD5rjENGRc >}}

*It is the shortened version of a lecture in my upcoming minicourse ["Workplace Automation With Go"](https://appliedgo.com/p/workplace-automation-with-go).*

- - -

Spreadsheet data is everywhere. You can find it in Excel sheets as well as when downloading business data from a website.

Package [`encoding/csv`](https://golang.org/pkg/encoding/csv/) from the Go standard library can help you processing that data and produce statistics, reports or other kinds of output from it. Here is how.

Let's assume we work at a stationery distributor. Every evening, we receive a spreadsheet containing the orders of the day for review. The data looks like this:

Date       | Order ID | Order Item   | Unit Price | Quantity
-----------|:--------:|--------------|-----------:|---------:
2017-11-17 | 1        | Ball Pen     | 1.99       | 50
2017-11-17 | 2        | Notebook     | 12.99      | 10
2017-11-17 | 3        | Binder       | 4.99       | 25
2017-11-20 | 4        | Pencil       | 0.99       | 100
2017-11-20 | 5        | Sketch Block | 2.99       | 40
2017-11-22 | 6        | Ball Pen     | 1.99       | 30
2017-11-23 | 7        | Sketch Block | 2.99       | 20
2017-11-24 | 8        | Ball Pen     | 1.99       | 60

We're interested in some information that is not directly contained in the data; especially, we want to know

* the total price for each order,
* the total sales volume,
* and the number of ball pens sold.

As we get a new copy every day, creating formulas within the spreadsheet is not an option. Instead, we decide to write a small tool to do the calculations for us. Also, the tool shall add the result to the table and write a new spreadsheet file.

But before starting to code, our first step is to export the spreadsheet data to CSV. To make things a bit more complicated, we export the data with a semicolon as the column separator.

(The exact steps vary, depending on the spreadsheet software used.)

The raw CSV data looks like this:

```
Date;Order ID;Order Item;Unit Price;Quantity
2017-11-17;1;Ball Pen;1.99;50
2017-11-17;2;Notebook;12.99;10
2017-11-17;3;Binder;4.99;25
2017-11-18;4;Pencil;0.99;100
2017-11-18;5;Sketch Block;2.99;40
2017-11-19;6;Ball Pen;1.99;30
2017-11-19;7;Sketch Block;2.99;20
2017-11-19;8;Ball Pen;1.99;60
```

We can see a header row and data rows, with data separated by semicolons.


Now let's dive into the Go code to process the data from this spreadsheet and from all spreadsheets that are still to come.

## Reading and processing CSV data with Go
*/

// ### Imports and main
package main

// We only use packages from the standard library here.
import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// In `main()`, we sketch out our program flow:
//
// * Read the CSV file,
// * calculate the desired numbers, and
// * write the results to a new CSV file.
func main() {
	rows := readOrders("orders.csv")
	rows = calculate(rows)
	writeOrders("ordersReport.csv", rows)
}

/*
### Reading CSV files

As the next step, we need to read in the header row, and then the data rows. The result shall be a two-dimensional slice of strings, or a slice of slices of strings.
*/

// `readOrders` takes a filename and returns a two-dimensional list of spreadsheet cells.
func readOrders(name string) [][]string {

	f, err := os.Open(name)
	// Usually we would return the error to the caller and handle
	// all errors in function `main()`. However, this is just a
	// small command-line tool, and so we use `log.Fatal()`
	// instead, in order to write the error message to the
	// terminal and exit immediately.
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", name, err.Error())
	}

	// After this point, the file has been successfully opened,
	// and we want to ensure that it gets closed when no longer
	// needed, so we add a deferred call to `f.Close()`.
	defer f.Close()

	// To read in the CSV data, we create a new CSV reader that
	// reads from the input file.
	//
	// The CSV reader is aware of the CSV data format. It
	// separates the input stream into rows and columns,
	// and returns a slice of slices of strings.
	r := csv.NewReader(f)

	// We can even adjust the reader to recognize a semicolon,
	// rather than a comma, as the column separator.
	r.Comma = ';'

	// Read the whole file at once. (We don't expect large files.)
	rows, err := r.ReadAll()

	// Again, we check for any error,
	if err != nil {
		log.Fatalln("Cannot read CSV data:", err.Error())
	}

	// and finally we can return the rows.
	return rows
}

/*
### Process the data

Now that the data is read in, we can loop over the rows, and read from or write to each row slice as needed.

This is where we can extract the desired information: The total price for each order, the total sales volume, and the number of ball pens sold.
*/

// `calculate` takes a spreadsheet, extracts and calculates the desired information, and returns the result as a new spreadsheet.
func calculate(rows [][]string) [][]string {

	sum := 0
	nb := 0

	// To process the data, we loop over the rows, and read from
	// or write to each row slice as needed.
	for i := range rows {

		// The first row is the header row. Here, we only want to
		// add a new header for the column that holds the total prices.
		if i == 0 {
			rows[0] = append(rows[0], "Total")
			continue
		}

		// From the next row onwards, we calculate the total
		// price, sum up all prices, and count the number of ball
		// pens being ordered.

		// This is fairly straightforward, as we know the indexes
		// of the item name, the unit price, and the quantity.
		// The only difficulty is that all columns are string
		// values but we need the price and quantity values as
		// numeric values.

		// We know that column 2 contains the item name.
		item := rows[i][2]

		// Another obstacle we are facing here is that the prices are floating-point values but for financial calculations, we want to use precise integer calculation only. Luckily, the [`strings`](https://golang.org/pkg/strings) and [`strconv`](https://golang.org/pkg/strconv/) packages have got us covered.

		// Column 3 contains the price. Remove the decimal point using `strings.Replace()`, and
		// turn the value into an integer (representing the value in cents) using `strconv.Atoi`.
		price, err := strconv.Atoi(strings.Replace(rows[i][3], ".", "", -1))
		if err != nil {
			log.Fatalf("Cannot retrieve price of %s: %s\n", item, err)
		}

		// Column 4 contains the ordered quantity. Again, we convert the value into an integer.
		qty, err := strconv.Atoi(rows[i][4])
		if err != nil {
			log.Fatalf("Cannot retrieve quantity of %s: %s\n", item, err)
		}

		// Calculate the total and append it to the current row.
		total := price * qty

		// We use a helper function to turn the total value (an integer) back into a floating-point value with two decimals, represented as a string (see below).
		rows[i] = append(rows[i], intToFloatString(total))

		// Update the total sum
		sum += total

		// and the # of ball pens.
		if item == "Ball Pen" {
			nb += qty
		}
	}

	// Here we append two new rows. The first one shows the total sum, and
	// the second one shows the number of ball pens ordered.
	rows = append(rows, []string{"", "", "", "Sum", "", intToFloatString(sum)})
	rows = append(rows, []string{"", "", "", "Ball Pens", fmt.Sprint(nb), ""})

	// Return the new spreadsheet.
	return rows
}

// `intToFloatString` takes an integer `n` and calculates the floating point value representing `n/100` as a string.
func intToFloatString(n int) string {
	intgr := n / 100
	frac := n - intgr*100
	return fmt.Sprintf("%d.%d", intgr, frac)
}

/*
### Write the new CSV data

Finally, we write the result to a new file, using `os.Create()` and a CSV writer that knows how to turn the slice of slices of strings back into a proper CSV file.

Note that we do not set the separator to semicolon here, as we  want to create a standard CSV format this time.
*/

// `writeOrders` takes a filename and a spreadsheet and writes the spreadsheet as CSV to the file.
func writeOrders(name string, rows [][]string) {

	f, err := os.Create(name)
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", name, err.Error())
	}

	// We are going to write to a file, so any errors are relevant and
	// need to be logged. Hence the anonymous func instead of a one-liner.
	defer func() {
		e := f.Close()
		if e != nil {
			log.Fatalf("Cannot close '%s': %s\n", name, e.Error())
		}
	}()

	w := csv.NewWriter(f)
	err = w.WriteAll(rows)
}

/*
When running this code, the output file should look like this:

```
Date,Order ID,Order Item,Unit Price,Quantity,Total
2017-11-17,1,Ball Pen,1.99,50,99.50
2017-11-17,2,Notebook,12.99,10,129.90
2017-11-17,3,Binder,4.99,25,124.75
2017-11-18,4,Pencil,0.99,100,99.0
2017-11-18,5,Sketch Block,2.99,40,119.60
2017-11-19,6,Ball Pen,1.99,30,59.70
2017-11-19,7,Sketch Block,2.99,20,59.80
2017-11-19,8,Ball Pen,1.99,60,119.40
,,,Sum,,811.65
,,,Ball Pens,140,
```

And we can open it in our spreadsheet app, or in a CSV viewer, to get a nicely formatted table.

Date       | Order ID | Order Item   | Unit Price | Quantity | **Total**
-----------|:--------:|--------------|-----------:|---------:|-------:
2017-11-17 | 1        | Ball Pen     | 1.99       | 50       | **99.50**
2017-11-17 | 2        | Notebook     | 12.99      | 10       | **129.90**
2017-11-17 | 3        | Binder       | 4.99       | 25       | **124.75**
2017-11-18 | 4        | Pencil       | 0.99       | 100      | **99.0**
2017-11-18 | 5        | Sketch Block | 2.99       | 40       | **119.60**
2017-11-19 | 6        | Ball Pen     | 1.99       | 30       | **59.70**
2017-11-19 | 7        | Sketch Block | 2.99       | 20       | **59.80**
2017-11-19 | 8        | Ball Pen     | 1.99       | 60       | **119.40**
           |          | **Sum**      |            |          | **811.65**
           |          | **Ball Pens**|            | **140**  |

Here we can see our new Totals column, as well as the two new rows that show the overall sum and the number of ball pens ordered.


## How to get and run the code

Step 1: `go get` the code. Note the `-d` flag that prevents auto-installing
the binary into `$GOPATH/bin`.

This time, also note the `/...` postfix that downloads all files, not only those imported by the main package.

    go get -d github.com/appliedgo/spreadsheet/...

Step 2: `cd` to the source code directory.

    cd $GOPATH/src/github.com/appliedgo/spreadsheet

Step 3. Run the binary.

    go run spreadsheet.go

You should then find a file named `ordersReport.csv` in the current directory. Verify that it contains the expected result.


### Q&A: Why CSV?

I use CSV here, rather than the file formats used by Excel or Open/Libre Office or Numbers, in order to stay as flexible and vendor-independent as possible. If you specifically want to work with Excel sheets, a [quick search on GitHub](https://github.com/search?o=desc&q=excel+language%3Ago&s=stars&type=Repositories&utf8=%E2%9C%93) should return a couple of useful third-party libraries. I have not used any of them yet, so I can neither share any experience nor recommend a particular one.

## Links
[Wikipedia: Comma-separated values](https://en.wikipedia.org/wiki/Comma-separated_values) - Details about the CSV format.

**Happy coding!**

*/

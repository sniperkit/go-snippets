/*
Package main (gorearrange.go) :

This is a CLI tool to interactively rearrange a text data on a terminal.

When a text data is rearranged, there are many applications for automatically sorting text data. But there are a little CLI applications for manually rearranging it using. Furthermore, I just had to create an application for manually and interactively rearranging data. So I created this.

# Features of "go-rearrange" are as follows.

1. Data can be interactively rearranged on your terminal as a CLI tool.

2. Output rearranged data as ``[]string``.

3. Retrieve selected values and select history.

https://github.com/tanaikech/gorearrange/

You can read the detail information there.


---------------------------------------------------------------

# Usage

$ cat sample.txt | gorearrange

You can use the standard output ``>`` to output the result as a file. If you use the command prompt on windows dos, please use ``type sample.txt | gorearrange``. For example, if you use msys2, you can use ``winpty gorearrange -i sample.txt``. You can use an option ``-o outputfile`` to output the result as a file.

---------------------------------------------------------------
*/
package main

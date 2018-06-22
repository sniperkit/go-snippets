package dataset

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var (
	skipHTMLRe      = regexp.MustCompile(`^(head|script|style|object)$`)
	htmlEntityRe    = regexp.MustCompile("&((#\\d+)|([a-zA-Z]+));")
	codeBeginByte   = []byte{2}
	codeEndByte     = []byte{3}
	blockBorderByte = []byte{4}
)

func parseHTMLEntity(entName []byte) []byte {
	entNameStr := strings.ToLower(string(entName[1 : len(entName)-1]))

	if entNameStr[0] == '#' {
		val, err := strconv.Atoi(entNameStr[1:])
		if err != nil {
			return entName
		}
		return []byte(string(rune(val)))
	}
	// the list is not full
	switch entNameStr {
	case "nbsp":
		return []byte(" ")
	case "gt":
		return []byte(">")
	case "lt":
		return []byte("<")
	case "amp":
		return []byte("&")
	case "quot":
		return []byte("\"")
	case "apos":
		return []byte("'")
	case "cent":
		return []byte("¢")
	case "pound":
		return []byte("£")
	case "yen":
		return []byte("¥")
	case "euro":
		return []byte("€")
	case "copy":
		return []byte("©")
	case "reg":
		return []byte("®")
	case "ldquo":
		return []byte("\"")
	case "rdquo":
		return []byte("\"")
	case "lsquo":
		return []byte("'")
	case "rsquo":
		return []byte("'")
	case "sbquo":
		return []byte("\"")
	case "rbquo":
		return []byte("\"")
	case "bdquo":
		return []byte("\"")
	case "ndash":
		return []byte("-")
	case "mdash":
		return []byte("-")
	case "bull":
		return []byte("*")
	case "hellip":
		return []byte("...")
	case "prime":
		return []byte("'")
	case "lsaquo":
		return []byte("'")
	case "rsaquo":
		return []byte("'")
	case "trade":
		return []byte("™")
	case "minus":
		return []byte("-")
	case "raquo":
		return []byte("\"")
	case "laquo":
		return []byte("\"")
	case "deg":
		return []byte("°")
	case "sect":
		return []byte("*")
	case "iexcl":
		return []byte("¡")
	default:
		return entName
	}
}

// convertHTMLToText parses HTML as plain text. E.g. it rips all the tags.
// It marks <code> and </code> positions as codeBeginByte and codeEndByte
func convertHTMLToText(htmlSource string, output io.Writer) {
	doc := html.NewTokenizer(strings.NewReader(htmlSource))
	skip := false
	for token := doc.Next(); token != html.ErrorToken; token = doc.Next() {
		tagName, _ := doc.TagName()
		if skipHTMLRe.Match(tagName) {
			if doc.Token().Type != html.SelfClosingTagToken {
				skip = !skip
			}
			continue
		}
		if skip {
			continue
		}
		text := doc.Text()
		text = htmlEntityRe.ReplaceAllFunc(text, parseHTMLEntity)
		text = bytes.Replace(text, []byte("\u00a0"), []byte(" "), -1)
		output.Write(text)
		strTagName := string(tagName)
		if strTagName == "br" {
			output.Write([]byte{'\n'})
		} else if strTagName == "code" {
			if doc.Token().Type == html.StartTagToken {
				output.Write(codeBeginByte)
			} else {
				output.Write(codeEndByte)
			}
		} else if strTagName == "a" {
			for key, val, _ := doc.TagAttr(); key != nil; key, val, _ = doc.TagAttr() {
				if string(key) == "href" {
					output.Write(val)
					output.Write([]byte{' '})
					break
				}
			}
		}
	}
}

func main() {
	csvReader := csv.NewReader(os.Stdin)
	i := 0
	for record, err := csvReader.Read(); err != io.EOF; record, err = csvReader.Read() {
		i++
		if i == 1 {
			// skip the header
			continue
		}
		if i > 2 {
			os.Stdout.Write(blockBorderByte)
		}
		convertHTMLToText(record[len(record)-1], os.Stdout)
	}
}

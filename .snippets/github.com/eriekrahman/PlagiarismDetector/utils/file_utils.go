package utils

import (
	"log"
	"io/ioutil"
	"strings"
)

/*****************************************************************************************************
 * Function: LoadFile
 * Purpose: To load the content of a file
 * Parameter: filePath as the complete path of a file
 * Output: content of a file
 *****************************************************************************************************/
func LoadFile(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return ""
	} else {
		return string(content)
	}
}

/*****************************************************************************************************
 * Function: RemoveUnusedPunctuation
 * Purpose: To remove unused punctuation character from the text
 * Parameter: text as a file's content (comprises 1 or more than 1 sentence)
 * Output: truncated text
 *****************************************************************************************************/
func RemoveUnusedPunctuation(text string) string {
	strings.Replace(text, ",", "", -1)
	strings.Replace(text, ";", "", -1)
	strings.Replace(text, ":", "", -1)
	strings.Replace(text, "&", "", -1)
	strings.Replace(text, "?", "", -1)
	strings.Replace(text, "!", "", -1)
	strings.Replace(text, "\"", "", -1)
	strings.Replace(text, "'", "", -1)	
	return text
}

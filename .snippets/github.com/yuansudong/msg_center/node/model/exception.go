package model

import (
	"fmt"
)

// UnknownProtoCode is used as exception code from remote computer
func UnknownProtoCode(requestContent []byte) (data []byte) {
	fmt.Printf("have receive unknown prtocol code, protocol content is %s\n", string(requestContent))
	return data
}

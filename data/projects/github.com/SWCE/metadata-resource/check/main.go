package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/swce/metadata-resource/models"
)

func main() {
	var request models.CheckRequest
	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err.Error())
		os.Exit(1)
	}
	t := strconv.FormatInt(time.Now().UnixNano(),10)
	versions := models.CheckResponse {
		models.TimestampVersion {
			Version: t,
		},
	}
	json.NewEncoder(os.Stdout).Encode(versions)
}

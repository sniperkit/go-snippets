package utils

import (
	"os"
)

// Options terminal args
type Options struct {
	News        bool
	Github      bool
	Films       bool
	V2EX        bool
	Zhihu       bool
	JianDan     bool
	HackNews    bool
	Reddit      bool
	ProductHunt bool
	Medium      bool
	Query       string
}

// ParseArgs parse terminal arguments
func ParseArgs(args []string) *Options {
	if len(os.Args) == 1 {
		return nil
	}
	options := &Options{}

	switch args[1] {
	case "news":
		options.News = true
		break
	case "github":
		options.Github = true
		break
	case "films":
		options.Films = true
		break
	case "zhihu":
		options.Zhihu = true
		break
	case "jiandan":
		options.JianDan = true
		break
	case "reddit":
		options.Reddit = true
		break
	case "medium":
		options.Medium = true
		break
	case "v2ex":
		options.V2EX = true
		break
	case "hn":
		options.HackNews = true
		break
	case "ph":
		options.ProductHunt = true
		break
	case "help":
		return nil
	default:
		break
	}
	if len(os.Args) == 3 {
		options.Query = args[2]
	}
	return options
}


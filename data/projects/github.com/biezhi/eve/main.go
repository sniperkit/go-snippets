package main

import (
	"fmt"
	"os"

	"github.com/biezhi/eve/utils"
)

var (
	logo = `
______   __   __ ______    
/\  ___\ /\ \ / //\  ___\   
\ \  __\ \ \ \'/ \ \  __\   
 \ \_____\\ \__|  \ \_____\ 
  \/_____/ \/_/    \/_____/ 
  
eve v0.0.4

everyday explore.

Inspired by https://github.com/biezhi/eve
`
)

func green(s string) {
	fmt.Println(utils.Green(s))
}

func displayUsage() {
	fmt.Println(utils.Red(logo))

	green("Usage:\n")
	green("    eve <option>\n")
	green("    eve news    (show all news)")
	green("    eve github  (show github trending or search repo)")
	green("    eve zhihu   (show zhihu hot topics)")
	green("    eve jiandan (show jiandan hot topics)")
	green("    eve reddit  (show reddit hot topics)")
	green("    eve medium  (show medium hot topics)")
	green("    eve v2ex    (show v2ex hot topics)")
	green("    eve hn      (show HackNews hot topics)")
	green("    eve ph      (show Product Hunt hot product)")
	green("    eve help    (show usage)")
}

func main() {
	options := utils.ParseArgs(os.Args)
	if options == nil {
		displayUsage()
		os.Exit(0)
	}
	if options.News {
		ShowNews("")
	}
	if options.Github {
		//ShowGithubInfo(options.Query)
		ShowNews("github")
	}
	if options.Films {
		ShowFilms(options.Query)
	}
	if options.Zhihu {
		ShowNews("zhihu")
	}
	if options.JianDan {
		ShowNews("jiandan")
	}
	if options.HackNews {
		ShowNews("hn")
	}
	if options.Reddit {
		ShowNews("reddit")
	}
	if options.Medium {
		ShowNews("medium")
	}
	if options.ProductHunt {
		ShowNews("ph")
	}
	if options.V2EX {
		ShowNews("v2ex")
	}
}

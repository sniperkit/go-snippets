package main

import (
	"fmt"

	"github.com/gophergala/cobs/hunter"
)

func main() {
	ghc := hunter.ParseGitHubURL("https://github.com/projectchrono/chrono/blob/dbff7a6564285b8494dfe481036b78cbdea061ae/.clang-format")
	fmt.Println(ghc)
	ghc = hunter.ParseGitHubURL("https://raw.githubusercontent.com/projectchrono/chrono/dbff7a6564285b8494dfe481036b78cbdea061ae/.clang-format")
	fmt.Println(ghc)
	ghc = hunter.ParseGitHubURL("https://raw.githubusercontent.com/projectchrono/chrono/dbff7a6564285b8494dfe481036b78cbdea061ae/.clang-format/moo")
	fmt.Println(ghc)
	ghc = hunter.ParseGitHubURL("https://raw.githubusercontent.com/projectchrono/chrono")
	fmt.Println(ghc)

}

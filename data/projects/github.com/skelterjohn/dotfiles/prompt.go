package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func color(code, text string) string {
	return fmt.Sprintf(`\[\e[%s\]%s\[\e[0m\]`, code, text)
}

func yellow(text string) string {
	return color("1;33m", text)
}

func green(text string) string {
	return color("1;32m", text)
}

func red(text string) string {
	return color("1;31m", text)
}

var (
	CITC_ROOT = "/google/src/cloud/"+os.Getenv("USER")
)

func printCompressedAbsPath(abspath string) {
	dir := "/"

	colorfulTilde := green("~")

	home := os.Getenv("HOME")
	if strings.HasPrefix(abspath, home) {
		abspath = colorfulTilde + abspath[len(home):]
		dir = home
	}

	tokens := strings.Split(abspath, "/")

	for _, t := range tokens[:len(tokens)-1] {
		if len(t) == 0 {
			// since we lead with a slash, an empty first element indicates root
			continue
		}
		end := 5
		if end > len(t) {
			end = len(t)
		}
		shortestUnique := t[:end]

		f, err := os.Open(dir)
		if err == nil {
			if names, err := f.Readdirnames(-1); err == nil {
				for _, name := range names {
					if name == t {
						continue
					}
					for shortestUnique != t && strings.HasPrefix(name, shortestUnique) {
						shortestUnique = t[:len(shortestUnique)+1]
					}
				}
			}
		}
		if t != colorfulTilde {
			dir = path.Join(dir, t)
			if shortestUnique != t {
				shortestUnique += yellow("*")
			}
		} else {
			shortestUnique = colorfulTilde
		}

		fmt.Printf("%s/", shortestUnique)
	}
	fmt.Printf("%s", tokens[len(tokens)-1])
}

func gitRepoAndBranch(path string) (repoRoot, repoPath, branch string, ok bool) {
	cmd := exec.Command("sh", "-c", `git branch | grep '^\*'`)
	cmd.Dir = path
	out, err := cmd.Output()
	if err == nil && len(out) >= 3 {
		branch = string(out[2:])
	}

	cmd = exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = path
	out, err = cmd.Output()
	if err != nil || len(out) < 1 {
		return
	}
	repoRoot = strings.TrimSpace(string(out))
	repoPath, err = filepath.Rel(repoRoot, path)
	repoPath = strings.TrimSpace(repoPath)
	if err != nil {
		return
	}
	branch = strings.TrimSpace(branch)
	ok = true
	return
}

func main() {
	pwd, _ := os.Getwd()

	fmt.Print("[\\t]")

	if root, path, branch, ok := gitRepoAndBranch(pwd); ok {
		printCompressedAbsPath(root)
		fmt.Printf("%s%s%s/", red("{"), yellow(branch), red("}"))
		if path != "." {
			fmt.Print(path)
		}
		fmt.Print("$ ")
		return
	}

	if strings.HasPrefix(pwd, CITC_ROOT) {
		subpath := pwd[len(CITC_ROOT):]
		var tokens []string
		rest := subpath
		var client string
		for {
			newrest, token := filepath.Split(rest)
			newrest = strings.Trim(newrest, "/")
			if newrest == "" {
				client = token
				break
			}
			if newrest == rest {
				client = newrest
				break
			}
			rest = newrest
			tokens = append([]string{token}, tokens...)
		}
		subpath = filepath.Join(tokens...)
		fmt.Printf(`%s:%s/%s`, green("citc"), yellow(client), subpath)
		fmt.Print("$ ")
		return
	}

	printCompressedAbsPath(pwd)
	fmt.Print("$ ")
}

package main

import (
	//"bufio"
	"fmt"
	//"os"
	"regexp"
	//"strings"
	"telnet"
)

var ResultR = regexp.MustCompile("(?P<status>[[:alnum:]-]+)[[:space:]]+{(?P<result>[[:alnum:][:space:]_-]+)}")

var Session *telnet.Session

func main() {
	res, err := Invoke("AgtHelp ListObjects")
	if err != nil {
		panic(err)
	}

	fmt.Println(res)

	res, err = Invoke("AgtHelp ListMethods AgtSessionManager")
	if err != nil {
		panic(err)
	}

	fmt.Println(res)

	res, err = Invoke("AgtHe")
	if err != nil {
		fmt.Println(err.Error())
	}

	res, err = Run("AgtHe")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res)
}

func Run(method ...string) (string, error) {
	if len(method) == 0 {
		return "", fmt.Errorf("You must give the method to run")
	}

	var cmd string
	for _, p := range method {
		cmd += fmt.Sprintf(" %s", p)
	}
	_, err := Session.WriteLine(cmd)
	if err != nil {
		return "", fmt.Errorf("Run %s with error: %s", cmd, err.Error())
	}

	return GetCommandResult()
}

func GetCommandResult() (string, error) {
	var line []byte

	b, err := Session.ReadByte()
	if err != nil {
		return "", fmt.Errorf("Cannot get command result: ", err.Error())
	}
	if b == 'i' {
		line, err = Session.ReadUntilSkip([]string{"\""}, []string{"name"})
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, line))
	} else {
		sline, err := Session.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, sline))
	}

	return string(line), nil
}

func Invoke(method ...string) (string, error) {
	cmd := fmt.Sprintf("%s ", "invoke")
	for _, p := range method {
		cmd += fmt.Sprintf(" %s", p)
	}
	_, err := Session.WriteLine(cmd)
	if err != nil {
		return "", fmt.Errorf("Run %s with error: %s", cmd, err.Error())
	}

	line, err := GetCommandResult()
	if err != nil {
		return "", fmt.Errorf("Cannot get result of: %s with error: %s", cmd, err.Error())
	}
	fmt.Println(line)

	res := ResultR.FindStringSubmatch(line)
	if len(res) != 3 {
		return "", fmt.Errorf("Run %s with invalid result: %s", cmd, line)
	}

	if res[1] == "0" {
		return res[2], nil
	}

	return "", fmt.Errorf("Run %s with result: (%s: %s)", cmd, res[1], res[2])
}

func init() {
	sess, err := telnet.New3("10.71.20.231:9001")
	if err != nil {
		panic(err)
	}

	Session = sess

	/*
		go func(r *telnet.Session) {
			for {
				line, err := r.ReadString('\n')
				if err != nil {
					panic(err)
				}
				fmt.Printf("\r> %s", line)
				fmt.Printf("\r>")
			}
		}(sess)

		go func(r *telnet.Session) {
			fmt.Println(">")
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if scanner.Text() == "quit" {
					os.Exit(0)
				}

				line := scanner.Text()
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "run ") {
					line = strings.TrimPrefix(line, "run ")
					if strings.TrimSpace(line) != "" {
						_, err := r.WriteLine(line)
						if err != nil {
							panic(err)
						}
					}
				}
				fmt.Printf(">")
			}
		}(sess)
	*/
}

package publikey

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/codegangsta/cli"
	"github.com/gerred/publikey/server"
	"github.com/mitchellh/go-homedir"
)

func NewRegisterCommand() cli.Command {
	return cli.Command{
		Name:      "register",
		ShortName: "r",
		Usage:     "create a new publikey user",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "username, u",
				Usage: "publikey username or email (required)",
			},
			cli.StringFlag{
				Name:  "password, p",
				Usage: "publikey password (required)",
			},
		},
		Action: handleRegisterCommand,
	}
}

func handleRegisterCommand(c *cli.Context) {
	username := c.String("username")
	password := c.String("password")
	var user server.NewUser
	if username == "" || password == "" {
		fmt.Println("[err] username and password required to register a publikey account")
	}
	resp, _ := http.PostForm("http://"+c.GlobalString("host")+"/users/new", url.Values{"email": {username}, "password": {password}})

	if resp.StatusCode != 200 {
		fmt.Println("error occurred")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &user)

	if err != nil {
		return
	}

	apiKey := user.ApiKey
	dir, _ := homedir.Dir()
	rcPath := path.Join(dir, ".publikeyrc")

	ioutil.WriteFile(rcPath, []byte(apiKey), 0644)

	fmt.Println("Logged in as " + username)
}

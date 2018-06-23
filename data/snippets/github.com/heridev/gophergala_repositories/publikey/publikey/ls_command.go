package publikey

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/codegangsta/cli"
	"github.com/gerred/publikey/server"
)

func NewListCommand() cli.Command {
	return cli.Command{
		Name:   "ls",
		Usage:  "list keys for a user, defaults to logged in user",
		Action: handleListCommand,
	}
}

func handleListCommand(c *cli.Context) {
	resp, _ := http.Get("http://" + c.GlobalString("host") + "/users/" + c.GlobalString("user") + "/keys")
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var user server.User

	json.Unmarshal(body, &user)

	fmt.Printf("Keys for %s:\n", c.GlobalString("user"))
	for _, key := range user.Keys {
		fmt.Println(key.Value)
	}
}

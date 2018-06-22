package publikey

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
)

func NewAddCommand() cli.Command {
	return cli.Command{
		Name:      "add",
		ShortName: "a",
		Usage:     "add key to publikey under the logged in user",
		Action:    handleAddCommand,
	}
}

func handleAddCommand(c *cli.Context) {
	homeDir, _ := homedir.Dir()
	rcPath := path.Join(homeDir, ".publikeyrc")
	dat, _ := ioutil.ReadFile(rcPath)
	apiKey := string(dat)
	fileName := c.Args()[0]
	dir, _ := filepath.Abs(filepath.Dir(fileName))
	sshKey, _ := ioutil.ReadFile(path.Join(dir, fileName))

	client := &http.Client{}
	data := url.Values{"key": {string(sshKey)}}
	req, _ := http.NewRequest("POST", "http://"+c.GlobalString("host")+"/keys", strings.NewReader(data.Encode()))
	req.Header.Add("X-PUBLIKEY-API-KEY", apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client.Do(req)

}

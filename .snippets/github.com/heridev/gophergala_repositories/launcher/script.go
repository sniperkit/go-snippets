package main

import (
	"code.google.com/p/go.crypto/ssh"
	"io"
	"io/ioutil"
	"os/user"
	"strconv"
)

type Host struct {
	Name     string
	User     string
	Password string
	Port     int
}

type Script struct {
	Name    string
	Host    string
	Content string
}

func (self *Script) Execute(host *Host, out io.Writer) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	if host.User == "" {
		host.User = usr.Username
	}
	cfg := &ssh.ClientConfig{
		User: host.User,
	}
	if host.Password != "" {
		cfg.Auth = []ssh.AuthMethod{
			ssh.Password(host.Password),
		}
	} else {
		content, err := ioutil.ReadFile(usr.HomeDir + "/.ssh/id_rsa")
		if err != nil {
			content, err = ioutil.ReadFile(usr.HomeDir + "/.ssh/id_dsa")
			if err != nil {
				return err
			}
		}
		key, err := ssh.ParsePrivateKey(content)
		if err != nil {
			return err
		}
		cfg.Auth = []ssh.AuthMethod{ssh.PublicKeys(key)}
	}
	client, err := ssh.Dial("tcp", host.Name+":"+strconv.Itoa(host.Port), cfg)
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stdout = out
	session.Stderr = out
	if err := session.Run(self.Content); err != nil {
		return err
	}
	return nil
}

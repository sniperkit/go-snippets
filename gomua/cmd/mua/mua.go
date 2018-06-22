package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gophergala/gomua"
)

// client handles common data as a user navigates the MUA.
// messages is the list of all mail (in the current folder?)
// current is the currently selected Mail
// displayN is the # of Mail to display on the screen at one time.
// user is the user's email address, for sending.
type client struct {
	messages []gomua.Mail
	current  gomua.Mail
	displayN int
	user     string
	dir      string
}

// reads from the config file, creates a new client
func NewClient(filename string) (*client, error) {
	c := new(client)

	// TODO: much of the following is duplicated from send.go, refactor
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("Client: missing " + gomua.ConfigLocation + " file.")
	}

	var client string
	sections := strings.Split(string(b), "[")
	for _, sec := range sections {
		if strings.HasPrefix(sec, "client]") {
			client = "[" + sec
		}
	}

	lines := strings.Split(client, "\n")
	for _, l := range lines {
		switch {
		case strings.HasPrefix(l, "Maildir="):
			c.dir = strings.TrimPrefix(l, "Maildir=")
		case strings.HasPrefix(l, "DisplayN="):
			c.displayN, _ = strconv.Atoi(strings.TrimPrefix(l, "DisplayN="))
		case strings.HasPrefix(l, "User="):
			c.user = strings.TrimPrefix(l, "User=")
		}
	}

	if c.dir == "" || c.user == "" || c.displayN == 0 {
		return nil, errors.New("Client: incorrect " + gomua.ConfigLocation + " file.")
	}

	return c, nil
}

// takes a Maildir directory, scans for messages, and returns a slice of Message structs.
func (c *client) scanMailDir(dir string) {
	var msgs []gomua.Mail

	// scan new and cur folder
	newmail := gomua.Scan(dir + "new/")
	curmail := gomua.Scan(dir + "cur/")

	for _, m := range newmail {
		folder, name := filepath.Split(m.Filename)          // grab the base name and the path leading to it
		root := strings.TrimRight(folder, "/new/")          // slice off the expected "/new/" from the path to get the root Maildir
		newname := filepath.Join(root, "cur", name) + ":2," // now add /cur/ to the root and the base name and the processed flag
		err := os.Rename(m.Filename, newname)               // finally move the file to the new name as calculated above
		m.Filename = newname
		if err != nil {
			log.Fatal(err)
		}
		// TODO: I have *absolutely* no idea why the following line isn't necessary...
		//msgs = append(msgs, m)
	}

	for _, cm := range curmail {
		msgs = append(msgs, cm)
	}

	c.messages = msgs
}

// takes a slice of Messages and prints a numbered list of summaries
func viewMailList(msgs []gomua.Mail, start int, w io.Writer) {
	var unread string
	for i, msg := range msgs {
		switch m := msg.(type) {
		case *gomua.Message:
			if m.Unread() {
				unread = color("(Unread) ", "34")
			} else {
				unread = ""
			}
		}
		fmt.Fprintf(w, "%d. %s%s\n", i+start+1, unread, msg.Summary())
	}
}

// prints a single mail message to the screen
func viewMail(msg gomua.Mail, w io.Writer) {
	fmt.Fprint(w, msg)

	switch m := msg.(type) {
	case *gomua.Message:
		m.Flag("S")
	}
}

// prompts the user for the response content, and sends a reply to the mail
func replyMessage(old *gomua.Message, user string) (reply *mail.Message) {
	oldid := old.Header.Get("Message-ID")
	oldref := old.Header.Get("References")

	to := "To: " + old.Header.Get("From") + "\r\n"
	from := fmt.Sprintf("From: %s\r\n", user)

	subject := fmt.Sprintf("Subject: RE: %s\r\n", old.Header.Get("Subject"))
	inreplyto := fmt.Sprintf("In-Reply-To: %s\r\n", oldid)
	references := fmt.Sprintf("References: %s%s\r\n", oldref, oldid)

	content := "\r\n\r\n" + gomua.WriteContent(os.Stdin)
	quote := bufio.NewScanner(strings.NewReader(old.SanitizeContent()))
	for quote.Scan() {
		line := quote.Text()
		token := "> "
		if line != "" && line[0] == '>' {
			token = ">"
		}
		content += "\n" + token + line
	}

	buf := bytes.NewBufferString(inreplyto)
	buf.WriteString(references)
	buf.WriteString(to)
	buf.WriteString(from)
	buf.WriteString(subject)
	buf.WriteString(content)

	reply, err := mail.ReadMessage(bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Fatal(err)
	}
	return reply
}

// adds ANSI colors to text
func color(s string, color string) string {
	return "\033[" + color + "m" + s + "\033[0m"
}

// helper func to check that view doesn't overflow []. Returns end.
func (c *client) printList(start, end int) (newstart int, newend int) {
	if end = start + c.displayN; end > len(c.messages) {
		end = len(c.messages)
	}
	m := c.messages[start:end]
	viewMailList(m, start, os.Stdout)
	start = end
	return start, end
}

// user input loop
func (c *client) input(exit chan bool) {
	var start, end int

	fmt.Println("\n\nWelcome to GoMUA! Type 'help' for help.\n")
	cli := bufio.NewScanner(os.Stdin)
	for {
		cli.Scan()
		input := cli.Text()

		switch {
		case input == "help", input == "h":
			fmt.Println(help())
		case input == "main", input == "view":
			start = 0
			start, end = c.printList(start, end)
		case input == "more":
			start, end = c.printList(start, end)
		case strings.HasPrefix(input, "reply"):
			num, err := strconv.Atoi(strings.TrimPrefix(input, "reply "))
			if err != nil {
				fmt.Println(err)
			}
			old := c.messages[num-1].(*gomua.Message)
			reply := replyMessage(old, c.user)
			gomua.Send(reply)
			old.Flag("R")
		case input == "exit", input == "x", input == "quit", input == "q":
			exit <- true
		case strings.ContainsAny(input, "01234566789"):
			num, _ := strconv.Atoi(input)
			if num <= len(c.messages) && num > 0 {
				viewMail(c.messages[num-1], os.Stdout)
			}
		}

	}

	exit <- true
}

func help() string {
	output := fmt.Sprint(
		"  help                 prints this help\n",
		"  main                 view the list of mail in your mailbox\n",
		"  more                 prints more mail listings, if not all were printed previously\n",
		"  #                    prints the details of the message #\n",
		"  reply #              prompts for the text of your reply the message #, then sends it\n",
		"  exit                 exits the program\n")

	return output
}

func main() {
	u, _ := user.Current()
	client, err := NewClient(u.HomeDir + gomua.ConfigLocation[1:])
	if err != nil {
		log.Fatal(err)
	}
	client.scanMailDir(client.dir)

	exit := make(chan bool, 1)
	go client.input(exit)
	<-exit
	os.Exit(2)
}

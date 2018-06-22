package gomua

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os/user"
	"strconv"
	"strings"
	"time"
)

const ConfigLocation = "~/.gomua/gomua.cfg"

// SMTPServer describes a connection to an SMTP server for sending mail.
type SMTPServer struct {
	name     string
	username string
	password string
	address  string
	port     int
	tlsB     bool
}

// NewSMTPServer reads from a configuration file, and returns a new SMTPServer struct ready to use.
func NewSMTPServer(filename string) (*SMTPServer, error) {
	s := new(SMTPServer)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("SMTP: missing " + ConfigLocation + " file.")
	}

	var smtp string
	sections := strings.Split(string(b), "[")
	for _, sec := range sections {
		if strings.HasPrefix(sec, "smtp]") {
			smtp = "[" + sec
		}
	}

	lines := strings.Split(smtp, "\n")
	for _, l := range lines {
		switch {
		case strings.HasPrefix(l, "Name="):
			s.name = strings.TrimPrefix(l, "Name=")
		case strings.HasPrefix(l, "Username="):
			s.username = strings.TrimPrefix(l, "Username=")
		case strings.HasPrefix(l, "Password="):
			s.password = strings.TrimPrefix(l, "Password=")
		case strings.HasPrefix(l, "Address="):
			s.address = strings.TrimPrefix(l, "Address=")
		case strings.HasPrefix(l, "Port="):
			s.port, _ = strconv.Atoi(strings.TrimPrefix(l, "Port="))
		case strings.HasPrefix(l, "TLS="):
			str := strings.TrimPrefix(l, "TLS=")
			if str == "true" {
				s.tlsB = true
			} else {
				s.tlsB = false
			}
		}
	}

	if s.name == "" || s.username == "" || s.password == "" || s.address == "" || s.port == 0 {
		return nil, errors.New("SMTP: incorrect " + ConfigLocation + " file.")
	}
	return s, nil
}

// Connects and authenticates to an SMTPServer, returns a client connection ready to write.
// This client *must be Quit()ed after finished using, preferably with defer.
func connectSMTP(s *SMTPServer) (*smtp.Client, error) {
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", s.address, s.port))
	if err != nil {
		return nil, err
	}

	if err := c.StartTLS(&tls.Config{ServerName: s.name}); err != nil {
		return nil, err
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.address)
	if err := c.Auth(auth); err != nil {
		return nil, err
	}

	return c, nil
}

// SendSMTP takes a SMTP server and a message, connects to the server, sends the message, and quits the connection to the server.
func sendSMTP(server *SMTPServer, msg *mail.Message) error {
	// connect to SMTP server
	var c *smtp.Client
	c, err := connectSMTP(server)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Quit()

	from, _ := msg.Header.AddressList("From")
	if len(from) != 0 {
		for _, t := range from {
			if err := c.Mail(t.Address); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		if err := c.Mail(server.username); err != nil {
			log.Fatal(err)
		}

	}

	to, _ := msg.Header.AddressList("To")
	for _, t := range to {
		if err := c.Rcpt(t.Address); err != nil {
			log.Fatal(err)
		}
	}

	// Send email body
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(wc, "Date: %v\r\n", time.Now().Local().Format(time.RFC822))
	for key, heads := range msg.Header {
		fmt.Fprintf(wc, "%s: ", key)
		for i, h := range heads {
			fmt.Fprint(wc, h)
			if len(heads)-1 > i {
				fmt.Fprint(wc, ",")
			}

		}
		fmt.Fprint(wc, "\r\n")
	}

	body, _ := ioutil.ReadAll(msg.Body)

	_, err = fmt.Fprintf(wc, "\n%s\n", body)
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Send opens a new SMTP server connection from the config file and sends a message.
func Send(msg *mail.Message) {
	// Look for a SMTPServer configuration file in ~/.gomua/send.cfg
	u, _ := user.Current()
	srv, err := NewSMTPServer(u.HomeDir + ConfigLocation[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nSending...")

	if sendSMTP(srv, msg) == nil {
		fmt.Println("Message Sent")
	}
}

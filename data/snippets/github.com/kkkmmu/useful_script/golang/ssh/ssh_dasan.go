package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"time"
)

type SSHSession struct {
	session *ssh.Session
	in      chan string
	out     chan string
	w       io.WriteCloser
	r       *bufio.Reader
}

func NewSSHSession(user, password, ipPort string) (*SSHSession, error) {
	sshSession := new(SSHSession)
	if err := sshSession.createConnection(user, password, ipPort); err != nil {
		return nil, err
	}
	if err := sshSession.muxShell(); err != nil {
		return nil, err
	}
	if err := sshSession.start(); err != nil {
		return nil, err
	}
	return sshSession, nil
}

func (this *SSHSession) createConnection(user, password, ipPort string) error {
	client, err := ssh.Dial("tcp", ipPort, &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 20 * time.Second,
		Config: ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com",
				"arcfour256", "arcfour128", "aes128-cbc", "aes256-cbc", "3des-cbc", "des-cbc",
			},
		},
	})
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	this.session = session
	return nil
}

func (this *SSHSession) muxShell() error {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := this.session.RequestPty("vt100", 80, 40, modes); err != nil {
		return err
	}
	w, err := this.session.StdinPipe()
	if err != nil {
		return err
	}

	this.w = w
	r, err := this.session.StdoutPipe()
	if err != nil {
		return err
	}
	this.r = bufio.NewReader(r)
	return nil
}

func (this *SSHSession) start() error {
	if err := this.session.Shell(); err != nil {
		return err
	}
	return nil
}

func (this *SSHSession) Write(buf []byte) (int, error) {
	return this.w.Write(buf)
}

func (this *SSHSession) WriteLine(data string) (int, error) {
	return this.w.Write([]byte(data + "\n"))
}

func (this *SSHSession) WriteString(data string) (int, error) {
	return this.w.Write([]byte(data))
}

func (this *SSHSession) RunCommand(cmd string) (string, error) {
	this.WriteLine(cmd)
	return "", nil
}

func (this *SSHSession) readUntil(read bool, delims ...string) ([]byte, int, error) {
	if len(delims) == 0 {
		return nil, 0, nil
	}
	p := make([]string, len(delims))
	for i, s := range delims {
		if len(s) == 0 {
			return nil, 0, nil
		}
		p[i] = s
	}
	var line []byte
	for {
		b, err := this.r.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		if read {
			line = append(line, b)
		}
		for i, s := range p {
			if s[0] == b {
				if len(s) == 1 {
					return line, i, nil
				}
				p[i] = s[1:]
			} else {
				p[i] = delims[i]
			}
		}
	}
	panic(nil)
}

func (this *SSHSession) Read(buf []byte) (int, error) {
	return this.r.Read(buf)
}

// SkipBytes works like ReadBytes but skips all read data.
func (this *SSHSession) SkipBytes(delim byte) error {
	for {
		b, err := this.r.ReadByte()
		if err != nil {
			return err
		}
		if b == delim {
			break
		}
	}
	return nil
}

// ReadString works like bufio.ReadString
func (this *SSHSession) ReadString(delim byte) (string, error) {
	bytes, err := this.r.ReadBytes(delim)
	return string(bytes), err
}

// ReadUntilIndex reads from connection until one of delimiters occurs. Returns
// read data and an index of delimiter or error.
func (this *SSHSession) ReadUntilIndex(delims ...string) ([]byte, int, error) {
	line, n, err := this.readUntil(true, delims...)
	//log.Println(string(line))
	return line, n, err
}

// ReadUntil works like ReadUntilIndex but don't return a delimiter index.
func (this *SSHSession) ReadUntil(delims ...string) ([]byte, error) {
	line, _, err := this.readUntil(true, delims...)
	//log.Println(string(line))
	return line, err
}

// SkipUntilIndex works like ReadUntilIndex but skips all read data.
func (this *SSHSession) SkipUntilIndex(delims ...string) (int, error) {
	_, i, err := this.readUntil(false, delims...)
	return i, err
}

// SkipUntil works like ReadUntil but skips all read data.
func (this *SSHSession) SkipUntil(delims ...string) error {
	_, _, err := this.readUntil(false, delims...)
	return err
}

func main() {
	cmds := make([]string, 0)
	cmds = append(cmds, "enable")
	cmds = append(cmds, "show running-config")
	cmds = append(cmds, "show system")
	cmds = append(cmds, "show ip interface brief")

	sess, err := NewSSHSession("liwei", "Lee123!@#", "10.71.20.102:22")
	if err != nil {
		panic(err)
	}

	res, err := sess.ReadUntil("SWITCH>")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(res))

	for _, cmd := range cmds {
		_, err := sess.RunCommand(cmd)
		if err != nil {
			panic(err)
		}

		res, err := sess.ReadUntil("SWITCH#")
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(string(res))
	}

	time.Sleep(5 * time.Second)
}

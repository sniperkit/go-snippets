package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type SSHSession struct {
	session     *ssh.Session
	in          chan string
	out         chan string
	brand       string
	lastUseTime time.Time
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
	sshSession.lastUseTime = time.Now()
	sshSession.brand = ""
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
	r, err := this.session.StdoutPipe()
	if err != nil {
		return err
	}

	in := make(chan string, 1024)
	out := make(chan string, 1024)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				panic(err)
				return
			}
		}()
		for cmd := range in {
			_, err := w.Write([]byte(cmd + "\n"))
			if err != nil {
				panic(err)
				return
			}
		}
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				panic(err)
				return
			}
		}()
		var (
			buf [65 * 1024]byte
			t   int
		)
		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				return
			}
			t += n
			out <- string(buf[:t])
			t = 0
		}
	}()
	this.in = in
	this.out = out
	return nil
}

func (this *SSHSession) start() error {
	if err := this.session.Shell(); err != nil {
		return err
	}
	//等待登录信息输出
	//	this.ReadChannelExpect(time.Second, "#", ">", "]")
	this.readChannelData()
	return nil
}

func (this *SSHSession) WriteChannel(cmds ...string) {
	for _, cmd := range cmds {
		this.in <- cmd
	}
}

func (this *SSHSession) ClearChannel() {
	time.Sleep(time.Millisecond * 100)
	this.readChannelData()
}

func (this *SSHSession) readChannelData() string {
	output := ""
	for {
		select {
		case channelData, ok := <-this.out:
			if !ok {
				//如果out管道已经被关闭，则停止读取，否则<-this.out会进入无限循环
				return output
			}
			output += channelData
			return output
			/*
				default:
					return output
			*/
		}
	}
}

func (this *SSHSession) RunCommands(cmds ...string) (string, error) {
	this.WriteChannel(cmds...)
	result := this.readChannelData()
	return result, nil
}

func main() {
	cmds := make([]string, 0)
	cmds = append(cmds, "show running-config")
	cmds = append(cmds, "show system")

	sess, err := NewSSHSession("test", "Lee123!@#", "10.71.20.191:22")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		go func(i int) {
			result, err := sess.RunCommands(cmds...)
			if err != nil {
				panic(err)
			}
			fmt.Printf("RunCommands<%d> result:\n%s", i, result)
		}(i)
	}

	time.Sleep(5 * time.Second)

}

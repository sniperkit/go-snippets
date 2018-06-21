// Package telnet provides simple interface for interacting with Telnet
// connection.
package telnet

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	//"time"
	"strings"
	"unicode"
)

const (
	CR = byte('\r')
	LF = byte('\n')
)

const (
	cmdSE   = 240
	cmdNOP  = 241
	cmdData = 242

	cmdBreak = 243
	cmdGA    = 249
	cmdSB    = 250

	cmdWill = 251
	cmdWont = 252
	cmdDo   = 253
	cmdDont = 254

	cmdIAC = 255
)

const (
	optEcho            = 1
	optSuppressGoAhead = 3
	//	optTerminalType    = 24
	optNAWS = 31
)

// Session implements net.Conn interface for Telnet protocol plus some set of
// Telnet specific methods.
type Session struct {
	net.Conn
	r *bufio.Reader

	unixWriteMode bool

	cliSuppressGoAhead bool
	cliEcho            bool
}

func newSession(addr string) (*Session, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	//conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	c := Session{
		Conn: conn,
		r:    bufio.NewReaderSize(conn, 256),
	}
	return &c, nil
}

// SetUnixWriteMode sets flag that applies only to the Write method.
// If set, Write converts any '\n' (LF) to '\r\n' (CR LF).
func (s *Session) SetUnixWriteMode(uwm bool) {
	s.unixWriteMode = uwm
}

func (s *Session) do(option byte) error {
	//log.Println("do:", option)
	_, err := s.Conn.Write([]byte{cmdIAC, cmdDo, option})
	return err
}

func (s *Session) dont(option byte) error {
	//log.Println("dont:", option)
	_, err := s.Conn.Write([]byte{cmdIAC, cmdDont, option})
	return err
}

func (s *Session) will(option byte) error {
	//log.Println("will:", option)
	_, err := s.Conn.Write([]byte{cmdIAC, cmdWill, option})
	return err
}

func (s *Session) wont(option byte) error {
	//log.Println("wont:", option)
	_, err := s.Conn.Write([]byte{cmdIAC, cmdWont, option})
	return err
}

func (s *Session) sub(opt byte, data ...byte) error {
	if _, err := s.Conn.Write([]byte{cmdIAC, cmdSB, opt}); err != nil {
		return err
	}
	if _, err := s.Conn.Write(data); err != nil {
		return err
	}
	_, err := s.Conn.Write([]byte{cmdIAC, cmdSE})
	return err
}

func (s *Session) deny(cmd, opt byte) (err error) {
	switch cmd {
	case cmdDo:
		err = s.wont(opt)
	case cmdDont:
		// nop
	case cmdWill, cmdWont:
		err = s.dont(opt)
	}
	return
}

func (s *Session) skipSubneg() error {
	for {
		if b, err := s.r.ReadByte(); err != nil {
			return err
		} else if b == cmdIAC {
			if b, err = s.r.ReadByte(); err != nil {
				return err
			} else if b == cmdSE {
				return nil
			}
		}
	}
}

func (s *Session) cmd(cmd byte) error {
	switch cmd {
	case cmdGA:
		return nil
	case cmdDo, cmdDont, cmdWill, cmdWont:
		// Process cmd after this switch.
	case cmdSB:
		return s.skipSubneg()
	default:
		return fmt.Errorf("unknown command: %d", cmd)
	}
	// Read an option
	o, err := s.r.ReadByte()
	if err != nil {
		return err
	}
	//log.Println("received cmd:", cmd, o)
	switch o {
	case optEcho:
		// Accept any echo configuration.
		switch cmd {
		case cmdDo:
			if !s.cliEcho {
				s.cliEcho = true
				err = s.will(o)
			}
		case cmdDont:
			if s.cliEcho {
				s.cliEcho = false
				err = s.wont(o)
			}
		case cmdWill:
			if !s.cliEcho {
				s.cliEcho = true
				err = s.do(o)
			}
		case cmdWont:
			if s.cliEcho {
				s.cliEcho = false
				err = s.dont(o)
			}
		}
	case optSuppressGoAhead:
		// We don't use GA so can allways accept every configuration
		switch cmd {
		case cmdDo:
			if !s.cliSuppressGoAhead {
				s.cliSuppressGoAhead = true
				err = s.will(o)
			}
		case cmdDont:
			if s.cliSuppressGoAhead {
				s.cliSuppressGoAhead = false
				err = s.wont(o)
			}
		case cmdWill:
			if !s.cliSuppressGoAhead {
				s.cliSuppressGoAhead = true
				err = s.do(o)
			}
		case cmdWont:
			if s.cliSuppressGoAhead {
				s.cliSuppressGoAhead = false
				err = s.dont(o)
			}
		}
	case optNAWS:
		if cmd != cmdDo {
			err = s.deny(cmd, o)
			break
		}
		if err = s.will(o); err != nil {
			break
		}
		// Reply with max window size: 65535x65535
		err = s.sub(o, 255, 255, 255, 255)
	default:
		// Deny any other option
		err = s.deny(cmd, o)
	}
	return err
}

func (s *Session) tryReadByte() (b byte, retry bool, err error) {
	b, err = s.r.ReadByte()
	if err != nil || b != cmdIAC {
		return
	}
	b, err = s.r.ReadByte()
	if err != nil {
		return
	}
	if b != cmdIAC {
		err = s.cmd(b)
		if err != nil {
			return
		}
		retry = true
	}
	return
}

// SetEcho tries to enable/disable echo on server side. Typically telnet
// servers doesn't support this.
func (s *Session) SetEcho(echo bool) error {
	if echo {
		return s.do(optEcho)
	}
	return s.dont(optEcho)
}

// ReadByte works like bufio.ReadByte
func (s *Session) ReadByte() (b byte, err error) {
	retry := true
	for retry && err == nil {
		b, retry, err = s.tryReadByte()
	}
	return
}

// ReadRune works like bufio.ReadRune
func (s *Session) ReadRune() (r rune, size int, err error) {
loop:
	r, size, err = s.r.ReadRune()
	if err != nil {
		return
	}
	if r != unicode.ReplacementChar || size != 1 {
		// Properly readed rune
		return
	}
	// Bad rune
	err = s.r.UnreadRune()
	if err != nil {
		return
	}
	// Read telnet command or escaped IAC
	_, retry, err := s.tryReadByte()
	if err != nil {
		return
	}
	if retry {
		// This bad rune was a begining of telnet command. Try read next rune.
		goto loop
	}
	// Return escaped IAC as unicode.ReplacementChar
	return
}

// Read is for implement an io.Reader interface
func (s *Session) Read(buf []byte) (int, error) {
	var n int
	for n < len(buf) {
		b, retry, err := s.tryReadByte()
		if err != nil {
			return n, err
		}
		if !retry {
			buf[n] = b
			n++
		}
		if n > 0 && s.r.Buffered() == 0 {
			// Don't block if can't return more data.
			return n, err
		}
	}
	return n, nil
}

// ReadBytes works like bufio.ReadBytes
func (s *Session) ReadBytes(delim byte) ([]byte, error) {
	var line []byte
	for {
		b, err := s.ReadByte()
		if err != nil {
			return nil, err
		}
		line = append(line, b)
		if b == delim {
			break
		}
	}
	return line, nil
}

// SkipBytes works like ReadBytes but skips all read data.
func (s *Session) SkipBytes(delim byte) error {
	for {
		b, err := s.ReadByte()
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
func (s *Session) ReadString(delim byte) (string, error) {
	bytes, err := s.ReadBytes(delim)
	return string(bytes), err
}

func (s *Session) readUntil(read bool, delims ...string) ([]byte, int, error) {
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
		b, err := s.ReadByte()
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

// ReadUntilIndex reads from connection until one of delimiters occurs. Returns
// read data and an index of delimiter or error.
func (s *Session) ReadUntilIndex(delims ...string) ([]byte, int, error) {
	line, n, err := s.readUntil(true, delims...)
	//log.Println(string(line))
	return line, n, err
}

// ReadUntil works like ReadUntilIndex but don't return a delimiter index.
func (s *Session) ReadUntil(delims ...string) ([]byte, error) {
	line, _, err := s.readUntil(true, delims...)
	//log.Println(string(line))
	return line, err
}

func (s *Session) ReadUntilSkip(delims, sdelims []string) ([]byte, error) {
	if len(sdelims) > 0 {
		result := make([]byte, 0, 64)
		for {
		AGAIN:
			line, err := s.ReadUntil(delims...)
			if err != nil {
				return result, err
			}
			result = append(result, line...)
			for _, ed := range sdelims {
				if strings.Contains(string(line), ed) {
					goto AGAIN
				}
			}

			return result, err
		}
	}

	return s.ReadUntil(delims...)
}

// SkipUntilIndex works like ReadUntilIndex but skips all read data.
func (s *Session) SkipUntilIndex(delims ...string) (int, error) {
	_, i, err := s.readUntil(false, delims...)
	return i, err
}

// SkipUntil works like ReadUntil but skips all read data.
func (s *Session) SkipUntil(delims ...string) error {
	_, _, err := s.readUntil(false, delims...)
	return err
}

// Write is for implement an io.Writer interface
func (s *Session) Write(buf []byte) (int, error) {
	search := "\xff"
	if s.unixWriteMode {
		search = "\xff\n"
	}
	var (
		n   int
		err error
	)
	for len(buf) > 0 {
		var k int
		i := bytes.IndexAny(buf, search)
		if i == -1 {
			k, err = s.Conn.Write(buf)
			n += k
			break
		}
		k, err = s.Conn.Write(buf[:i])
		n += k
		if err != nil {
			break
		}
		switch buf[i] {
		case LF:
			k, err = s.Conn.Write([]byte{CR, LF})
		case cmdIAC:
			k, err = s.Conn.Write([]byte{cmdIAC, cmdIAC})
		}
		n += k
		if err != nil {
			break
		}
		buf = buf[i+1:]
	}
	return n, err
}

func (s *Session) WriteString(str string) (int, error) {
	return s.Write([]byte(str))
}

func (s *Session) WriteLine(str string) (int, error) {
	return s.WriteString(str + "\n")
}

func New(user, pass, ip, port string) (*Session, error) {
	s, err := newSession(ip + ":" + port)
	if err != nil {
		log.Println("error happend when connect to: ", ip, " with: ", err.Error())
		return nil, errors.New("Cannot connect to host")
	}

	s.WriteLine("\n") //For serial server
	s.SetUnixWriteMode(true)
	_, err = s.ReadUntil("ogin:")
	if err != nil {
		fmt.Println("Error happend when get login: ", err.Error())
		return nil, err
	}

	s.WriteLine(user)
	_, err = s.ReadUntil("assword:")
	if err != nil {
		fmt.Println("Error happend when get login prompt: ", err.Error())
		return nil, err
	}

	s.WriteLine(pass)
	_, err = s.ReadUntil(">")
	if err != nil {
		fmt.Println("Error happend when login: ", err.Error())
		return nil, err
	}

	return s, nil
}

/* New Telnet client with login prompt, password prompt and normal prompt */
func New2(user, pass, ip, port, up, pp, np string) (*Session, error) {
	s, err := newSession(ip + ":" + port)
	if err != nil {
		log.Println("error happend when connect to: ", ip, " with: ", err.Error())
		return nil, errors.New("Cannot connect to host")
	}

	s.WriteLine("\n") //For serial server
	s.SetUnixWriteMode(true)
	_, err = s.ReadUntil(up)
	if err != nil {
		fmt.Println("Error happend when get login: ", err.Error())
		return nil, err
	}

	s.WriteLine(user)
	_, err = s.ReadUntil(pp)
	if err != nil {
		fmt.Println("Error happend when get login prompt: ", err.Error())
		return nil, err
	}

	s.WriteLine(pass)
	_, err = s.ReadUntil(np)
	if err != nil {
		fmt.Println("Error happend when login: ", err.Error())
		return nil, err
	}

	return s, nil
}

func New3(addr string) (*Session, error) {
	s, err := newSession(addr)
	if err != nil {
		log.Println("error happend when connect to: ", addr, " with: ", err.Error())
		return nil, errors.New("Cannot connect to host")
	}
	s.SetUnixWriteMode(true)

	return s, nil
}

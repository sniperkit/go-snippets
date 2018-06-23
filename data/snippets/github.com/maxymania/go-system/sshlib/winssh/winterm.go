/*
 * Copyright(C) 2016 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 */

/*
 This package implements the session handling (experimental) on windows machines.
 */
package winssh

import "fmt"
import "io"
import "bufio"
import "os/exec"
import "golang.org/x/crypto/ssh"
import "github.com/maxymania/go-system/sshlib"

func HandleSess(sess *sshlib.ShellSession, cmd *exec.Cmd) {
	defer sess.Ch.Close()
	xout,_ := cmd.StdoutPipe()
	xerr,_ := cmd.StderrPipe()
	xin,_ := cmd.StdinPipe()
	cmd.Start()
	go copyout(xout,sess.Ch)
	go copyerr(xerr,sess.Ch)
	go copyin(xin,sess.Ch,cmd)
	cmd.Wait()
}

func copyout(i io.ReadCloser,c ssh.Channel){
	defer i.Close()
	io.Copy(c,i)
}
func copyerr(i io.ReadCloser,c ssh.Channel){
	defer i.Close()
	io.Copy(c.Stderr(),i)
}

func copyin(o io.WriteCloser,c ssh.Channel, cmd *exec.Cmd){
	defer o.Close()
	cc := bufio.NewReader(c)
	cmdbuffer := make([]byte,0,200)
	esc := make([]byte,10)
	esci := 0
	buffer := make([]byte,1)
	for{
		b,e := cc.ReadByte()
		if e!=nil { return }
		
		/*
		if b>=0x20 {
			fmt.Printf("Char '%c' \\x%02x\n",int(b),int(b))
		}else{
			fmt.Printf("Char '#' \\x%02x\n",int(b))
		}
		*/
		
		switch b{
		case '\n': break // ignore
		case '\r':
			cmdbuffer = append(cmdbuffer,'\r','\n')
			o.Write(cmdbuffer)
			cmdbuffer = cmdbuffer[:0]
			fmt.Fprint(c,"\r\n")
		case 0x7f:
			l := len(cmdbuffer)-1
			if l>=0 {
				cmdbuffer = cmdbuffer[:l]
				buffer[0]=b
				c.Write(buffer)
			}
		case 0x1b: // ANSI / VT100 escape sequence
			esci=0
			esc[esci]='^'; esci++
			esc[esci],e = cc.ReadByte() ; if e!=nil { return }
			esci++
			switch esc[esci-1]{
			default:
				c.Write(esc[:esci])
				cmdbuffer = append(cmdbuffer,esc[:esci]...)
			}
		case 0x04: //ctrl-d
			fmt.Fprintf(c,"\r\ntype 'exit' to logout\r\n")
		case 0x03: //ctrl-c
			cmd.Process.Kill()
			o.Close()
		default:
			buffer[0]=b
			cmdbuffer = append(cmdbuffer,b)
			c.Write(buffer)
		}
	}
}

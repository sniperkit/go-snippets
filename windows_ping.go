package main


import (
	"os/exec"
	"bufio"

	"github.com/axgle/mahonia"
	"fmt"
	"io"
)

type PingSt struct {
	SendPk			string
	RecvPk			string
	LossPk			string
	MinDelay		string
	AvgDelay		string
	MaxDelay		string
}


func main() {
	//pingWindows("139.198.11.147","5")
	pingLinux("139.198.11.147","10")
}

func pingWindows(Addr string, cnt string) {
	//var ps PingSt
	cmd := exec.Command("ping","-n", cnt, Addr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	for {
		l, Merr := reader.ReadString('\n')
		if Merr != nil || io.EOF == Merr{
			break
		}
		var line string
		var dec mahonia.Decoder
		dec = mahonia.NewDecoder("gbk")
		line = dec.ConvertString(l)
		fmt.Println(line)
	}
	cmd.Wait()
}

func pingLinux(Addr string, cnt string) {
	cmd := exec.Command("ping","-c", cnt, Addr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	cmd.Start()
	reader := bufio.NewReader(stdout)

	for {
		line, Merr := reader.ReadString('\n')
		if Merr != nil || io.EOF == Merr {
			break
		}

		fmt.Println(line)

	}
	cmd.Wait()
}
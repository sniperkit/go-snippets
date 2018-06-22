package scpm

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gronpipmaster/pb"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var homeFolder = os.Getenv("HOME")

type Host struct {
	User     string
	Addr     string
	Output   string
	Client   *ssh.Client
	sess     *ssh.Session
	Identity *ssh.ClientConfig
}

func NewHost(host string, key string, port int) (h *Host, err error) {
	h = new(Host)
	if strings.Index(host, "@") == -1 {
		h.User = os.Getenv("LOGNAME")
	} else {
		arrHost := strings.Split(host, "@")
		h.User = arrHost[0]
		host = strings.Replace(host, h.User+"@", "", -1)
	}
	if strings.Index(host, ":") == -1 {
		err = errors.New("host incorrect")
		return
	}
	arrHost := strings.Split(host, ":")
	h.Addr = arrHost[0] + ":" + fmt.Sprint(port)
	h.Output = arrHost[1]
	keys := []string{
		key,
		homeFolder + "/.ssh/id_rsa",
		homeFolder + "/.ssh/id_dsa",
		homeFolder + "/.ssh/id_ecdsa",
	}
	h.Identity = &ssh.ClientConfig{User: h.User}
	for _, k := range keys {
		//create ssh.Config from private keys
		f, err := os.Open(k)
		if err != nil {
			continue
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return h, err
		}
		signer, err := ssh.ParsePrivateKey(data)
		if err != nil {
			return h, err
		}
		h.Identity.Auth = append(h.Identity.Auth, ssh.PublicKeys(signer))
		f.Close()
	}
	return
}

func (h Host) String() string {
	str := fmt.Sprintf("%s ", h.Addr+" "+h.Output)
	//TODO fixed str size if > 50
	return str
}

func (h *Host) Auth() error {
	var err error
	h.Client, err = ssh.Dial("tcp", h.Addr, h.Identity)
	if err != nil {
		h.Identity.Auth = append(h.Identity.Auth, ssh.Password(dialogPassword(h.Addr)))
		h.Client, err = ssh.Dial("tcp", h.Addr, h.Identity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Host) Copy(tree *Tree, wg *sync.WaitGroup, bar *pb.ProgressBar) {
	defer func() {
		bar.Finish()
		wg.Done()
	}()
	for _, file := range tree.Files {
		in := file.Dir + string(os.PathSeparator) + file.Info.Name()
		out := strings.Replace(in, tree.BaseDir, h.Output, -1)
		if err := h.cp(in, out, bar); err != nil && err != io.EOF {
			log.Println(err)
		}
	}
}

func (h *Host) exec(cmd string) error {
	var err error
	h.sess, err = h.Client.NewSession()
	if err != nil {
		return err
	}
	defer h.sess.Close()
	return h.sess.Run(cmd)
}

const (
	cmdCat   string = "cat > %s"
	cmdStat  string = "stat %s"
	cmdMkDir string = "mkdir -p %s"
)

//remote cp
func (h *Host) cp(in, out string, bar *pb.ProgressBar) error {
	var err error
	//create remote dir
	dir := filepath.Dir(out)
	//check folder exists
	if err := h.exec(fmt.Sprintf(cmdStat, dir)); err != nil {
		if err := h.mkdir(dir); err != nil {
			return err
		}
	}
	//open fd file
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()
	//open ssh session
	h.sess, err = h.Client.NewSession()
	if err != nil {
		return err
	}
	defer h.sess.Close()
	dest, err := h.sess.StdinPipe()
	if err != nil {
		return err
	}
	if err = h.sess.Start(fmt.Sprintf(cmdCat, out)); err != nil {
		return err
	}
	writer := io.MultiWriter(dest, bar)
	_, err = io.Copy(writer, f)
	return err
}

//remote mkdir
func (h *Host) mkdir(dir string) error {
	return h.exec(fmt.Sprintf(cmdMkDir, dir))
}

type Scp struct {
	hosts []*Host
	tree  *Tree
	wg    *sync.WaitGroup
	done  chan bool
}

func New(hosts []*Host, path string) (scp *Scp, err error) {
	if strings.Index(path, "~") != -1 {
		path = strings.Replace(path, "~", homeFolder, -1)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	fmt.Println("Start copy", absPath)
	if len(hosts) == 0 {
		err = errors.New("hosts is nil")
		return
	}
	scp = new(Scp)
	scp.hosts = hosts
	scp.tree, err = NewTree(absPath)
	if err != nil {
		return
	}
	scp.wg = new(sync.WaitGroup)
	scp.done = make(chan bool)
	return
}

func (s *Scp) Run(quit chan bool) {
	pool := &pb.Pool{}
	for _, host := range s.hosts {
		if err := host.Auth(); err != nil {
			fmt.Println("Auth err:", err)
			continue
		}
		s.wg.Add(1)
		bar := pb.New(int(s.tree.Size)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
		bar.ShowSpeed = true
		bar.Prefix(host.String())
		pool.Add(bar)
		go host.Copy(s.tree, s.wg, bar)
	}
	pool.Start()
	go func() {
		s.wg.Wait()
		time.Sleep(10 * time.Millisecond)
		s.done <- true
	}()
	for {
		select {
		case <-s.done:
			quit <- true
		}
	}
}

type Tree struct {
	BaseDir string
	Size    int64
	Files   []File
}

type File struct {
	Info os.FileInfo
	Dir  string
}

func NewTree(path string) (t *Tree, err error) {
	t = new(Tree)
	t.BaseDir = filepath.Dir(path)
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	if info.IsDir() {
		err = filepath.Walk(path, t.Scan)
		return
	}
	t.Files = append(t.Files, File{Info: info, Dir: t.BaseDir})
	t.Size = info.Size()
	return
}

func (t *Tree) Scan(path string, fileInfo os.FileInfo, errInp error) (err error) {
	if errInp != nil {
		log.Println(errInp)
		return nil
	}
	if fileInfo.IsDir() {
		return nil
	}
	t.Files = append(t.Files, File{Info: fileInfo, Dir: filepath.Dir(path)})
	t.Size += fileInfo.Size()
	return
}

func dialogPassword(addr string) (pass string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(addr)
	fmt.Print("Enter Password: ")
	pass, _ = reader.ReadString('\n')
	return strings.TrimSpace(pass)
}

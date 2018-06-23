package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func SaveToFile(name string, data []byte) {
	file, err := os.Create(name)
	if err != nil {
		log.Println("Cannot create file: ", name, " ", err.Error())
		return
	}

	file.Write(data)
	file.Sync()
	defer file.Close()
}

func AppendToFile(name string, data []byte) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(data)
}

func GenerateSessionID() string {
	var id = "0x"

	rand.Seed(time.Now().Unix())
	result := rand.Perm(128)
	for _, i := range result {
		id = id + strconv.Itoa(i)
	}

	return id
}

func GenerateSessionIDByUserNameAndPassword(name, pass string) string {
	hash := sha1.New()
	return hex.EncodeToString(hash.Sum([]byte(name + pass)))
}

func Download(link, name string) error {
	if name == "" || link == "" {
		return errors.New("Invalid paramenters: " + link + ":" + name)
	}

	resp, err := http.Get(link)
	if err != nil {
		return errors.New("Do http request error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Do http request error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Read http request body error: " + err.Error())
	}

	if err := ioutil.WriteFile(name, body, 0644); err != nil {
		return errors.New("Write content to file error: " + err.Error())
	}

	return nil
}

func Wait(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func DiffFile(before string, after string) {
	os.Remove(after + ".line.diff")
	content1, _ := ioutil.ReadFile(before)
	content2, _ := ioutil.ReadFile(after)

	lines1 := strings.Split(string(content1), "\n")
	lines2 := strings.Split(string(content2), "\n")

	for i, line := range lines1 {
		if strings.TrimSpace(line) == "" {
			continue
		}
		dmp := diffmatchpatch.New()

		diffs := dmp.DiffMain(string(line), string(lines2[i]), false)

		if len(diffs) == 1 {
			continue
		}

		var buff bytes.Buffer
		for _, d := range diffs {
			if d.Type == diffmatchpatch.DiffInsert {
				_, _ = buff.WriteString("\x1b[32m")
				_, _ = buff.WriteString(d.Text)
				_, _ = buff.WriteString("\x1b[0m")
			} else if d.Type == diffmatchpatch.DiffDelete {
				_, _ = buff.WriteString("\x1b[31m")
				_, _ = buff.WriteString(d.Text)
				_, _ = buff.WriteString("\x1b[0m")
			} else if d.Type == diffmatchpatch.DiffEqual {
				_, _ = buff.WriteString(d.Text)
			}
		}
		AppendToFile(after+".line.diff", buff.Bytes())
	}
}

func DiffFileToText(after string, before string) {
	os.Remove(after + ".diff")
	content1, _ := ioutil.ReadFile(after)
	content2, _ := ioutil.ReadFile(before)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(content1), string(content2), false)
	SaveToFile(after+".diff", []byte(dmp.DiffPrettyText(diffs)))
}

func DiffFileToHtml(after string, before string) {
	os.Remove(after + ".html")
	content1, _ := ioutil.ReadFile(after)
	content2, _ := ioutil.ReadFile(before)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(content1), string(content2), false)
	SaveToFile(after+".html", []byte(fmt.Sprintf("<html><head><title>%s</title></head><body>%s</body></html>", after, dmp.DiffPrettyHtml(diffs))))
}

func DiffToText(after, before string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(after, before, false)
	SaveToFile("diff.diff", []byte(dmp.DiffPrettyText(diffs)))
}

func DiffToHtml(after string, before string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(after, before, false)
	SaveToFile("diff.html", []byte(fmt.Sprintf("<html><head><title>%s</title></head><body>%s</body></html>", after, dmp.DiffPrettyHtml(diffs))))
}

func Diff(after, before string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(after, before, false)

	fmt.Println(dmp.DiffPrettyText(diffs))
}

func Diff2(after, before string) {
	lines1 := strings.Split(string(after), "\n")
	lines2 := strings.Split(string(before), "\n")

	for i, line := range lines1 {
		if strings.TrimSpace(line) == "" {
			continue
		}
		dmp := diffmatchpatch.New()

		diffs := dmp.DiffMain(string(line), string(lines2[i]), false)

		if len(diffs) == 1 {
			continue
		}

		var buff bytes.Buffer
		for _, d := range diffs {
			if d.Type == diffmatchpatch.DiffInsert {
				_, _ = buff.WriteString("\x1b[32m")
				_, _ = buff.WriteString(d.Text)
				_, _ = buff.WriteString("\x1b[0m")
			} else if d.Type == diffmatchpatch.DiffDelete {
				_, _ = buff.WriteString("\x1b[31m")
				_, _ = buff.WriteString(d.Text)
				_, _ = buff.WriteString("\x1b[0m")
			} else if d.Type == diffmatchpatch.DiffEqual {
				_, _ = buff.WriteString(d.Text)
			}
		}
		fmt.Println(string(buff.Bytes()))
	}
}

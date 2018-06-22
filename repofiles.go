package repofiles

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Repo struct {
	user     string
	repo     string
	revision string
	list     List
	files    []File
}

type List struct {
	Sha  string
	Url  string
	Tree []struct {
		Path   string
		Mode   string
		Format string
		Sha    string
		Size   int
		Url    string
	}
	Trunctated bool
}

type File struct {
	Path     string
	Contents string
}

type Credentials struct {
	User  string
	Token string
}

func (r *Repo) List(credentials Credentials) List {
	client := &http.Client{}
	req, err := http.NewRequest("GET", r.listEndpoint(), nil)
	if credentials.User != "" {
		req.SetBasicAuth(credentials.User, credentials.Token)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	list := List{}
	json.Unmarshal([]byte(respBody), &list)

	r.list = list

	return list
}

func (r *Repo) Files(filter string, credentials Credentials) []File {
	re := regexp.MustCompile(filter)
	type jsonFile struct {
		Sha      string
		Size     int
		Url      string
		Content  string
		Encoding string
	}
	var file File
	var parsedFile jsonFile
	client := &http.Client{}
	for _, v := range r.list.Tree {
		file = File{Path: v.Path}
		if v.Url != "" && re.Match([]byte(file.Name())) {
			req, err := http.NewRequest("GET", v.Url, nil)
			if credentials.User != "" {
				req.SetBasicAuth(credentials.User, credentials.Token)
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Error : %s", err)
			}
			respBody, _ := ioutil.ReadAll(resp.Body)

			json.Unmarshal([]byte(respBody), &parsedFile)
			contents, _ := base64.StdEncoding.DecodeString(parsedFile.Content)
			r.files = append(r.files, File{Path: v.Path, Contents: string(contents)})
		}
	}
	return r.files
}

func (r *Repo) listEndpoint() string {
	segments := []string{
		"https://api.github.com/repos", r.user, r.repo,
		"git/trees/" + r.revision + "?recursive=1",
	}
	return strings.Join(segments, "/")
}

func (f *File) Name() string {
	steps := strings.Split(f.Path, "/")
	return steps[len(steps)-1]
}

func NewRepo(user string, repo string, revision string) Repo {
	return Repo{user: user, repo: repo, revision: revision}
}

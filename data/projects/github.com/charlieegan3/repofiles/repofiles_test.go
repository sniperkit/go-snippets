package repofiles

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var credentials Credentials

func TestMain(m *testing.M) {
	credentials = Credentials{User: os.Getenv("GITHUB_USER"), Token: os.Getenv("GITHUB_TOKEN")}
	fmt.Println(credentials)
	code := m.Run()
	os.Exit(code)
}

func TestInit(t *testing.T) {
	repo := NewRepo("charlieegan3", "dotfiles", "efba099dae2acc7daac7d0638ba1be39a45914ab")
	expected := Repo{user: "charlieegan3", repo: "dotfiles", revision: "efba099dae2acc7daac7d0638ba1be39a45914ab"}
	assert.Equal(t, expected, repo, "should be equal")
}

func TestList(t *testing.T) {
	repo := NewRepo("charlieegan3", "dotfiles", "efba099dae2acc7daac7d0638ba1be39a45914ab")
	list := repo.List(credentials)
	var files []string
	for _, v := range list.Tree {
		files = append(files, v.Path)
	}
	expectedFiles :=
		[]string{".bashrc", ".gitconfig", ".gitignore", "com.googlecode.iterm2.plist"}
	assert.Equal(t, expectedFiles, files, "should be equal")
}

func TestFiles(t *testing.T) {
	repo := NewRepo("charlieegan3", "dotfiles", "efba099dae2acc7daac7d0638ba1be39a45914ab")
	repo.List(credentials)
	repo.Files("", credentials)

	assert.Equal(t, repo.files[0].Path, ".bashrc", "should be equal")
}

func TestFilesFilter(t *testing.T) {
	repo := NewRepo("charlieegan3", "dotfiles", "efba099dae2acc7daac7d0638ba1be39a45914ab")
	repo.List(credentials)
	repo.Files("bash", credentials)

	assert.Equal(t, len(repo.files), 1, "should be equal")
	assert.Equal(t, repo.files[0].Path, ".bashrc", "should be equal")
}

func TestFileName(t *testing.T) {
	file := File{Path: "./folder1/folder2/file.ext"}
	assert.Equal(t, "file.ext", file.Name(), "should be equal")
	file = File{Path: "file.ext"}
	assert.Equal(t, "file.ext", file.Name(), "should be equal")
}

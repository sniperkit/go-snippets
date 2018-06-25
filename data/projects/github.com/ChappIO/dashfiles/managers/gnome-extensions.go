package managers

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
)

type GnomeExtensions struct {
	PackageManager
}

const chromeGnomeShell = "chrome-gnome-shell"

var EPSILON float64 = 0.00000001

var gnomeExec *exec.Cmd
var writer io.WriteCloser
var reader io.ReadCloser
var started bool

func init() {
	Register(&GnomeExtensions{})
}

func (gnome *GnomeExtensions) GetName() string {
	return "gnome_extensions"
}

func (gnome *GnomeExtensions) GetFileFormat() string {
	return "packages"
}

func (gnome *GnomeExtensions) IsInstalled() bool {
	_, err := exec.LookPath(chromeGnomeShell)
	return err == nil
}

func (gnome *GnomeExtensions) InstallFromFile(file string) error {
	if !started {
		if err := gnomeExec.Start(); err != nil {
			log.Fatal(err)
		}
		started = true
	}
	packages := getToBeInstalledPackages(gnome, file)
	for _, extension := range packages {
		if err := Install(extension); err != nil {
			fmt.Println("Error: Could not install gnome extension " + extension + ": " + err.Error())
		}
	}
	return nil
}

func (gnome *GnomeExtensions) Update() error {
	return nil
}

func (gnome *GnomeExtensions) GetPackages() []string {
	// Go packages are reinstalled every time
	exts, err := ListExtensions()
	if err != nil {
		fmt.Println("Failed to list extensions: " + err.Error())
		return nil
	}
	result := make([]string, 0)
	for _, extension := range exts {
		if math.Abs(extension.State-1.0) < EPSILON {
			result = append(result, extension.Uuid)
		}
	}
	return result
}

func init() {
	gnomeExec = exec.Command("chrome-gnome-shell")
	gnomeExec.Stderr = os.Stderr
	writer, _ = gnomeExec.StdinPipe()
	reader, _ = gnomeExec.StdoutPipe()
}

type Extension struct {
	Uuid        string  `json:"uuid"`
	Description string  `json:"description"`
	Url         string  `json:"url"`
	Path        string  `json:"path"`
	Name        string  `json:"name"`
	State       float64 `json:"state"`
}

func ListExtensions() (exts map[string]Extension, err error) {
	data, err := execute(`{"execute": "listExtensions"}`)
	var extensions struct {
		Extensions map[string]Extension `json:"extensions"`
	}
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &extensions)
	exts = extensions.Extensions
	return
}

func Install(uuid string) (err error) {
	_, err = execute(`{
		"execute": "installExtension",
		"uuid": "` + uuid + `"
	}`)
	return
}

func execute(message string) (out []byte, err error) {
	byteMessage := []byte(message)
	binary.Write(writer, binary.LittleEndian, int32(len(byteMessage)))
	writer.Write(byteMessage)
	var messageLength int32
	binary.Read(reader, binary.LittleEndian, &messageLength)
	data := make([]byte, messageLength)
	reader.Read(data)
	out = data
	return
}

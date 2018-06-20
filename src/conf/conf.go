package conf

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

var ErrConfNotExist = errors.New("Config File Not Exist")

const (
	dirFileMode  os.FileMode = 0755
	confFileMode os.FileMode = 0644
)

func createDirectories(pathStr string) error {
	return os.MkdirAll(path.Dir(pathStr), dirFileMode)
}

func createConfigFile(pathStr string, deft interface{}) error {
	data, err := json.MarshalIndent(deft, "", "\t")
	if err != nil {
		return err
	}

	if err = createDirectories(pathStr); err != nil {
		return err
	}

	if err = ioutil.WriteFile(pathStr, data, confFileMode); err != nil {
		return err
	}

	return ErrConfNotExist
}

// LoadConfig load config data from config file in pathStr, if config file is not exist,
// this function while create a new config file according to deft
func LoadConfig(pathStr string, deft interface{}, val interface{}) error {
	file, err := os.Open(pathStr)
	if err != nil {
		if os.IsNotExist(err) {
			return createConfigFile(pathStr, deft)
		}

		return err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, val)
}

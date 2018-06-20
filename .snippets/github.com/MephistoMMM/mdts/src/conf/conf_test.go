package conf

import (
	"os"
	"testing"
)

type TestConf struct {
	K1 int    `json:"Key1"`
	K2 string `json:"Key2"`
}

var defaultConf = TestConf{
	K1: 1,
	K2: "2",
}

func TestLoadConf(t *testing.T) {
	var config TestConf

	if err := LoadConfig("./test/config.json", &defaultConf, &config); err != nil {
		if err == ErrConfNotExist {
			t.Log("Config File Not Exist, Created It.")
		} else {
			t.Error(err)
		}
	}

	t.Log(config)

}

func TestEnvConf(t *testing.T) {
	confMap := map[string]string{
		"app": "1",
		"bpp": "2",
	}

	InitConfMapFromEnv(confMap)
	if confMap["bpp"] != "2" {
		t.Error("value of 'bpp' in confMap is error.")
	}

	os.Setenv("BPP", "3")
	InitConfMapFromEnv(confMap)
	if confMap["bpp"] != "3" {
		t.Error("value of 'bpp' in confMap is error.")
	}
}

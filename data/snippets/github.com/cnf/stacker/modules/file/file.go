package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
	"gopkg.in/fsnotify.v1"

	"github.com/cnf/stacker/engine"
)

const graceTime = 3 * time.Second

// ConfStore object
type ConfStore struct {
	Config   *Config
	eng      *engine.Engine
	fileList map[string]time.Time
	stopper  chan bool
}

// File represents a single config file
type File struct {
	Path string
	MD5  string
}

// NewConfStore returns a new file ConfStore object
func NewConfStore() engine.ConfStore {
	return &ConfStore{}
}

// RegisterConfStore returns the module name
func RegisterConfStore() {
	engine.RegisterConfStore("file", NewConfStore)
}

// Setup the file confstore
func (cs *ConfStore) Setup(eng *engine.Engine, cfg engine.Config) error {
	cs.Config = cfg.(*Config)
	cs.eng = eng
	cs.fileList = make(map[string]time.Time)
	return nil
}

// Start the ConfStore
func (cs *ConfStore) Start() error {
	if glog.V(1) {
		glog.Infof("starting file configstore")
	}
	// TODO: make this work
	cs.stopper = make(chan bool)
	cs.parseFiles()
	if cs.Config.Watch {
		go cs.watch()
	}
	// TODO: parse files on start
	return nil
}

// Stop the ConfStore
func (cs *ConfStore) Stop() error {
	close(cs.stopper)
	return nil
}

func (cs *ConfStore) parseFiles() error {
	list, err := getFileList(cs.Config.ConfigDir)
	if err != nil {
		// TODO: error handeling
		return err
	}
	for _, v := range list {
		cs.fileOp(v)
	}
	return nil
}

func (cs *ConfStore) watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		// TODO: Error handeling
		glog.Error(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				cs.fileOp(event.Name)
			case err := <-watcher.Errors:
				glog.Error(err)
			case <-cs.stopper:
				return
			}
		}
	}()

	if err := watcher.Add(cs.Config.ConfigDir); err != nil {
		// TODO: ...
		glog.Error(err)
	}

	<-cs.stopper
	glog.Info("stopping file watcher")
}

func (cs *ConfStore) fileOp(file string) {
	// glog.Info(event)
	// TODO: sanitize
	name := filepath.Clean(file)
	fp, err := os.Open(name)
	if err != nil {
		// TODO: errors
		glog.Error(err.Error())
		return
	}
	defer fp.Close()
	fi, err := fp.Stat()
	if err != nil {
		// TODO: errors
		glog.Error(err.Error())
	}
	if fi.IsDir() {
		return
	}
	ext := filepath.Ext(name)
	if ext != ".toml" {
		glog.Info("not a toml file")
		return
	}
	// TODO: grace time
	// if time.Now().Sub(fi.ModTime()) < graceTime {
	// 	glog.Info("modified less than 3 seconds ago")
	// 	return
	// }
	if _, ok := cs.fileList[name]; ok {
		if fi.ModTime().Sub(cs.fileList[name]) < graceTime {
			glog.Info("modified less than 3 seconds ago")
			return
		}
	}
	cs.fileList[name] = fi.ModTime()
	data := make(map[string]interface{})
	_, err = toml.DecodeReader(fp, data)
	if err != nil {
		// TODO: errors
		glog.Error(err.Error())
	}
	cs.sentEvent(data)
}

func (cs *ConfStore) sentEvent(data map[string]interface{}) error {
	src := &engine.NameSpace{
		Type:   "confstore",
		Module: "file",
	}
	// we take the section heading from the toml file, and add it to _type
	// so for sectiom [[container]], we do _type = container
	for k, v := range data {
		switch v.(type) {
		case map[string]interface{}:
			v.(map[string]interface{})["_type"] = k
			cmd := &engine.Command{
				Name:    "config",
				Source:  src,
				Payload: v,
			}
			cs.eng.Event(cmd)
		case []map[string]interface{}:
			for _, iv := range v.([]map[string]interface{}) {
				// iv.(map[string]interface{})["_type"].(string) = k
				iv["_type"] = k
				cmd := &engine.Command{
					Name:    "config",
					Source:  src,
					Payload: iv,
				}
				cs.eng.Event(cmd)
			}
		default:
			glog.Info("not the right config format")
			continue
		}
	}
	return nil
}

func getFileList(cfgpath string) ([]string, error) {
	var list []string

	// TODO: sanitize paths
	cfgpath = filepath.Clean(cfgpath)

	files, _ := ioutil.ReadDir(cfgpath)
	for _, v := range files {
		if v.Mode().IsDir() {
			rpath := filepath.Join(cfgpath, v.Name())
			rlist, err := getFileList(rpath)
			if err != nil {
				continue
			}
			list = append(list, rlist...)
			continue
		}

		newpath := filepath.Join(cfgpath, v.Name())
		list = append(list, newpath)

		if !v.Mode().IsRegular() {
			continue
		}
	}
	return list, nil
}

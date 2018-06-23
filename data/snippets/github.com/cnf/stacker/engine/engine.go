package engine

import (
	// "fmt"
	"errors"
	"path/filepath"
	"time"

	"github.com/golang/glog"

	"github.com/kr/pretty"
)

// Engine is the heart of stacker, it runs everything.
type Engine struct {
	Config  map[string]Config
	cfgfile string

	ContainerRuntime ContainerRuntime
	ContainerList    map[string]Container
	FileStore        ConfStore
	ConfStore        ConfStore
	Action           map[string]Action
	Reaction         map[string]Reaction

	stopper chan bool
	ch      chan *Command
	chcount int
	// actionlist map[string][]*Reaction
}

const chlen int = 5

// New creates a new Engine object, and returns it.
func New() (*Engine, error) {
	eng := &Engine{}
	return eng, nil
}

// State gives a dump of the engine state
// This is for testing, and will be removes.
func (e *Engine) State() error {
	glog.Info("==========")
	// glog.Infof("configDecoderList length: %d", len(configDecoderList))
	// glog.Infof("confStoreList length: %d", len(confStoreList))
	// glog.Infof("actionList length: %d", len(actionList))
	// glog.Infof("reactionList length: %d", len(reactionList))
	glog.Infof("items in ContainerList: %d", len(e.ContainerList))
	glog.Info("==========")
	return nil
}

// SetConfigFile sets the path to the config file
func (e *Engine) SetConfigFile(cfgfile string) error {
	if e.cfgfile == "" {
		cdp, err := filepath.Abs(cfgfile)
		if err != nil {
			return err
		}
		e.cfgfile = filepath.Clean(cdp)
	}
	return nil
}

// Initialize the engine
func (e *Engine) Initialize() error {
	glog.Info("initializing the engine")
	if e.ch == nil {
		e.ch = make(chan *Command, chlen)
	}
	if e.Config == nil {
		e.Config = make(map[string]Config)
	}
	tmp, err := parse(e.cfgfile)
	if err != nil {
		return err
	}
	for k, v := range tmp {
		switch v.(type) {
		case map[string]interface{}:
			cfg, err := e.parseConfig(k, tmp[k].(map[string]interface{}))
			if err != nil {
				glog.Warningf("config decoder for [%s] failed: %s", k, err.Error())
				continue
			}
			e.Config[k] = cfg
		}
	}
	e.ensureDefault("stacker")
	e.ensureDefault("docker")

	if err := e.setupContainerRuntime(); err != nil {
		glog.Fatal(err.Error())
	}
	e.setupConfStore()
	e.setupActions()
	e.setupReactions()

	// TODO: if err := dw.Connect(cfg.Docker); err != nil {
	// TODO: if err := dockr.BuildNewLists(); err != nil {
	return nil
}

// Start the engine
func (e *Engine) Start() error {
	glog.Info("starting the engine")
	if e.ch == nil {
		e.ch = make(chan *Command, chlen)
	}
	e.stopper = make(chan bool)
	go e.ticker()
	go e.runner()

	if e.ContainerRuntime == nil {
		return errors.New("can not run without a container runtime")
	}
	if err := e.ContainerRuntime.Start(); err != nil {
		return err
	}

	if err := e.FileStore.Start(); err != nil {
		return err
	}
	// TODO: no confstore should not crash
	if e.ConfStore != nil {
		e.ConfStore.Start()
	}
	for k := range e.Action {
		e.Action[k].Start()
	}
	for k := range e.Reaction {
		e.Reaction[k].Start()
	}
	return nil
}

// Stop the engine
func (e *Engine) Stop() error {
	glog.Info("stopping the engine")

	e.ContainerRuntime.Stop()
	e.FileStore.Stop()
	if e.ConfStore != nil {
		e.ConfStore.Stop()
	}
	for k := range e.Action {
		e.Action[k].Stop()
	}
	for k := range e.Reaction {
		e.Reaction[k].Stop()
	}
	close(e.stopper)
	return nil
}

// Restart the engine.
func (e *Engine) Restart() error {
	if err := e.Stop(); err != nil {
		return err
	}
	// TODO: reinitialize?
	if err := e.Initialize(); err != nil {
		return err
	}
	if err := e.Start(); err != nil {
		return err
	}
	return nil
}

// Event triggers an incoming event
// TODO: clean up
func (e *Engine) Event(cmd *Command) error {
	// glog.Infof("+ putting '%s' on channel with fill %d", cmd.Name, e.chcount)
	e.ch <- cmd
	e.chcount++
	return nil
}

// RegisterContainer registers an container on the container list
func (e *Engine) RegisterContainer(name string, c Container) error {
	if _, ok := e.ContainerList[name]; ok {
		// TODO: already exists, do we replace?
	}
	e.ContainerList[name] = c
	return nil
}

// DeregisterContainer registers an action on the container list
func (e *Engine) DeregisterContainer(name string) error {
	if _, ok := e.ContainerList[name]; ok {
		e.ContainerList[name] = nil
		return nil
	}
	return errors.New("container did not exist in the list")
}

func (e *Engine) ticker() {
	c := time.Tick(15 * time.Second)
	for {
		select {
		case <-c:
			glog.Infof(">> channel fill: %d", e.chcount)
		case <-e.stopper:
			return
		}
	}
}

func (e *Engine) runner() {
	// run the engine loop
	// TODO: numeric events?
	// TODO: REFACTOR!!!!!!!!
	for {
		select {
		case cmd := <-e.ch:
			e.chcount--
			// glog.Infof("- reading '%s' from channel with fill %d", cmd.Name, e.chcount)
			switch cmd.Name {
			case "config":
				go e.dispatchConfig(cmd)
			case "need":
				// TODO: DEBUG
				if glog.V(10) {
					pretty.Print("\nneed\n", cmd, "\n")
				}
				go e.dispatchRequest(cmd)
			case "done":
				go e.dispatchNotify(cmd)
				// case "update", "create":
			case "log":
				go e.dispatchReactions(cmd)
			case "depcheck":
				go e.dispatchCRT(cmd)
			case "error":
				go e.dispatchReactions(cmd)
				// TODO: more dispatching
			default:
				switch cmd.Destination.Type {
				case "crt":
					go e.dispatchCRT(cmd)
					go e.dispatchReactions(cmd)
				case "actio-n":
					// TODO: dispatch actions...
				default:
					glog.Warningf("do not know how to dispatch %s > %s (%T) ", cmd.Name, cmd.Source.String(), cmd.Payload)
				}

				// glog.Infof("default for %s > %s (%T) ", cmd.Name, cmd.Destination.String(), cmd.Payload)
			}
		case <-e.stopper:
			return
		}
	}
}

func (e *Engine) dispatchConfig(cmd *Command) {
	if cmd.Destination == nil {
		e.newConfig(cmd)
		return
	}
	switch cmd.Destination.Type {
	case "crt":
		if cmd.Destination.Module == "container" {
			if _, ok := e.ContainerList[cmd.Destination.ID]; ok {
				e.ContainerList[cmd.Destination.ID].Command(cmd)
				return
			}
			e.ContainerRuntime.Command(cmd)
			return
		}
		e.ContainerRuntime.Command(cmd)
		return
	default:
		glog.Warningf("unknown config destination: %s", cmd.Destination.String())
		return
	}
	return
}

func (e *Engine) dispatchCRT(cmd *Command) {
	if cmd.Destination == nil {
		for k := range e.ContainerList {
			if err := e.ContainerList[k].Command(cmd); err != nil {
				glog.Warningf("'%s': %s", k, err.Error())
			}
		}
		if err := e.ContainerRuntime.Command(cmd); err != nil {
			glog.Warningf("'%s': %s", "CRT", err.Error())
		}
		return
	}
	if cmd.Destination.Module == "container" && cmd.Destination.ID != "" {
		if _, ok := e.ContainerList[cmd.Destination.ID]; ok {
			if err := e.ContainerList[cmd.Destination.ID].Command(cmd); err != nil {
				glog.Warningf("'%s': %s", cmd.Destination.ID, err.Error())
			}
			return
		}
		if err := e.ContainerRuntime.Command(cmd); err != nil {
			glog.Warningf("'%s': %s", "CRT", err.Error())
		}
		return
	}
	// TODO: other modules
	if err := e.ContainerRuntime.Command(cmd); err != nil {
		glog.Warningf("'%s': %s", "CRT", err.Error())
	}
	return
}

func (e *Engine) dispatchReactions(cmd *Command) {
	for k := range e.Reaction {
		if glog.V(2) {
			glog.Infof("sending '%s' to '%s'", cmd.Name, k)
		}
		if err := e.Reaction[k].Command(cmd); err != nil {
			glog.Warningf("reaction %s failed: %s", k, err.Error())
			// TODO: react to errors
		}
	}
}

func (e *Engine) dispatchNotify(cmd *Command) {
	if cmd.Destination == nil {
		e.dispatchCRT(cmd)
		// TODO: react to done events?
		// e.dispatchReactions(cmd)
		return
	}
	switch cmd.Destination.Type {
	case "crt":
		e.dispatchCRT(cmd)
		return
	case "reaction":
		e.dispatchReactions(cmd)
		return
	}
}

func (e *Engine) dispatchRequest(cmd *Command) {
	// TODO: DEBUG
	if glog.V(10) {
		pretty.Print("\nneed\n", cmd, "\n")
	}
	if cmd.Destination.ID == "" {
		glog.Warning("nothing to dispatch to")
		return
	}
	switch cmd.Destination.Type {
	case "crt":
		e.dispatchCRT(cmd)
	}
}

func (e *Engine) setupContainerRuntime() error {
	// setup the container runtime
	e.ContainerList = make(map[string]Container)
	for k, v := range containerRTList {
		if e.Config[k] != nil {
			e.ContainerRuntime = v()
			if err := e.ContainerRuntime.Setup(e, e.Config[k]); err != nil {
				return err
			}
			break
		}
	}
	if e.ContainerRuntime == nil {
		return ErrNoContainerRuntime
	}
	return nil
}

func (e *Engine) setupConfStore() error {
	// setup the confstore
	if confStoreList["file"] != nil {
		e.FileStore = confStoreList["file"]()
		e.FileStore.Setup(e, e.Config["stacker"])
	}
	for k, v := range confStoreList {
		if k == "file" {
			continue
		}
		if e.Config[k] != nil {
			e.ConfStore = v()
			e.ConfStore.Setup(e, e.Config[k])
			break
		}
	}
	return nil
}

func (e *Engine) setupActions() error {
	// setup all actions
	e.Action = make(map[string]Action)
	for k, v := range actionList {
		if e.Config[k] != nil {
			e.Action[k] = v()
			e.Action[k].Setup(e, e.Config[k])
		}
	}
	return nil
}

func (e *Engine) setupReactions() error {
	// setup all reactions
	e.Reaction = make(map[string]Reaction)
	for k, v := range reactionList {
		if e.Config[k] != nil {
			e.Reaction[k] = v()
			e.Reaction[k].Setup(e, e.Config[k])
		}
	}
	return nil
}

func (e *Engine) newConfig(cmd *Command) {
	switch cmd.Payload.(type) {
	case map[string]interface{}:
		// get the config type, this is set by the confstore module.
		t, ok := cmd.Payload.(map[string]interface{})["_type"].(string)
		if !ok {
			glog.Info("no type specified")
			return
		}
		cfg, err := e.parseConfig(t, cmd.Payload.(map[string]interface{}))
		if err != nil {
			// TODO: ERROR
			return
		}
		src := &NameSpace{Type: "engine"}
		cfgcmd := &Command{
			Name:        "config",
			Source:      src,
			Destination: cfg.Target(),
			Payload:     cfg.(interface{}),
		}
		e.Event(cfgcmd)
	}
}

func (e *Engine) parseConfig(t string, data map[string]interface{}) (Config, error) {
	if configDecoderList[t] == nil {
		glog.Warningf("config entry [%s] has no config decoder", t)
		return nil, ErrNoConfigDecoder
	}
	cfg, err := configDecoderList[t].DecodeConfig(data)
	if err != nil {
		glog.Warningf("config decoder for [%s] failed: %s", t, err.Error())
		return nil, ErrConfigDecoderFailed
	}
	return cfg, nil
}

func (e *Engine) ensureDefault(name string) error {
	if _, ok := e.Config[name]; !ok {
		cfg, err := e.parseConfig(name, make(map[string]interface{}))
		if err != nil {
			glog.Fatal("can not make '%s' config: %s", name, err.Error())
		}
		e.Config[name] = cfg
	}
	return nil
}

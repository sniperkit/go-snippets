package container

import (
	"fmt"
	"strings"
	"sync"

	"github.com/cnf/structhash"
	"github.com/docker/docker/nat"
	"github.com/docker/docker/pkg/units"
	dockerapi "github.com/fsouza/go-dockerclient"
	"github.com/golang/glog"

	"github.com/cnf/stacker/docker/wrapper"
	"github.com/cnf/stacker/engine"
)

// Container represents a container
type Container struct {
	sync.Mutex
	ID   string
	Name string
	Hash string

	DockerState *dockerapi.State
	Config      *Config

	action       Action
	namespace    *engine.NameSpace
	hasImage     bool
	eng          *engine.Engine
	dependencies map[string]bool
	depcheck     bool
	Channel      chan string
}

// NewContainer makes a new Container object
func NewContainer(name string, eng *engine.Engine) *Container {
	ns := &engine.NameSpace{
		Type:   "crt",
		Module: "container",
		ID:     name,
	}
	return &Container{
		Config:       NewConfig(),
		Name:         name,
		namespace:    ns,
		eng:          eng,
		dependencies: make(map[string]bool),
	}
}

const maxLoops int = 10

// Command takes a Command object, and triggers events
func (c *Container) Command(cmd *engine.Command) error {
	count := 0
	event := cmd.Name
	// TODO: DEBUG
	if glog.V(10) {
		glog.Info("==============")
		glog.Infof("receiving event '%s' for container '%s'", event, c.Name)
	}
	c.Lock()
	// rand := md5.Sum([]byte(time.Now().String()))
	// glog.Warningf("ENTERING : %s[%s] %x", c.Name, cmd.Name, rand)
	// defer glog.Warningf("EXITING  : %s[%s] %x", c.Name, cmd.Name, rand)
	defer c.Unlock()
	if c.action == Error && event != "config" {
		// TODO: Propper error
		return fmt.Errorf("container in error, not doing anything")
	}
	switch event {
	case "need":
		if c.action == Done {
			c.notifyDone(cmd.Source)
			return nil
		}
	case "done":
		if c.action != Needing {
			return nil
		}
		if err := c.updateDependencies(cmd); err != nil {
			return err
		}
		if c.hasDependencies() {
			return nil
		}
	case "depcheck":
		if err := c.updateDependencies(cmd); err != nil {
			return err
		}
		if !c.depcheck {
			return nil
		}
		// TODO: ...
	case "update":
	case "create":
	case "start":
	case "stop":
	case "die":
		c.DockerState = nil
	case "destroy":
		c.DockerState = nil
		c.ID = ""
	case "pause", "unpause":
		glog.Warning("we ignore pause/unpause events for now")
		return nil
	case "config":
		c.newConfig(cmd)
		return nil
	default:
		glog.Warningf("unknown event `%s`", event)
		return engine.ErrUnknownEvent
	}
EventLoop:
	for {
		if count >= maxLoops {
			c.action = Error
			return engine.ErrLoopDetected
		}
		count++
		c.updateDockerState()
		c.update()
		switch c.action {
		case Done:
			if glog.V(2) {
				glog.Infof("done with container %s", c.Name)
			}
			c.notifyDone(nil)
			break EventLoop
		case Creating:
			if glog.V(1) {
				glog.Infof("creating container %s", c.Name)
			}
			if err := c.create(); err != nil {
				glog.Warningf("can not create container: %s", err.Error())
			}
		case Starting:
			if glog.V(1) {
				glog.Infof("starting container %s", c.Name)
			}
			if err := c.start(); err != nil {
				glog.Warningf("can not start container: %s", err.Error())
				// TODO: handle this cleanly
				c.action = Error
				return nil
			}
		case Stopping:
			if glog.V(1) {
				glog.Infof("stopping container %s", c.Name)
			}
			c.stop()
		case Removing:
			if glog.V(1) {
				glog.Infof("removing container %s", c.Name)
			}
			c.remove()
		case Needing:
			glog.Infof("container %s needs stuff", c.Name)
			break EventLoop
		default:
			glog.Warningf("unknown state for container %s [%d]", c.Name, c.action)
			break EventLoop
		}
	}
	// TODO: DEBUG
	if glog.V(10) {
		glog.Info("==============")
	}
	return nil
}

// Dependencies of this container
func (c *Container) Dependencies() []string {
	// TODO
	return append(c.Config.Dependencies, c.Config.VolumesFrom...)
}

func (c *Container) doesHashMatch() bool {
	hv := structhash.Version(c.Hash)
	if hv == -1 {
		return false
	}
	hash, err := structhash.Hash(c.Config, hv)
	if err != nil {
		return false
	}
	if c.Hash == hash {
		return true
	}
	return false
}

func (c *Container) isStateCurrent() bool {
	// TODO: all states
	if c.DockerState == nil && c.Config.State == "removed" {
		return true
	}
	if c.DockerState.Running == true && c.Config.State == "running" {
		return true
	}
	if c.DockerState.Running == false && c.Config.State == "stopped" {
		return true
	}
	return false
}

func (c *Container) update() {
	c.action = Done
	if c.Config == nil {
		glog.Infof("not ours to manage: '%s'", c.Name)
		return
	}

	if c.DockerState != nil && c.Config.State == "created" {
		glog.Infof("'%s' exists, not touching it", c.Name)
		return
	}

	if !c.depcheck {
		c.action = Needing
		return
	}

	if c.Config.State == "removed" {
		if c.DockerState != nil {
			if c.DockerState.Running == true {
				glog.Infof("stopping %s", c.Name)
				c.action = Stopping
				return
			}
			glog.Infof("removing %s", c.Name)
			c.action = Removing
			return
		}
		return
	}

	// TODO: Dependency checking

	if c.DockerState == nil && c.Config.State != "removed" {
		glog.Infof("creating %s", c.Name)
		c.action = Creating
		return
	}

	if c.DockerState.Running == true && c.Config.State == "stopped" {
		c.action = Stopping
		return
	}

	if c.Config.Restart == "never" {
		glog.Infof("no restart allowed: %s", c.Name)
		return
	}

	h := c.doesHashMatch()

	if !h && c.Config.Restart == "always" {
		glog.Infof("hash does not match, updating: %s", c.Name)
		if c.DockerState == nil {
			glog.Infof("creating %s", c.Name)
			c.action = Creating
			return
		}
		if c.DockerState.Running == true {
			glog.Infof("stopping %s", c.Name)
			c.action = Stopping
			return
		}
		if c.DockerState.Running == false {
			glog.Infof("removing %s", c.Name)
			c.action = Removing
			return
		}
	}

	if c.DockerState.Running == false && c.Config.State == "running" {
		glog.Infof("starting %s", c.Name)
		c.action = Starting
		return
	}

	if glog.V(2) {
		glog.Infof("container '%s' is in sync", c.Name)
	}
}

func (c *Container) create() error {
	// TODO: check for, and create if needed, the image used
	// TODO: dependencies
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	opts := c.getCreateContaineroptions()
	nc, err := dw.Client.CreateContainer(opts)
	if err != nil {
		if err == dockerapi.ErrNoSuchImage {
			c.hasImage = false
		}
		return err
	}
	c.ID = nc.ID
	c.DockerState = &nc.State
	cHash, err := structhash.Hash(c.Config, engine.LatestHashVersion)
	if err != nil {
		return err
	}
	c.Hash = cHash
	return nil
}

func (c *Container) start() error {
	// TODO: dependencies
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	if c.ID == "" {
		return fmt.Errorf("container does not exist, can not start")
	}
	opts := c.getHostConfig()
	if err := dw.Client.StartContainer(c.ID, opts); err != nil {
		return err
	}
	return c.updateDockerState()
}

func (c *Container) stop() error {
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	if c.ID == "" {
		return fmt.Errorf("container does not exist, can not stop")
	}
	if err := dw.Client.StopContainer(c.ID, 10); err != nil {
		return err
	}
	nc, err := dw.Client.InspectContainer(c.ID)
	if err != nil {
		return err
	}
	c.DockerState = &nc.State
	return nil
}

func (c *Container) remove() error {
	if c.Config.Remove == false {
		glog.Warningf("can not remove container '%s', not allowed by config", c.Name)
		return engine.ErrNotAllowed
	}
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	opts := c.getRemoveContainerOptions()
	if err := dw.Client.RemoveContainer(opts); err != nil {
		return err
	}
	c.DockerState = nil
	c.ID = ""
	return nil
}

func (c *Container) newConfig(cmd *engine.Command) {
	if glog.V(1) {
		glog.Infof("updating config for '%s'", c.Name)
	}
	// TODO: validation
	c.Config = cmd.Payload.(*Config)
	c.buildDependencies()

	c.updateDependencies(cmd)
	c.notifyNeeds()
}

func (c *Container) notifyDone(rqst *engine.NameSpace) {
	cmd := &engine.Command{
		Name:        "done",
		Source:      c.namespace,
		Destination: rqst,
		Payload:     nil,
	}
	c.eng.Event(cmd)
}

func (c *Container) notifyNeeds() {
	if !c.hasImage {
		c.needImage()
	}

	if !c.depcheck {
		c.needDepCheck()
	}

	for name, has := range c.dependencies {
		if !has {
			c.needContainer(name)
		}
	}

}

func (c *Container) buildDependencies() error {
	dlist := make(map[string]bool)
	// dlist[c.Config.Image] = false
	for i := range c.Config.Dependencies {
		dlist[c.Config.Dependencies[i]] = false
	}
	c.dependencies = dlist
	return nil
}

func (c *Container) updateDependencies(cmd *engine.Command) error {
	switch cmd.Name {
	case "config":
		c.hasImage = false
		c.depcheck = false
		c.action = Needing
		return nil
	case "depcheck":
		if _, ok := cmd.Payload.(bool); !ok {
			return engine.ErrInvalidPayload
		}
		if cmd.Payload.(bool) {
			c.depcheck = true
			return nil
		}
		c.depcheck = false
		c.action = Error
		return engine.ErrCircularDependency
	case "done":
		if cmd.Source == nil {
			glog.Errorf("got 'done' without a source")
			return nil
		}
		switch cmd.Source.Module {
		case "image":
			if cmd.Source.ID == c.Config.Image {
				c.hasImage = true
				return nil
			}
		case "container":
			for name := range c.dependencies {
				if cmd.Source.ID == name {
					// TODO: REMOVE
					glog.Info("W0000000T!")
					c.dependencies[name] = true
					return nil
				}
			}
		}
		// if cmd.Source.Module == "image" && cmd.Source.ID == c.Config.Image {
		// 	c.hasImage = true
		// 	return nil
		// }
	}
	return nil
}

func (c *Container) hasDependencies() bool {
	if !c.hasImage {
		return true
	}
	if !c.depcheck {
		return true
	}
	for i := range c.dependencies {
		if !c.dependencies[i] {
			return true
		}
	}
	return false
}

func (c *Container) getCreateContaineroptions() dockerapi.CreateContainerOptions {
	var parsedMemory int64
	var err error
	if c.Config.Memory != "" {
		parsedMemory, err = units.FromHumanSize(c.Config.Memory)
	}
	if err != nil {
		glog.Warningf("can not parse memory `%s`: %s", c.Config.Memory, err.Error())
		parsedMemory = 0
	}
	var stdin, stderr, stdout bool
	for k := range c.Config.Attach {
		if strings.ToLower(c.Config.Attach[k]) == "stdin" {
			stdin = true
		}
		if strings.ToLower(c.Config.Attach[k]) == "stdout" {
			stdout = true
		}
		if strings.ToLower(c.Config.Attach[k]) == "stderr" {
			stderr = true
		}
	}
	// TODO: actually map them
	parsedExposedPorts := make(map[dockerapi.Port]struct{})
	ports, _, err := nat.ParsePortSpecs(c.Config.Expose)
	for port := range ports {
		parsedExposedPorts[dockerapi.Port(port)] = struct{}{}
	}
	cHash, err := structhash.Hash(c.Config, engine.LatestHashVersion)
	if err != nil {
		glog.Warning(err)
	}
	parsedEnv := append(c.Config.Env, fmt.Sprintf("STACKER_CONTAINER_CFG=%s", cHash))
	parsedVolumes := make(map[string]struct{})
	for _, volume := range c.Config.Volumes {
		if strings.Contains(volume, ":") {
			continue
		}
		parsedVolumes[volume] = struct{}{}
	}
	parsedNetworkDisabled := false
	if c.Config.Net == "none" {
		parsedNetworkDisabled = true
	}
	return dockerapi.CreateContainerOptions{
		Name: c.Name,
		Config: &dockerapi.Config{
			Hostname: c.Config.Hostname,
			// Domainname: ,
			User:   c.Config.User,
			Memory: parsedMemory,
			// MemorySwap: ,
			CPUShares:    c.Config.CPUShares,
			CPUSet:       c.Config.CPUSet,
			AttachStdin:  stdin,
			AttachStdout: stdout,
			AttachStderr: stderr,
			// PortSpecs: , // Deprecated
			ExposedPorts: parsedExposedPorts,
			Tty:          c.Config.TTY,
			OpenStdin:    false,
			StdinOnce:    false,
			Env:          parsedEnv,
			Cmd:          c.Config.Cmd,
			// Dns: c.Config.DNS, // API 1.9 and below, which we do not support
			Image:   c.Config.Image,
			Volumes: parsedVolumes,
			// VolumesFrom: , // Do this at start
			WorkingDir:      c.Config.Workdir,
			Entrypoint:      c.Config.Entrypoint,
			NetworkDisabled: parsedNetworkDisabled,
		},
	}
}

func (c *Container) getHostConfig() *dockerapi.HostConfig {
	// parsedBinds := make([]string, 0)
	var parsedBinds []string
	for _, volume := range c.Config.Volumes {
		if !strings.Contains(volume, ":") {
			continue
		}
		parsedBinds = append(parsedBinds, volume)
	}
	parsedPortBindings := make(map[dockerapi.Port][]dockerapi.PortBinding)
	_, binds, err := nat.ParsePortSpecs(c.Config.Publish)
	if err != nil {
		glog.Error(err.Error())
	}
	for port, bind := range binds {
		for _, entry := range bind {
			pb := dockerapi.PortBinding{
				HostIP:   entry.HostIp,
				HostPort: entry.HostPort,
			}
			parsedPortBindings[dockerapi.Port(port)] = append(parsedPortBindings[dockerapi.Port(port)], pb)
		}

	}
	return &dockerapi.HostConfig{
		Binds: parsedBinds,
		// Device mappings // TODO:
		CapAdd:          c.Config.CapAdd,
		CapDrop:         c.Config.CapDrop,
		ContainerIDFile: c.Config.CIDFile,
		// LxcConf: ,// No support atm
		Privileged:      c.Config.Privileged,
		PortBindings:    parsedPortBindings,
		Links:           c.Config.Links,
		PublishAllPorts: c.Config.PublishAll,
		DNS:             c.Config.DNS,
		DNSSearch:       c.Config.DNSSearch,
		VolumesFrom:     c.Config.VolumesFrom,
		NetworkMode:     c.Config.Net,
		RestartPolicy:   dockerapi.NeverRestart(),
	}
}

func (c *Container) getRemoveContainerOptions() dockerapi.RemoveContainerOptions {
	return dockerapi.RemoveContainerOptions{
		ID: c.ID,
		// RemoveVolumes // No idea what this really does...
		// Make this configurable
		Force: false,
	}
}

func (c *Container) updateDockerState() error {
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		return err
	}
	if glog.V(2) {
		glog.Infof("updating docker state for container '%s'", c.Name)
	}
	if c.ID == "" {
		c.DockerState = nil
		return nil
	}
	nc, err := dw.Client.InspectContainer(c.ID)
	if err == nil {
		c.DockerState = &nc.State
		return nil
	}
	opts := dockerapi.ListContainersOptions{All: true}
	// dcl, err := dw.Client.ListContainers(opts)
	dcl, err := dw.Client.ListContainers(opts)
	if err != nil {
		return err
	}
	var idMatch bool
	for _, v := range dcl {
		if v.ID == c.ID {
			idMatch = true
			break
		}
	}
	if !idMatch {
		c.ID = ""
		c.DockerState = nil
	}
	return nil
}

func (c *Container) needImage() {
	dst := &engine.NameSpace{
		Type:   "crt",
		Module: "image",
		ID:     c.Config.Image,
	}
	cmd := &engine.Command{
		Name:        "need",
		Source:      c.namespace,
		Destination: dst,
		Payload:     nil,
	}
	c.eng.Event(cmd)
}

func (c *Container) needContainer(name string) {
	dst := &engine.NameSpace{
		Type:   "crt",
		Module: "container",
		ID:     name,
	}
	cmd := &engine.Command{
		Name:        "need",
		Source:      c.namespace,
		Destination: dst,
		Payload:     nil,
	}
	c.eng.Event(cmd)
}

func (c *Container) needDepCheck() {
	dst := &engine.NameSpace{
		Type:   "crt",
		Module: "docker",
	}
	cmd := &engine.Command{
		Name:        "depcheck",
		Source:      c.namespace,
		Destination: dst,
		Payload:     c.Config.Dependencies,
	}
	c.eng.Event(cmd)
}

package docker

import (
	"github.com/golang/glog"

	"github.com/cnf/stacker/engine"
)

var deplist map[string]*Checker

// Checker ...
type Checker struct {
	Name string
	Deps []*Checker
}

func (c *Checker) hasCircularDependency() bool {
	dstack := make([]string, 0, 10)
	return c.hasDependency(dstack[0:])
}

func (c *Checker) hasDependency(dstack []string) bool {
	for ds := range dstack {
		if dstack[ds] == c.Name {
			return true
		}
	}
	dstackNew := append(dstack, c.Name)
	for i := range c.Deps {
		// glog.Warningf("dep check for '%s' <> '%s'", c.Name, c.Deps[i].Name)
		if c.Deps[i].hasDependency(dstackNew[0:]) {
			return true
		}
	}
	return false
}

func (d *Docker) checkDependencies(cmd *engine.Command) {
	if deplist == nil {
		deplist = make(map[string]*Checker)
	}
	if _, ok := cmd.Payload.([]string); !ok {
		// TODO:...
		return
	}
	if cmd.Source == nil {
		return
	}
	name := cmd.Source.ID
	if _, ok := deplist[name]; !ok {
		deplist[name] = &Checker{Name: name}
	}

	// newlist := make([]*Checker, 0)
	var newlist []*Checker
	for _, v := range cmd.Payload.([]string) {
		if val, ok := deplist[v]; ok {
			newlist = append(newlist, val)
		} else {
			newdep := &Checker{Name: v}
			deplist[v] = newdep
			newlist = append(newlist, newdep)
		}
	}
	deplist[name].Deps = newlist
	if deplist[name].hasCircularDependency() {
		glog.Warningf("!! '%s' has circular dependencies", name)
		d.notifyDepcheckNotOK(name)
		return
	}
	d.notifyDepcheckOK(name)

}

func (d *Docker) notifyDepcheckOK(name string) {
	d.notifyDepcheck(name, true)
}

func (d *Docker) notifyDepcheckNotOK(name string) {
	d.notifyDepcheck(name, false)
}

func (d *Docker) notifyDepcheck(name string, ok bool) {
	dst := &engine.NameSpace{
		Type:   "crt",
		Module: "container",
		ID:     name,
	}
	cmd := &engine.Command{
		Name: "depcheck",
		// Source:      c.namespace,
		Source:      d.ns,
		Destination: dst,
		Payload:     ok,
	}
	d.eng.Event(cmd)
}

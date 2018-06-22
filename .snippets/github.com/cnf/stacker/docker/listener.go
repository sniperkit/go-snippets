package docker

import (
	"net"
	"strconv"
	"strings"

	"github.com/golang/glog"

	"github.com/cnf/stacker/docker/wrapper"
	"github.com/cnf/stacker/engine"
)

// listener listenes to docker events on the docker api,
// and sends them to the engine.
func (d *Docker) listener() {
	defer glog.Warning("docker listener exiting")
	src := &engine.NameSpace{
		Type:   "crt",
		Module: "api",
	}
	// Start listening to Docker events
	for {
		select {
		case event := <-dockerwrapper.Channel:
			switch event.Status {
			case "untag", "delete":
				// TODO: this doesn't work...
				cd, err := getContainerData(event.ID)
				if err != nil {
					continue
					// TODO: OOPS
				}

				cName := strings.TrimLeft(cd.Name, "/")
				dst := &engine.NameSpace{
					Type:   "crt",
					Module: "image",
					ID:     cName,
				}
				cmd := &engine.Command{
					Name:        event.Status,
					Source:      src,
					Destination: dst,
					Payload:     cd,
				}
				d.eng.Event(cmd)
			case "update", "create", "start", "stop", "die", "destroy", "pause", "unpause":
				cd, err := getContainerData(event.ID)
				var cName string
				if err != nil {
					cName = d.Containers.GetNameFromID(event.ID)
				} else {
					cName = strings.TrimLeft(cd.Name, "/")
				}

				dst := &engine.NameSpace{
					Type:   "crt",
					Module: "container",
					ID:     cName,
				}
				cmd := &engine.Command{
					Name:        event.Status,
					Source:      src,
					Destination: dst,
					Payload:     cd,
				}
				d.eng.Event(cmd)
			default:
				continue
			}
		case <-d.stopper:
			return
		}
	}
}

func getContainerData(ID string) (*engine.ContainerData, error) {
	// TODO: error handeling
	dw, err := dockerwrapper.GetObject()
	if err != nil {
		glog.Errorf("could not get docker connection: %#v", err.Error())
		return nil, err
	}
	nc, err := dw.Client.InspectContainer(ID)
	if err != nil {
		return nil, err
	}
	cd := engine.NewContainerData()

	cd.ID = nc.ID
	cd.Name = strings.TrimLeft(nc.Name, "/")
	cd.Created = nc.Created

	cd.Network.IPAddress = nc.NetworkSettings.IPAddress
	cd.Network.Mask = nc.NetworkSettings.IPPrefixLen
	cd.Network.Hostname = nc.Config.Hostname
	cd.Network.Domainname = nc.Config.Domainname
	cd.Network.DNS = nc.Config.DNS
	for k, v := range nc.NetworkSettings.Ports {
		port := engine.Port(k)
		cd.Network.Ports = append(cd.Network.Ports, port)
		if v != nil {
			for _, publish := range v {
				pport, err := strconv.ParseUint(publish.HostPort, 10, 16)
				if err != nil {
					continue
				}
				pub := engine.Published{
					IP:   net.ParseIP(publish.HostIP),
					Port: uint16(pport),
				}
				cd.Network.Published[port] = append(cd.Network.Published[port], pub)
			}
		}
	}

	return cd, nil
}

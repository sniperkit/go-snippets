package hyperdb

import "time"

// datasetList maps dataset names to the pointer to a dataset
type datasetList map[string]*dataset

// Dataset contains information about which write and read servers are
// available for the dataset
type dataset struct {
	write                *Server
	read                 datacenterList // dc -> prio -> []*Server
	allReadServers       serverList     // prio -> []*Server
	lastReadServerCached time.Time
	lastReadServer       *Server
}

// addServer adds a server to the dataset (read and/or write list as
// appropriate)
func (d *dataset) addServer(server Server) {
	if server.WritePriority > 0 {
		d.write = &server
	}
	if server.ReadPriority > 0 {
		d.read.addServer(&server)
		if _, ok := d.allReadServers[server.ReadPriority]; !ok {
			d.allReadServers[server.ReadPriority] = []*Server{}
		}
		d.allReadServers[server.ReadPriority] = append(d.allReadServers[server.ReadPriority], &server)
	}
}

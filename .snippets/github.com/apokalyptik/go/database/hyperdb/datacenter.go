package hyperdb

// datasetlist is a map of datacenter to a serverList
type datacenterList map[string]serverList

// addServer adds a server into the datacenterList
func (d datacenterList) addServer(server *Server) {
	if _, ok := d[server.Datacenter]; !ok {
		d[server.Datacenter] = make(serverList)
	}
	dc := d[server.Datacenter]
	if _, ok := dc[server.ReadPriority]; !ok {
		dc[server.ReadPriority] = []*Server{}
	}
	dc[server.ReadPriority] = append(dc[server.ReadPriority], server)
}

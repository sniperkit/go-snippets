package discovery

// BrokerInfo is used to register broker into etcd
type BrokerInfo struct {
	HostPort  string
	Type      string
	Hash      string
	StartTime int64
}

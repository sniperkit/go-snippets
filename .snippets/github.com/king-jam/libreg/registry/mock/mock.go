package mock

import (
	"errors"
	"sync"

	"github.com/king-jam/libreg"
	"github.com/king-jam/libreg/registry"
)

type Mock struct {
	Catalog map[*registry.CatalogRegistration]bool
	sync.RWMutex
}

func Register() {
	libreg.AddRegistry(registry.MOCK, New)
}

func New(endpoints []string, options *registry.Config) (registry.Registry, error) {
	m := &Mock{}
	m.Catalog = make(map[*registry.CatalogRegistration]bool)
	return m, nil
}

// GetSearchTerms returns the whitelist of SSDP URNs to act on
func (m *Mock) GetCatalog() map[registry.CatalogRegistration]bool {
	c := make(map[registry.CatalogRegistration]bool)

	m.RLock()
	for r := range m.Catalog {
		c[*r] = true
	}
	m.RUnlock()

	return c
}

// Register creates a new node, service or check
func (m *Mock) Register(reg *registry.CatalogRegistration, options *registry.WriteOptions) error {
	if reg.Address == "1.1.1.1" {
		return errors.New("Forced error: reg.Address=1.1.1.1:1")
	}
	m.Lock()
	m.Catalog[reg] = true
	m.Unlock()
	return nil
}

// Deregister removes a node, service or check
func (m *Mock) Deregister(dereg *registry.CatalogDeregistration, options *registry.WriteOptions) error {
	m.Lock()
	for reg := range m.Catalog {
		if dereg.Node == reg.Node &&
			dereg.Address == reg.Address &&
			dereg.Datacenter == reg.Datacenter &&
			dereg.ServiceID == reg.Service.ID &&
			dereg.CheckID == reg.Check.CheckID {
			delete(m.Catalog, reg)
		}
	}
	m.Unlock()

	return nil
}

// Datacenters lists known datacenters
func (m *Mock) Datacenters() ([]string, error) {
	d := make(map[string]bool)
	for reg := range m.GetCatalog() {
		d[reg.Datacenter] = true
	}

	var keys []string
	for k := range d {
		keys = append(keys, k)
	}
	return keys, nil
}

// Nodes lists all nodes in a given DC
func (m *Mock) Nodes(options *registry.QueryOptions) ([]*registry.Node, error) {
	nodes := make(map[registry.Node]bool)
	for reg := range m.GetCatalog() {
		n := registry.Node{
			Node:    reg.Node,
			Address: reg.Address,
		}
		nodes[n] = true
	}

	var keys []*registry.Node
	for k := range nodes {
		keys = append(keys, &k)
	}
	return keys, nil
}

// Services lists all services in a given DC
func (m *Mock) Services(options *registry.QueryOptions) (map[string][]string, error) {
	serviceMap := make(map[string][]string)

	catalog := m.GetCatalog()
	for r := range catalog {
		serviceMap[r.Service.Service] = append(serviceMap[r.Service.Service], r.Service.Tags...)
		if serviceMap[r.Service.Service] == nil {
			s := []string{}
			serviceMap[r.Service.Service] = s
		}
	}

	return serviceMap, nil
}

// Service lists the nodes in a given service
func (m *Mock) Service(service, tag string, options *registry.QueryOptions) ([]*registry.CatalogService, error) {
	var c []*registry.CatalogService

	for r := range m.GetCatalog() {
		containsTag := false

		for _, serviceTag := range r.Service.Tags {
			if tag == serviceTag {
				containsTag = true
				break
			}
		}
		if tag == "" {
			containsTag = true
		}

		if r.Service.Service == service &&
			containsTag == true {
			catalogService := registry.CatalogService{
				Node:                     r.Node,
				Address:                  r.Address,
				ServiceID:                r.Service.ID,
				ServiceName:              r.Service.Service,
				ServiceAddress:           r.Service.Address,
				ServiceTags:              r.Service.Tags,
				ServicePort:              r.Service.Port,
				ServiceEnableTagOverride: r.Service.EnableTagOverride,
			}
			c = append(c, &catalogService)
		}
	}
	return c, nil
}

// Node lists the services provided by a given node
func (m *Mock) Node(node string, options *registry.QueryOptions) (*registry.CatalogNode, error) {
	if node == "" {
		return nil, errors.New("Node UUID is empty")
	}

	catalogNode := registry.CatalogNode{
		Node:     &registry.Node{},
		Services: make(map[string]*registry.AgentService),
	}

	for c := range m.GetCatalog() {
		if c.Node == node {
			catalogNode.Node.Node = c.Node
			catalogNode.Node.Address = c.Address
			catalogNode.Services[c.Service.Service] = c.Service
		}
	}
	return &catalogNode, nil
}

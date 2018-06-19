package sample

import (
	"github.com/CenturyLinkLabs/pmxadapter"
)

type SampleAdapter struct {
}

// Implementation of the PanamaxAdapter GetServices interface
func (a *SampleAdapter) GetServices() ([]pmxadapter.ServiceDeployment, error) {
	return []pmxadapter.ServiceDeployment{{}}, nil
}

// Implementation of the PanamaxAdapter GetService interface
func (a *SampleAdapter) GetService(id string) (pmxadapter.ServiceDeployment, error) {
	return pmxadapter.ServiceDeployment{}, nil
}

// Implementation of the PanamaxAdapter CreateServices interface
func (a *SampleAdapter) CreateServices(services []*pmxadapter.Service) ([]pmxadapter.ServiceDeployment, error) {
	return []pmxadapter.ServiceDeployment{{}}, nil
}

// Implementation of the PanamaxAdapter UpdateService interface
func (a *SampleAdapter) UpdateService(s *pmxadapter.Service) error {
	return nil
}

// Implementation of the PanamaxAdapter DestroyService interface
func (a *SampleAdapter) DestroyService(id string) error {
	return nil
}

// GetMetadata returns metadata for the adapter.
func (a *SampleAdapter) GetMetadata() pmxadapter.Metadata {
	return pmxadapter.Metadata{Type: "Sample", Version: "0.1"}
}

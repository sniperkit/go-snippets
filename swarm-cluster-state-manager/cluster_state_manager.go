// Copyright 2015 Google, Inc All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/samalba/dockerclient"
)

// ClusterState holds the Swarm cluster state declaration.
type ClusterState struct {
	Count int
	Image string
	Name  string

	// filters specifies an JSON encoded value of the filters used when
	// processing containers list via the Docker API.
	filters string
}

// ClusterStateStatus holds the Swarm cluster status status.
type ClusterStateStatus struct {
	Containers   []dockerclient.Container
	CurrentCount int
	DesiredCount int
	Image        string
	Name         string
}

// A ClusterStateManager defines parameters for running a Swarm cluster state manager server.
type ClusterStateManager struct {
	// DockerClient specifies the Docker client by which individual
	// Swarm requests are made.
	DockerClient *dockerclient.DockerClient

	// Store specifies the in-memory data store for cluster state by name.
	Store map[string]*ClusterState
	sync.RWMutex
}

func newClusterStateManager(daemonUrl string, tlsConfig *tls.Config) (*ClusterStateManager, error) {
	store := make(map[string]*ClusterState)
	client, err := dockerclient.NewDockerClient(daemonUrl, tlsConfig)
	if err != nil {
		return nil, err
	}
	return &ClusterStateManager{DockerClient: client, Store: store}, nil
}

func (csm *ClusterStateManager) Submit(cs *ClusterState) error {
	csm.Lock()
	defer csm.Unlock()

	log.Printf("updating %s cluster state", cs.Name)
	filters, err := json.Marshal(map[string][]string{
		"label": []string{fmt.Sprintf("swarm.cluster.state=%s", cs.Name)},
	})
	if err != nil {
		return err
	}
	cs.filters = string(filters)
	csm.Store[cs.Name] = cs
	return nil
}

func (csm *ClusterStateManager) Remove(name string) {
	csm.Lock()
	defer csm.Unlock()
	log.Printf("removing %s cluster state", name)
	delete(csm.Store, name)
}

func (csm *ClusterStateManager) ClusterStatus() ([]*ClusterStateStatus, error) {
	csm.RLock()
	defer csm.RUnlock()

	cs := make([]*ClusterStateStatus, 0)
	for _, desiredState := range csm.Store {
		containers, err := csm.DockerClient.ListContainers(false, false, desiredState.filters)
		if err != nil {
			log.Println(err)
			continue
		}

		cs = append(cs, &ClusterStateStatus{
			Containers:   containers,
			CurrentCount: len(containers),
			DesiredCount: desiredState.Count,
			Image:        desiredState.Image,
			Name:         desiredState.Name,
		})
	}

	return cs, nil
}

func (csm *ClusterStateManager) Sync() {
	for {
		clusterStatus, err := csm.ClusterStatus()
		if err != nil {
			log.Print(err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, status := range clusterStatus {
			delta := status.DesiredCount - status.CurrentCount

			if delta == 0 {
				continue
			}

			if delta > 0 {
				log.Printf("creating %d %s containers", delta, status.Name)

				labels := make(map[string]string)
				labels["swarm.cluster.state"] = status.Name

				for i := 0; i < delta; i++ {
					c := dockerclient.ContainerConfig{
						Image:  status.Image,
						Labels: labels,
					}
					id, err := csm.DockerClient.CreateContainer(&c, "")
					if err != nil {
						log.Print(err)
						continue
					}
					err = csm.DockerClient.StartContainer(id, nil)
					if err != nil {
						log.Print(err)
						continue
					}
					log.Printf("created %s container: %s", status.Name, id)
				}

				continue
			}

			log.Printf("removing %d %s containers", -delta, status.Name)
			for _, container := range status.Containers[status.DesiredCount:] {
				err = csm.DockerClient.RemoveContainer(container.Id, true, true)
				if err != nil {
					log.Print(err)
					continue
				}
				log.Printf("removed %s container: %s", status.Name, container.Id)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

package lib

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"registry-factory/client"
	"time"
)

//Executor ...
type Executor struct {
	hostOn    string
	docker    *client.DockerClient
	harbor    string
	namespace string
}

//Environment ...
type Environment struct {
	Target    ProxyTarget
	RuntimeID string
}

//NewExecutor ...
func NewExecutor(dockerdHost string, hPort uint, harbor string) *Executor {
	dHost := ""
	if hPort > 0 {
		dHost = fmt.Sprintf("tcp://%s:%d", dockerdHost, hPort)
	}
	return &Executor{
		hostOn: dockerdHost,
		docker: &client.DockerClient{
			Host: dHost,
		},
		harbor:    harbor,
	}
}

//SetNamespace ...
func (e *Executor) SetNamespace(ns string){
	if len(ns) > 0 {
		e.namespace = ns
	}
}

//Exec ...
func (e *Executor) Exec(policy *SchedulePolicy) (Environment, error) {
	if len(policy.Image) == 0 {
		return Environment{}, errors.New("empty image")
	}

	if len(policy.Tag) == 0 {
		policy.Tag = "latest"
	}

	//Only keep the 1st port as target port
	bindPorts := []string{}
	targetPort := 0
	for _, port := range policy.BoundPorts {
		portOnHost := giveMePort()
		if targetPort == 0 {
			targetPort = (int)(portOnHost)
		}
		boundPort := fmt.Sprintf("%d:%d", portOnHost, port)
		bindPorts = append(bindPorts, boundPort)
	}

	image := fmt.Sprintf("%s:%s", policy.Image, policy.Tag)
	if !policy.UseHub {
		image = fmt.Sprintf("%s/%s/%s", e.harbor, e.namespace, image)
	}

	runID, err := e.docker.Run(image, "", "", true, true, bindPorts, policy.EnvVars)
	if err != nil {
		return Environment{}, err
	}

	//Check connection available
	done := make(chan bool)
	errCh := make(chan error)
	go func() {
		t := time.Tick(1 * time.Second)
		for {
			select {
			case <-t:
				if res, err := http.Get(fmt.Sprintf("http://%s:%d", e.hostOn, targetPort)); err == nil {
					log.Printf("Checking connection of http://%s:%d: %s\n", e.hostOn, targetPort, res.Status)
					if res.StatusCode == http.StatusOK {
						time.Sleep(2 * time.Second)
						done <- true
						return
					}
				}
			case <-time.After(10 * time.Second):
				errCh <- fmt.Errorf("Check connection %s:%d timeout", e.hostOn, targetPort)
				return
			}
		}
	}()

	select {
	case <-done:
		return Environment{
			Target:    (ProxyTarget)(fmt.Sprintf("%s:%d", e.hostOn, targetPort)),
			RuntimeID: runID,
		}, nil
	case err := <-errCh:
		return Environment{}, err
	}
}

//Destroy ...
func (e *Executor) Destroy(runtimeID string) error {
	if len(runtimeID) == 0 {
		return errors.New("nil runtime ID")
	}

	return e.docker.Destroy(runtimeID)
}

func giveMePort() int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return 30000 + r.Int31n(35530)
}

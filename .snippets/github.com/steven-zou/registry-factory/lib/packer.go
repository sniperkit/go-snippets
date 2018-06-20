package lib

import (
	"errors"
	"fmt"
	"log"
	"registry-factory/client"
)

//Packer ...
type Packer struct {
	hostOn    string
	docker    *client.DockerClient
	harbor    string
	namespace string
}

//NewPacker ...
func NewPacker(dockerdHost string, dockerdPort uint, harborHost string) *Packer {
	dHost := ""
	if dockerdPort > 0 {
		dHost = fmt.Sprintf("tcp://%s:%d", dockerdHost, dockerdPort)
	}

	return &Packer{
		hostOn: dockerdHost,
		docker: &client.DockerClient{
			Host: dHost,
		},
		harbor: harborHost,
	}
}

//SetNamespace ...
func (p *Packer) SetNamespace(ns string) {
	if len(ns) > 0 {
		p.namespace = ns
	}
}

//Build ...
func (p *Packer) Build(baseContainer string, image, tag string) error {
	if len(baseContainer) == 0 {
		return errors.New("empty base container")
	}

	newTag := tag
	if len(newTag) == 0 {
		newTag = "latest"
	}

	fullNamespace := fmt.Sprintf("%s/%s/%s", p.harbor, p.namespace, image)
	if err := p.docker.Commit(baseContainer, fullNamespace, newTag); err != nil {
		return err
	}

	//login
	if err := p.docker.Login(Config.Dockerd.Admin, Config.Dockerd.Password, p.harbor); err != nil {
		return err
	}
	backendImage := fmt.Sprintf("%s:%s", fullNamespace, newTag)
	if err := p.docker.Push(backendImage); err != nil {
		return err
	}

	//Just try to remove local image
	if err := p.docker.RMImage(backendImage); err != nil {
		log.Printf("rm image error: %s\n", err)
	}

	return nil
}

//BuildLocal ...
func (p *Packer) BuildLocal(baseContainer string, image, tag string) error {
	if len(baseContainer) == 0 {
		return errors.New("empty base container")
	}

	newTag := tag
	if len(newTag) == 0 {
		newTag = "latest"
	}
	return p.docker.Commit(baseContainer, image, newTag)
}

//RMImage remove the specified image
func (p *Packer) RMImage(image string) error {
	if len(image) == 0 {
		return errors.New("empty image")
	}

	return p.docker.RMImage(image)
}

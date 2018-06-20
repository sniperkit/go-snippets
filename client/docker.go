package client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"log"
)

const (
	dockerCmd = "docker"
)

//DockerClient : Run docker commands
type DockerClient struct {
	//The host sock of docker listening: unix:///var/docker
	Host string
}

//Status : Check if docker daemon is there
func (dc *DockerClient) Status() error {
	args := []string{"version"}

	return dc.runCommand(dockerCmd, dc.arguments(args))
}

//Pull : Pull image
func (dc *DockerClient) Pull(image string) error {
	if len(strings.TrimSpace(image)) == 0 {
		return errors.New("Empty image")
	}

	args := []string{"pull", image}

	return dc.runCommandWithOutputs(dockerCmd, dc.arguments(args))
}

//Tag :Tag image
func (dc *DockerClient) Tag(source, target string) error {
	if len(strings.TrimSpace(source)) == 0 ||
		len(strings.TrimSpace(target)) == 0 {
		return errors.New("Empty images")
	}

	args := []string{"tag", source, target}

	return dc.runCommandWithOutputs(dockerCmd, dc.arguments(args))
}

//Push : push image
func (dc *DockerClient) Push(image string) error {
	if len(strings.TrimSpace(image)) == 0 {
		return errors.New("Empty image")
	}

	args := []string{"push", image}

	return dc.runCommandWithOutputs(dockerCmd, dc.arguments(args))
}

//Login : Login docker
func (dc *DockerClient) Login(userName, password string, uri string) error {
	if len(strings.TrimSpace(userName)) == 0 ||
		len(strings.TrimSpace(password)) == 0 {
		return errors.New("Invalid credential")
	}

	args := []string{"login", "-u", userName, "-p", password, uri}

	return dc.runCommandWithOutputs(dockerCmd, dc.arguments(args))
}

//Run containers
func (dc *DockerClient) Run(image, name, cmd string, isInteractive, asDaemon bool, bindPorts []string, env map[string]string) (string, error) {
	if len(strings.TrimSpace(image)) == 0 {
		return "", errors.New("image must be specified")
	}

	args := []string{"run"}
	containerName := name
	if len(strings.TrimSpace(containerName)) == 0 {
		containerName = fmt.Sprintf("container-%d", time.Now().UnixNano())
	}
	args = append(args, "--name", containerName)
	if isInteractive {
		args = append(args, "-it")
	}

	//[]string{"5674:5674", "8080:80"}
	if len(bindPorts) > 0 {
		for _, portMapping := range bindPorts {
			args = append(args, "-p", portMapping)
		}
	}

	if asDaemon {
		args = append(args, "-d")
	}

	for k, v := range env {
		//TODO: not support escaping quotes.
		envStr := fmt.Sprintf(`%s=%s`, k, v)
		args = append(args, "-e", envStr)
	}

	args = append(args, image)
	if len(strings.TrimSpace(cmd)) > 0 {
		args = append(args, cmd)
	}

	return dc.runCommandWithOutput2(dockerCmd, dc.arguments(args))
}

//Destroy container
func (dc *DockerClient) Destroy(container string) error {
	if len(strings.TrimSpace(container)) == 0 {
		return errors.New("empty container")
	}

	args := []string{"rm", "-f", container}
	return dc.runCommandWithOutputs(dockerCmd, dc.arguments(args))
}

//Commit ...
func (dc *DockerClient) Commit(container string, image, tag string) error {
	if len(strings.TrimSpace(container)) == 0 {
		return errors.New("empty container")
	}

	if len(image) == 0 {
		return errors.New("empty image name")
	}

	newTag := tag
	if len(newTag) == 0 {
		newTag = "latest"
	}

	fullNS := fmt.Sprintf("%s:%s", image, tag)
	args := []string{"commit", container, fullNS}

	return dc.runCommandWithOutputs(dockerCmd, dc.arguments(args))
}

//RMImage ...
func (dc *DockerClient) RMImage(image string) error {
	if len(strings.TrimSpace(image)) == 0 {
		return errors.New("empty image name")
	}

	args := []string{"rmi", "-f", image}
	return dc.runCommandWithOutputs(dockerCmd, dc.arguments(args))
}

func (dc *DockerClient) runCommandWithOutput2(cmdName string, args []string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	log.Printf("command: %s\n", cmd.Args)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	if err := cmd.Run(); err != nil {
		return "", err
	}

	output := string(cmdOutput.Bytes())
	output = strings.TrimRight(output, "\n")

	return output, nil
}

func (dc *DockerClient) runCommand(cmdName string, args []string) error {
	return exec.Command(cmdName, args...).Run()
}

func (dc *DockerClient) runCommandWithOutput(cmdName string, args []string) error {
	cmd := exec.Command(cmdName, args...)
	log.Printf("command: %s\n", cmd.Args)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Printf("%s out | %s\n", cmdName, scanner.Text())
		}
	}()

	if err = cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func (dc *DockerClient) runCommandWithOutputs(cmdName string, args []string) error {
	cmd := exec.Command(cmdName, args...)
	log.Printf("command: %s\n", cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err
	}
	if len(output) > 0 {
		log.Printf("[%s] OUT: %s", cmdName, output)
	}

	errData, _ := ioutil.ReadAll(stderr)
	if len(errData) > 0 {
		log.Printf("[%s] ERROR: %s", cmdName, errData)
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (dc *DockerClient) arguments(args []string) []string {
	argList := []string{}
	if len(strings.TrimSpace(dc.Host)) > 0 {
		argList = append(argList, fmt.Sprintf("-H %s", dc.Host))
	}

	if len(args) > 0 {
		argList = append(argList, args...)
	}

	return argList
}

package main

import (
	"log"

	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/docker/docker/builder"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

func main() {

	//TODO repository name
	imageRepositoryName := "dwrap-image"

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	cmdName := os.Args[1]
	cmdArgs := []string{}
	if len(os.Args) > 2 {
		cmdArgs = os.Args[2:]
	}

	imageFullName := fmt.Sprintf("%s/%s", imageRepositoryName, cmdName)

	host := os.Getenv("DOCKER_HOST")
	if host == "" {
		os.Setenv("DOCKER_HOST", "unix:///var/run/docker.sock")
	}
	cl, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// search image
	images, err := cl.ImageList(ctx, types.ImageListOptions{
		MatchName: imageFullName,
	})

	if err != nil {
		log.Fatal(err)
	}
	if images == nil || len(images) == 0 {
		buildImage(ctx, cl, imageFullName, cmdName)
	}

	containerConfig := &container.Config{
		Image: imageFullName,
		Cmd:   cmdArgs,
		//Env:          nil,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Tty:          false,
		OpenStdin:    true,
		StdinOnce:    true,
	}
	hostConfig := &container.HostConfig{
		AutoRemove: true,
	}
	resp, err := cl.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
	if err != nil {
		log.Fatal(err)
	}
	containerID := resp.ID
	defer func() {
		status, err := cl.ContainerWait(ctx, containerID)
		if err != nil {
			log.Fatal(err)
		}
		if status == 0 {
			options := types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			}
			if err := cl.ContainerRemove(ctx, containerID, options); err != nil {
				log.Fatal(err)
			}
		}
	}()

	attachOptions := types.ContainerAttachOptions{
		Stdout: true,
		Stderr: true,
		Stdin:  true,
		Stream: true,
	}
	attachResp, err := cl.ContainerAttach(ctx, containerID, attachOptions)
	if err != nil {
		log.Fatal(err)
	}

	containerStartOptions := types.ContainerStartOptions{}
	if err := cl.ContainerStart(ctx, containerID, containerStartOptions); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := holdHijackedConnection(os.Stdin, os.Stdout, os.Stderr, attachResp); err != nil {
			log.Fatal(err)
		}
	}()
}

func printUsage() {
	usage := `
Run command on docker with generated with image based by the Alipne Linux.

Usage:
  dwrap [options] [COMMAND] [COMMAND ARGS...]

Options:
  No option supported in the current version.

Examples:
  # download html by curl.
  dwrap curl -L http://example.com/index.html > index.html

  # download html by wget.
  dwrap wget -O - http://example.com/index.html > index.html

  # echo json and print by jq command.
  echo any.json | dwrap jq "."
`
	fmt.Println(usage)
}

func buildImage(ctx context.Context, cl *client.Client, imageFullName string, cmdName string) {
	buildCtx, relDockerfile, err := builder.GetContextFromReader(createDockerFileText(cmdName), "")
	if err != nil {
		log.Fatal(err)
	}

	res, err := cl.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Tags:       []string{imageFullName},
		Dockerfile: relDockerfile,
	})
	if err != nil {
		log.Fatal(err)

	}

	defer res.Body.Close()
	err = sleepWhileBuild(ctx, cl, imageFullName, 60*time.Second) // TODO parameterize
	if err != nil {
		log.Fatal(err)
	}
}

func sleepWhileBuild(ctx context.Context, cl *client.Client, imageFullName string, timeout time.Duration) error {

	current := 0 * time.Second
	interval := 1 * time.Second
	for {
		images, err := cl.ImageList(ctx, types.ImageListOptions{
			MatchName: imageFullName,
		})
		if err != nil {
			log.Fatal(err)
		}
		if images != nil && len(images) > 0 {
			return nil
		}

		time.Sleep(interval)
		current += interval

		if timeout > 0 && current > timeout {
			return fmt.Errorf("Timeout: sleepWhileBuild:%s", imageFullName)
		}
	}
}

func createDockerFileText(cmdName string) io.ReadCloser {
	tmpl, err := template.New("dockerfile").Parse(dockerFileTemplate)
	if err != nil {
		log.Fatal(err)
	}

	p := make(map[string]string)
	p["cmd"] = cmdName

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, p); err != nil {
		log.Fatal(err)
	}
	return ioutil.NopCloser(strings.NewReader(buf.String()))
}

// holdHijackedConnection handles copying input to and output from streams to
// the connection. Copied from github.com/docker/docker/api/client.
func holdHijackedConnection(inputStream io.ReadCloser, outputStream, errorStream io.Writer, resp types.HijackedResponse) error {
	var err error

	receiveStdout := make(chan error, 1)
	if outputStream != nil || errorStream != nil {
		go func() {
			_, err = stdcopy.StdCopy(outputStream, errorStream, resp.Reader)
			//log.Printf("[hijack] End of stdout")
			receiveStdout <- err
		}()
	}

	stdinDone := make(chan struct{})
	go func() {
		if inputStream != nil {
			io.Copy(resp.Conn, inputStream)
			//log.Printf("[hijack] End of stdin")
		}

		if err := resp.CloseWrite(); err != nil {
			log.Printf("dwrap: couldn't send EOF: %s", err)
		}
		close(stdinDone)
	}()

	select {
	case err := <-receiveStdout:
		if err != nil {
			return fmt.Errorf("Error receiveStdout: %s", err)
		}
	case <-stdinDone:
		if outputStream != nil || errorStream != nil {
			err := <-receiveStdout
			if err != nil {
				return fmt.Errorf("Error receiveStdout: %s", err)
			}
		}
	}

	return nil
}

const dockerFileTemplate = `
FROM alpine:latest
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>

RUN set -x && if [ ! $(which {{.cmd}}) ]; then \
    apk add --no-cache {{.cmd}}; \
    fi

ENTRYPOINT ["{{.cmd}}"]
`

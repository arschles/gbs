package handlers

import (
	"bytes"
	"fmt"
	"io"

	"code.google.com/p/go-uuid/uuid"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	golangImage = "golang:1.5.2"
	pwd         = "pwd"
	absPwd      = "/" + pwd
)

func createContainer(cl *docker.Client, workdir, site, org, repo string) (*docker.Container, error) {
	return cl.CreateContainer(docker.CreateContainerOptions{
		Name: fmt.Sprintf("build-%s-%s-%s-%s", site, org, repo, uuid.New()),
		Config: &docker.Config{
			Env:   []string{"GO15VENDOREXPERIMENT=1", "CGO_ENABLED=0", "SITE=" + site, "ORG=" + org, "REPO=" + repo},
			Cmd:   []string{"/bin/bash", absPwd + "/build.sh"},
			Image: golangImage,
			Volumes: map[string]struct{}{
				workdir: struct{}{},
			},
			Mounts: []docker.Mount{
				docker.Mount{Name: "pwd", Source: workdir, Destination: absPwd, Mode: "rx"},
			},
		},
		HostConfig: &docker.HostConfig{},
	})
}

func startContainer(cl *docker.Client, con *docker.Container, workdir string) error {
	return cl.StartContainer(con.ID, &docker.HostConfig{
		Binds: []string{fmt.Sprintf("%s:%s", workdir, absPwd)},
	})
}

func attachContainer(cl *docker.Client, con *docker.Container) (io.Reader, io.Reader, error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	opts := docker.AttachToContainerOptions{
		Container:    con.ID,
		OutputStream: &stdoutBuf,
		ErrorStream:  &stderrBuf,
		Logs:         true,
		Stream:       true,
		Stdout:       true,
		Stderr:       true,
	}
	if err := cl.AttachToContainer(opts); err != nil {
		return nil, nil, err
	}

	return &stdoutBuf, &stderrBuf, nil
}

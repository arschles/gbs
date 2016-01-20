package handlers

import (
	"fmt"
	"io"

	"code.google.com/p/go-uuid/uuid"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	pwd    = "pwd"
	absPwd = "/" + pwd
)

func createContainerOpts(img, workdir, site, org, repo string, env ...string) docker.CreateContainerOptions {
	defaultEnv := []string{"GO15VENDOREXPERIMENT=1", "SITE=" + site, "ORG=" + org, "REPO=" + repo}
	return docker.CreateContainerOptions{
		Name: fmt.Sprintf("build-%s-%s-%s-%s", site, org, repo, uuid.New()),
		Config: &docker.Config{
			Env:   append(defaultEnv, env...),
			Image: img,
			Volumes: map[string]struct{}{
				workdir: struct{}{},
			},
			Mounts: []docker.Mount{
				docker.Mount{Name: "pwd", Source: workdir, Destination: absPwd, Mode: "rx"},
			},
		},
		HostConfig: &docker.HostConfig{},
	}
}

// attachContainerOpts returns docker.AttachToContainerOptions with output and error streams turned on
// as well as logs. the returned io.Reader will output both stdout and stderr
func attachToContainerOpts(containerID string) (docker.AttachToContainerOptions, io.Reader) {
	r, w := io.Pipe()
	// var stdoutBuf, stderrBuf bytes.Buffer
	opts := docker.AttachToContainerOptions{
		Container:    containerID,
		OutputStream: w,
		ErrorStream:  w,
		Logs:         true,
		Stream:       true,
		Stdout:       true,
		Stderr:       true,
	}

	return opts, r
}

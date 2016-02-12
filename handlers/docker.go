package handlers

import (
	"fmt"
	"io"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/pborman/uuid"
)

func containerName(site, org, repo string) string {
	return fmt.Sprintf("build-%s-%s-%s-%s", site, org, repo, uuid.New())
}

func packageName(site, org, repo string) string {
	return fmt.Sprintf("%s/%s/%s", site, org, repo)
}

func createAndStartContainerOpts(
	imageName,
	containerName string,
	cmd []string,
	env []string,
	workDir string,
	mounts []docker.Mount,
) (*docker.CreateContainerOptions, *docker.HostConfig) {

	vols := make(map[string]struct{})
	for _, mount := range mounts {
		vols[mount.Destination] = struct{}{}
	}
	binds := make([]string, len(mounts))
	for i, mount := range mounts {
		binds[i] = fmt.Sprintf("%s:%s", mount.Source, mount.Destination)
	}

	createOpts := &docker.CreateContainerOptions{
		Name: containerName,
		Config: &docker.Config{
			Env:        env,
			Cmd:        cmd,
			Image:      imageName,
			Volumes:    vols,
			Mounts:     mounts,
			WorkingDir: workDir,
		},
		HostConfig: &docker.HostConfig{},
	}
	hostConfig := &docker.HostConfig{Binds: binds}
	return createOpts, hostConfig
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

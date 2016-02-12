package handlers

import (
	"testing"

	"github.com/arschles/assert"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	testImageName     = "myimg"
	testContainerName = "myContainer"
	testCmd           = "mycmd"
	testEnv1          = "myenv1"
	testEnv2          = "myenv2"
	testWorkDir       = "mywd"
)

func TestCreateContainerOpts(t *testing.T) {
	mounts := []docker.Mount{
		docker.Mount{Name: "mymount", Source: "mysrc", Destination: "mydest", Mode: "rx"},
	}
	opts, _ := createAndStartContainerOpts(
		testImageName,
		testContainerName,
		[]string{testCmd},
		[]string{testEnv1, testEnv2},
		testWorkDir,
		mounts,
	)
	assert.Equal(t, opts.Name, testContainerName, "container name")
	assert.True(t, opts.Config != nil, "config was nil")
	assert.Equal(t, opts.Config.Image, testImageName, "image name")
	expectedEnv := []string{testEnv1, testEnv2}
	assert.Equal(t, opts.Config.Env, expectedEnv, "environment variables")
}

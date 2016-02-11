package handlers

import (
	"testing"

	"github.com/arschles/assert"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	img     = "myimg"
	workdir = "myworkdir"
	site    = "mysite"
	repo    = "myrepo"
	env1    = "myenv1"
	env2    = "myenv2"
)

func TestCreateContainerOpts(t *testing.T) {
	opts := createContainerOpts(img, workdir, site, repo, env1, env2)
	assert.True(t, strings.HasPrefix(opts.Name, fmt.Sprintf("build-%s-%s-%s-", site, org, repo)), "name was not valid")
	assert.Equal(t, opts.Config.Image, img, "docker image")
	expectedEnv := []string{
		"GO15VENDOREXPERIMENT=1",
		"SITE=" + site,
		"ORG=" + org,
		"REPO=" + repo,
		env1,
		env2,
	}
	assert.Equal(t, opts.Config.Env, expectedEnv, "environment variables")
	assert.Equal(t, opts.Volumes, map[string]struct{}{workdir: struct{}{}}, "volumes")
	assert.Equal(t, opts.Config.Mounts, []docker.Mount{
		docker.Mount{Name: "pwd", Source: workdir, Destination: absPwd, Mode: "rx"},
	})
	assert.True(t, opts.HostConfig != nil, "host config was nil")
	assert.Equal(t, *docker.HostConfig, docker.HostConfig{}, "host config")
}

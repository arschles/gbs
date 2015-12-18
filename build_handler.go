package main

import (
	"fmt"
	"net/http"

	"code.google.com/p/go-uuid/uuid"
	"github.com/arschles/gbs/log"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
)

func buildHandler(workdir string, dockerCl *docker.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		site, ok := mux.Vars(r)["site"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		org, ok := mux.Vars(r)["org"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		repo, ok := mux.Vars(r)["repo"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("building %s/%s/%s", site, org, repo)
		container, err := dockerCl.CreateContainer(docker.CreateContainerOptions{
			Name: fmt.Sprintf("build-%s-%s-%s-%s", site, org, repo, uuid.New()),
			Config: &docker.Config{
				Env:   []string{"GO15VENDOREXPERIMENT=1", "CGO_ENABLED=0", "SITE=" + site, "ORG=" + org, "REPO=" + repo},
				Cmd:   []string{"/bin/bash", "/pwd/build.sh"},
				Image: "golang:1.5.2",
				Mounts: []docker.Mount{
					docker.Mount{Name: "pwd", Source: workdir, Destination: "/pwd", Mode: "rx"},
				},
			},
			HostConfig: &docker.HostConfig{},
		})
		if err != nil {
			log.Errf("creating container (%s)", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := dockerCl.StartContainer(container.ID, &docker.HostConfig{}); err != nil {
			log.Errf("starting container (%s)", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := dockerCl.AttachToContainer(docker.AttachToContainerOptions{
			Container:    container.ID,
			OutputStream: w,
			ErrorStream:  w,
			Logs:         true,
			Stream:       true,
			Stdout:       true,
			Stderr:       true,
		}); err != nil {
			log.Errf("attaching to container (%s)", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		code, err := dockerCl.WaitContainer(container.ID)
		if err != nil {
			log.Errf("container exited with code %d (%s)", code, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// w.WriteHeader(http.StatusCreated)
		// w.Write([]byte(fmt.Sprintf("created as %s on the filesystem\n", repo)))
	})
}

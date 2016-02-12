package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arschles/gbs/log"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
)

const (
	defaultBuildImg = "quay.io/arschles/gbs-env:0.0.1"
	containerBinDir = "/gobin"
	containerGoPath = "/go"
)

func BuildURL() string {
	return fmt.Sprintf("/{%s}/{%s}/{%s}", site, org, repo)
}

func Build(workdir string, dockerCl *docker.Client) http.Handler {
	type req struct {
		BuildImage *string `json:"build_image"`
		CGOEnabled *bool   `json:"cgo_enabled"`
	}

	type resp struct {
		StatusURL string `json:"status_url"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			httpErrf(w, http.StatusInternalServerError, "server doesn't support flushing output")
			return
		}

		site, ok := mux.Vars(r)[site]
		if !ok {
			httpErrf(w, http.StatusBadRequest, "missing site in path")
			return
		}

		org, ok := mux.Vars(r)[org]
		if !ok {
			httpErrf(w, http.StatusBadRequest, "missing org in path")
			return
		}

		repo, ok := mux.Vars(r)[repo]
		if !ok {
			httpErrf(w, http.StatusBadRequest, "missing repo in path")
			return
		}

		buildImg := defaultBuildImg
		var env []string
		req := new(req)
		if err := json.NewDecoder(r.Body).Decode(req); err == nil {
			if req.BuildImage != nil {
				buildImg = *req.BuildImage
			}
			if req.CGOEnabled == nil || !*req.CGOEnabled {
				env = append(env, "CGO_ENABLED=0")
			} else {
				env = append(env, "CGO_ENABLED=1")
			}
		}
		defer r.Body.Close()

		env = append(
			env,
			"GO15VENDOREXPERIMENT=1",
			"SITE="+site,
			"ORG="+org,
			"REPO="+repo,
			"BIN_NAME="+repo,
			"BIN_DIR="+containerBinDir,
			"GOPATH="+containerGoPath,
		)
		createContainerOpts, hostConfig := createAndStartContainerOpts(
			buildImg,
			containerName(site, org, repo),
			nil,
			env,
			"/",
			[]docker.Mount{
				docker.Mount{Name: "bin", Source: workdir, Destination: containerBinDir, Mode: "rx"},
			},
		)

		container, err := dockerCl.CreateContainer(*createContainerOpts)
		if err != nil {
			log.Errf("creating container [%s]", err)
			httpErrf(w, http.StatusInternalServerError, "error creating container [%s]", err)
			return
		}

		if err := dockerCl.StartContainer(container.ID, hostConfig); err != nil {
			log.Errf("starting container [%s]", err)
			httpErrf(w, http.StatusInternalServerError, "error starting container [%s]", err)
			return
		}

		attachOpts, outputReader := attachToContainerOpts(container.ID)
		errCh := make(chan error)
		go func() {
			if err := dockerCl.AttachToContainer(attachOpts); err != nil {
				errCh <- err
			}
		}()

		w.WriteHeader(http.StatusCreated)
		go func(reader io.Reader) {
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				fmt.Fprintf(w, "%s\n", scanner.Text())
				flusher.Flush()
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(w, "error with scanner in attached container [%s]\n", err)
			}
		}(outputReader)

		code, err := dockerCl.WaitContainer(container.ID)
		if err != nil {
			log.Errf("waiting for container %s [%s]", container.ID, err)
			return
		}
		w.Write([]byte(fmt.Sprintf("exited with error code %d\n", code)))

		removeOpts := docker.RemoveContainerOptions{
			ID:            container.ID,
			RemoveVolumes: true,
			Force:         true,
		}
		if err := dockerCl.RemoveContainer(removeOpts); err != nil {
			log.Errf("removing container %s [%s]", container.ID, err)
		}
	})
}

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
	defaultBuildEnv = "quay.io/arschles/gbs-env:0.0.1"
)

func BuildURL() string {
	return fmt.Sprintf("/{%s}/{%s}/{%s}", site, org, repo)
}

type startBuildReq struct {
	BuildEnv string `json:"build_env"`
}

type startBuildResp struct {
	StatusURL string `json:"status_url"`
}

func Build(workdir string, dockerCl *docker.Client) http.Handler {
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

		buildEnv := defaultBuildEnv
		req := new(startBuildReq)
		if err := json.NewDecoder(r.Body).Decode(req); err == nil {
			buildEnv = req.BuildEnv
		}
		defer r.Body.Close()

		containerOpts := createContainerOpts(workdir, site, org, repo)
		container, err := dockerCl.CreateContainer(containerOpts)
		if err != nil {
			log.Errf("creating container [%s]", err)
			httpErrf(w, http.StatusInternalServerError, "error creating container [%s]", err)
			return
		}

		hostConfig := &docker.HostConfig{Binds: []string{fmt.Sprintf("%s:%s", workdir, absPwd)}}
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
	})
}

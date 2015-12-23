package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/arschles/gbs/log"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
)

type buildResp struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}

func Build(workdir string, dockerCl *docker.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		site, ok := mux.Vars(r)["site"]
		if !ok {
			httpErr(w, http.StatusBadRequest, "missing site in path")
			return
		}
		org, ok := mux.Vars(r)["org"]
		if !ok {
			httpErr(w, http.StatusBadRequest, "missing org in path")
			return
		}
		repo, ok := mux.Vars(r)["repo"]
		if !ok {
			httpErr(w, http.StatusBadRequest, "missing repo in path")
			return
		}
		log.Printf("building %s/%s/%s", site, org, repo)
		container, err := createContainer(dockerCl, workdir, site, org, repo)
		if err != nil {
			log.Errf("creating container [%s]", err)
			httpErr(w, http.StatusInternalServerError, "error creating container [%s]", err)
			return
		}

		log.Printf("starting container")
		if err := startContainer(dockerCl, container, workdir); err != nil {
			log.Errf("starting container [%s]", err)
			httpErr(w, http.StatusInternalServerError, "error starting container [%s]", err)
			return
		}

		log.Printf("attaching container")
		stdoutReader, stderrReader, err := attachContainer(dockerCl, container)
		if err != nil {
			log.Errf("attaching container [%s]", err)
			httpErr(w, http.StatusInternalServerError, "error attaching container [%s]", err)
			return
		}

		log.Printf("waiting for container")
		exitCode, err := dockerCl.WaitContainer(container.ID)
		if err != nil {
			log.Errf("waiting for container [%s]", err)
			httpErr(w, http.StatusInternalServerError, "error waiting for container [%s]", err)
			return
		}

		log.Printf("done. outputting results")
		stdout, err := ioutil.ReadAll(stdoutReader)
		if err != nil {
			log.Errf("reading stdout [%s]", err)
			httpErr(w, http.StatusInternalServerError, "error reading STDOUT [%s]", err)
			return
		}
		stderr, err := ioutil.ReadAll(stderrReader)
		if err != nil {
			log.Errf("reading stderr [%s]", err)
			httpErr(w, http.StatusInternalServerError, "error reading STDERR [%s]", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		flusher, ok := w.(http.Flusher)
		resp := &buildResp{Stdout: string(stdout), Stderr: string(stderr), ExitCode: exitCode}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Errf("encoding JSON [%s]", err)
			httpErr(w, http.StatusInternalServerError, "error encoding JSON [%s]", err)
			return
		}
	})
}

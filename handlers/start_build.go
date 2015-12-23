package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arschles/gbs/log"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
)

func StartBuildURL() string {
	return fmt.Sprintf("/{%s}/{%s}/{%s}", site, org, repo)
}

type startBuildResp struct {
	StatusURL string `json:"status_url"`
}

func StartBuild(workdir string, dockerCl *docker.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

		container, err := createContainer(dockerCl, workdir, site, org, repo)
		if err != nil {
			log.Errf("creating container [%s]", err)
			httpErrf(w, http.StatusInternalServerError, "error creating container [%s]", err)
			return
		}

		if err := startContainer(dockerCl, container, workdir); err != nil {
			log.Errf("starting container [%s]", err)
			httpErrf(w, http.StatusInternalServerError, "error starting container [%s]", err)
			return
		}

		statURL := buildStatusURL(container.ID)
		w.Header().Set("Location", statURL)
		w.WriteHeader(http.StatusCreated)
		resp := startBuildResp{StatusURL: statURL}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Errf("encoding JSON [%s]", err)
			httpErrf(w, http.StatusInternalServerError, "error encoding JSON [%s]", err)
			return
		}
	})
}

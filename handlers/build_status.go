package handlers

import (
	"fmt"
	"net/http"

	"github.com/arschles/gbs/log"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
)

func BuildStatusURL() string {
	return fmt.Sprintf("/status/{%s}", containerID)
}

func BuildStatus(dockerCl *docker.Client) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		containerID, ok := mux.Vars(r)[containerID]
		if !ok {
			httpErrf(w, http.StatusBadRequest, "missing container ID in path")
			return
		}

		stdoutReader, stderrReader, doneCh, errCh := attachContainer(dockerCl, containerID)

		go streamReader(stdoutReader, w, httpFlushWriter)
		go streamReader(stderrReader, w, httpFlushWriter)

		select {
		case <-doneCh:
		case err := <-errCh:
			log.Errf("error attaching to container [%s]", err)
		}

		exitCode, err := dockerCl.WaitContainer(containerID)
		if err != nil {
			log.Errf("waiting for container [%s]", err)
			httpErrf(w, http.StatusInternalServerError, "error waiting for container [%s]", err)
			return
		}

		w.Write([]byte(fmt.Sprintf("exit code %d", exitCode)))

	})
}

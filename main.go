package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arschles/gbs/handlers"
	"github.com/arschles/gbs/log"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

const appName = "gbs"

type config struct {
	Port int `envconfig:"port",default:8080`
}

func main() {
	var cfg config
	if err := envconfig.Process(appName, &cfg); err != nil {
		log.Errf("loading config [%s]", err)
		os.Exit(1)
	}
	port := 8080
	cwd, err := os.Getwd()
	if err != nil {
		log.Errf("geting current working dir (%s)", err)
		os.Exit(1)
	}
	dockerCl, err := docker.NewClientFromEnv()
	if err != nil {
		log.Errf("creating new docker client (%s)", err)
		os.Exit(1)
	}
	r := mux.NewRouter()
	r.Handle(handlers.StartBuildURL(), handlers.Build(cwd, dockerCl)).Methods("POST")

	log.Printf("listening on port %d", port)
	hostStr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(hostStr, r); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

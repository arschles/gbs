package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"

	"gopkg.in/v2/yaml"
)

const filename = "gbs.yaml"

type buildFile struct {
	Repos []string `yaml:"repos"`
}

var workdir string

func init() {
	w, err := os.Getwd()
	if err != nil {
		fmt.Println("can't get current working dir: ", err)
		os.Exit(1)
	}
	workdir = w
}

func buildHandler(w http.ResponseWriter, r *http.Request) {
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
	script := []string{
		fmt.Sprintf("git clone https://%s/%s/%s.git", site, org, repo),
		fmt.Sprintf("mkdir -p /go/src/%s/%s", site, org),
		fmt.Sprintf("mv %s /go/src/%s/%s/%s", repo, site, org, repo),
		fmt.Sprintf("cd /go/src/%s/%s/%s", site, org, repo),
		fmt.Sprintf("go build -o /pwd/%s", repo),
	}
	remainder := strings.Join(script, " && ")
	cmd := exec.Command("docker",
		"run",
		// "--rm",
		"-e GO15VENDOREXPERIMENT=1",
		"-e CGO_ENABLED=0",
		fmt.Sprintf(`-v "%s":"/pwd"`, workdir),
		"golang:1.5.2",
		fmt.Sprintf(`sh -c '%s'`, remainder),
	)
	fmt.Println(strings.Join(cmd.Args, " "))
	cmd.Env = os.Environ()
	b, err := cmd.CombinedOutput()
	fmt.Println(string(b))
	if err != nil {
		fmt.Println("ERROR: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("created as %s on the filesystem", repo)))
}

func main() {
	fb, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("ERROR reading %s (%s)\n", filename, err)
		os.Exit(1)
	}

	var builds buildFile
	if err := yaml.Unmarshal(fb, &builds); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/{site}/{org}/{repo}", buildHandler).Methods("POST")

	fmt.Println("listening on 8080")
	http.ListenAndServe(":8080", r)

}

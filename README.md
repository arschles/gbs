# gbs

This repository builds a server that can download code for Go repositories and then build it.

# Building it

You'll need [go](http://golang.org) and the standard [GNU make](https://www.gnu.org/software/make/) installed. Once you have those, run the following:

```console
make bootstrsap
go build -o gbs
```

# Running It

_Note: Running inside a Docker container is coming_

Once you have a build ready, simply run it like this:

```console
./gbs
```

You can then make requests against it like the following, to build [mc](https://github.com/minio/mc):

```console
curl -XPOST localhost:8080/github.com/minio/mc
```

# Notes

All builds are built with the following environment variables set:

- `GO15VENDOREXPERIMENT=1`
- `CGO_ENABLED=0`

# TODOs

See [issues](https://github.com/arschles/gbs) for TODOs

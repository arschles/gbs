DOCKER_CMD := docker run -e GO15VENDOREXPERIMENT=1 -e CGO_ENABLED=0 --rm -v ${PWD}:/go/src/github.com/arschles/gbs -w /go/src/github.com/arschles/gbs quay.io/deis/go-dev:0.3.0
VERSION ?= 0.0.1
bootstrap:
	${DOCKER_CMD} glide up

build:
	${DOCKER_CMD} go build -o gbs

docker-build:
	docker build -t quay.io/arschles/gbs:$VERSION .

docker-push:
	docker push quay.io/arschles/gbs:$VERSION

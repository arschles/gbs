DOCKER_CMD := docker run -e GO15VENDOREXPERIMENT=1 -e CGO_ENABLED=0 --rm -v ${CURDIR}:/go/src/github.com/arschles/gbs -w /go/src/github.com/arschles/gbs quay.io/deis/go-dev:0.5.0
VERSION ?= 0.0.1
IMAGE_NAME := quay.io/arschles/gbs:${VERSION}
TEST_SERVER_IP ?= $(shell docker-machine ip dev)

bootstrap:
	${DOCKER_CMD} glide install

glideup:
	${DOCKER_CMD} glide up

glideget:
	${DOCKER_CMD} glide get ${PACKAGE}

build:
	${DOCKER_CMD} go build -o gbs

test:
	${DOCKER_CMD} sh -c 'go test $$(glide nv)'

run:
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v ${CURDIR}:/pwd -w /pwd ubuntu:14.04 ./gbs

docker-build:
	docker build -t ${IMAGE_NAME} .

docker-push:
	docker push ${IMAGE_NAME}

docker-build-env:
	make -C build-env docker-build

docker-push-env:
	make -C build-env docker-push

test-integration:
	curl -v -XPOST ${TEST_SERVER_IP}:8080/github.com/minio/mc

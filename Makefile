DOCKER_CMD := docker run -e GO15VENDOREXPERIMENT=1 -e CGO_ENABLED=0 --rm -v ${CURDIR}:/go/src/github.com/arschles/gbs -w /go/src/github.com/arschles/gbs quay.io/deis/go-dev:0.3.0
VERSION ?= 0.0.1
IMAGE_NAME := quay.io/arschles/gbs:${VERSION}

bootstrap:
	${DOCKER_CMD} glide up

build:
	${DOCKER_CMD} go build -o gbs

run:
	docker run --rm --net=host -v ${CURDIR}:/pwd -w /pwd quay.io/deis/go-dev:0.3.0 ./gbs

docker-build:
	docker build -t ${IMAGE_NAME} .

docker-push:
	docker push ${IMAGE_NAME}

docker-build-env:
	make -C build-env docker-build

docker-push-env:
	make -C build-env

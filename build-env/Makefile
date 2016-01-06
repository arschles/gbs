ENV_VERSION ?= 0.0.1
ENV_IMAGE_NAME := quay.io/arschles/gbs-env:${ENV_VERSION}

docker-build:
	docker build -t ${ENV_IMAGE_NAME} .

docker-push:
	docker push ${ENV_IMAGE_NAME}

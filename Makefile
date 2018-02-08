REGISTRY=kube-registry.kube.local
IMAGE=gopenguin/minimal-ldap-proxy
TAG=$(shell git rev-parse --verify HEAD)

.PHONY: all clear

all: push

push: tag
	docker push ${REGISTRY}/${IMAGE}:${TAG}

tag: build
	docker tag ${IMAGE}:${TAG} ${REGISTRY}/${IMAGE}:${TAG}

build:
	docker build -t ${IMAGE}:${TAG} .

clear:
	docker rmi --force $(docker images --format="{{.ID}}\t{{.Repository}}" | grep "${IMAGE}" | cut -f 1)


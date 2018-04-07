REGISTRY=kube-registry.kube.local
IMAGE=gopenguin/minimal-ldap-proxy
TAG=$(shell git rev-parse --verify HEAD)

DOCKER_HLP=.docker-hlp

GO_BUILD_ARGS=-tags "json1 fts5 sqlite_omit_load_extension ${GO_BUILD_TAGS}" -a -installsuffix cgo -ldflags "-linkmode external -extldflags \"-static -lc\" -w -s"

GO_SRC=${shell find ${PWD} -type f -name "*.go" | grep -v vendor}

.PHONY: all clean

all: run

push: ${DOCKER_HLP}/push
tag: ${DOCKER_HLP}/tag
build: ${DOCKER_HLP}/build

${DOCKER_HLP}/push: ${DOCKER_HLP}/tag
	docker push ${REGISTRY}/${IMAGE}:${TAG}
	touch ${DOCKER_HLP}/push

${DOCKER_HLP}/tag: ${DOCKER_HLP}/build
	docker tag ${IMAGE}:${TAG} ${REGISTRY}/${IMAGE}:${TAG}
	touch ${DOCKER_HLP}/tag

${DOCKER_HLP}/build: ${GO_SRC}
	docker build -t ${IMAGE}:${TAG} .
	mkdir -p ${DOCKER_HLP} && touch ${DOCKER_HLP}/build

clean-docker:
	docker rmi --force $$(docker images --format="{{.ID}}\t{{.Repository}}" | grep "${IMAGE}" | cut -f 1)
	rm -rf ${DOCKER_HLP}

minimal-ldap-proxy-static: ${GO_SRC}
	CGO_ENABLED=1 GOOS=linux go build ${GO_BUILD_ARGS} -o minimal-ldap-proxy-static .

minimal-ldap-proxy: ${GO_SRC}
	go build --tags "${GO_BUILD_TAGS}" -o minimal-ldap-proxy .

run:
	@CGO_ENABLED=1 GOOS=linux go run ${GO_BUILD_ARGS} main.go

clean: clean-docker
	rm -rf minimal-ldap-proxy

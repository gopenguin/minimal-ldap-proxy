GO_BUILD_TAGS=postgres mysql
GO_BUILD_STATIC_ARGS=-tags "json1 fts5 sqlite_omit_load_extension ${GO_BUILD_TAGS}" -a -installsuffix cgo -ldflags "-linkmode external -extldflags \"-static -lc\" -w -s"

GO_SRC=${shell find . -name '*.go' -not -path './vendor/*'}

CMDS=minimal-ldap-proxy mlpcli

.PHONY: all
all: ${CMDS}

all-static: $(addsuffix -static, $(CMDS))

${CMDS}: %: ${GO_SRC}
	go build -tags "${GO_BUILD_TAGS}" -o $@ ./cmd/$@

$(addsuffix -static, $(CMDS)): %-static:
	CGO_ENABLED=1 GOOS=linux go build ${GO_BUILD_STATIC_ARGS} -o $@ ./cmd/$(@:%-static=%)/main.go

$(addprefix run-, $(CMDS)): run-%:
	@CGO_ENABLED=1 GOOS=linux go run -tags "${GO_BUILD_TAGS}" ./cmd/$(@:run-%=%)/main.go

.PHONY: clean
clean:
	rm -f ${CMDS}
	rm -f $(addsuffix -static, $(CMDS))

###############################################################################
### DOCKER
###############################################################################

DOCKER_HLP=.docker-hlp
REGISTRY=kube-registry.kube.local
IMAGE=gopenguin/minimal-ldap-proxy
TAG=$(shell git rev-parse --verify HEAD)

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

REPO?=bhakiyakalimuthu/fileshare
VERSION?=v1.0.0
GOCMD?=CGO_ENABLED=0 go
GO_SERVER_MAIN_SRC?=cmd/server/main.go
GO_CLIENT_MAIN_SRC?=cmd/server/main.go

DOCKER_PLATFORM := --platform linux/amd64

.PHONY: test
test:
	${GOCMD} test -v ./... -mod=vendor -count=1

.PHONY:image
image:
	docker build -f Dockerfile --build-arg=VERSION=${VERSION} -T ${REPO}:${VERSION}

.PHONY:build_server
build_server:
	${GOCMD} build -mod vendor --ldflag "-X main.serviceVersion=${VERSION}" -o go-app ${GO_SERVER_MAIN_SRC}

.PHONY:build_client
build_client:
	${GOCMD} build -mod vendor --ldflag "-X main.serviceVersion=${VERSION}" -o go-app ${GO_CLIENT_MAIN_SRC}




.PHONY: proto_image
proto_image: ## Build proto-generation docker image.
	rm -rf tmp/protoc-gen-event; \
	git clone git@github.com:voiapp/protoc-gen-event.git tmp/protoc-gen-event; \
	docker build $(DOCKER_PLATFORM) -f Dockerfile-proto --build-arg=VERSION=$(VERSION) -t $(PROTO_IMAGE) . ;\
	rm -rf tmp/protoc-gen-event;


proto: proto_image ## Run go generate in docker.
	docker run $(DOCKER_PLATFORM) --rm -it -w /go/in -v $(CURDIR):/go/in $(PROTO_IMAGE) go generate -mod=vendor ./...

generate: proto ## Run go generate in docker (alias).
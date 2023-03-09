GO111MODULE := on
DOCKER_TAG := $(or ${GITHUB_TAG_NAME}, latest)

.PHONY: proto
proto:
	make -C proto protoc

.PHONY: server
server:
	go build -tags netgo -o bin/server server/main.go
	strip bin/server

.PHONY: client
client:
	go build -tags netgo -o bin/client client/main.go
	strip bin/client

.PHONY: dockerimage
dockerimage:
	docker build -t ghcr.io/metal-stack/droptailer:${DOCKER_TAG} .
	docker build -f Dockerfile.client -t ghcr.io/metal-stack/droptailer-client:${DOCKER_TAG} .

.PHONY: dockerpush
dockerpush:
	docker push ghcr.io/metal-stack/droptailer:${DOCKER_TAG}
	docker push ghcr.io/metal-stack/droptailer-client:${DOCKER_TAG}

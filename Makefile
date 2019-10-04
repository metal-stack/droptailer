GO111MODULE := on

proto:
	protoc -I dropsink/ dropsink/dropsink.proto --go_out=plugins=grpc:dropsink

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
	docker build -t droptailer .
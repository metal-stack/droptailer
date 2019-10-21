GO111MODULE := on

.PHONY: proto
proto:
	protoc -I proto/ proto/droptailer.proto --go_out=plugins=grpc:proto

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
	#docker build -t metalpod/droptailer:latest .
	docker build -f Dockerfile.client -t metalpod/droptailer-client:latest .

.PHONY: dockerpush
dockerpush:
	docker push metalpod/droptailer:latest
	docker push metalpod/droptailer-client:latest

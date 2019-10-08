FROM golang:1.13-buster as builder
COPY / /work
WORKDIR /work
RUN apt-get update \
 && apt-get install libsystemd-dev \
 && make client server

FROM alpine:3.10
COPY --from=builder /work/bin/server /server
COPY --from=builder /work/bin/client /client
ENV SERVER_PORT=50051
EXPOSE ${SERVER_PORT}
USER nobody
CMD ["/server"]
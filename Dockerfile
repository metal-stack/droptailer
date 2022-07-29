FROM golang:1.18-buster as builder
COPY / /work
WORKDIR /work
RUN make server

FROM alpine:3.16
COPY --from=builder /work/bin/server /server
ENV SERVER_PORT=50051
EXPOSE ${SERVER_PORT}
USER nobody
CMD ["/server"]

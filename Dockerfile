FROM golang:1.16-buster as builder
COPY / /work
WORKDIR /work
RUN make server

FROM alpine:3.13
COPY --from=builder /work/bin/server /server
ENV SERVER_PORT=50051
EXPOSE ${SERVER_PORT}
USER nobody
CMD ["/server"]

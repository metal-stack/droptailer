FROM golang:1.22-bookworm as builder
COPY / /work
WORKDIR /work
RUN make server

FROM alpine:3.19
COPY --from=builder /work/bin/server /server
ENV SERVER_PORT=50051
EXPOSE ${SERVER_PORT}
USER nobody
CMD ["/server"]

FROM golang:1.24-bookworm AS builder
COPY / /work
WORKDIR /work
RUN make server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /work/bin/server /server
ENV SERVER_PORT=50051
EXPOSE ${SERVER_PORT}
USER nobody
CMD ["/server"]

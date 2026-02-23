FROM golang:1.26-trixie AS builder
COPY / /work
WORKDIR /work
RUN make server

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=builder /work/bin/server /server
ENV SERVER_PORT=50051
EXPOSE ${SERVER_PORT}
USER nobody
CMD ["/server"]

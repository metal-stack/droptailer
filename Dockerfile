FROM gcr.io/distroless/static-debian13:nonroot
COPY bin/server /server
ENV SERVER_PORT=50051
EXPOSE ${SERVER_PORT}
USER nobody
CMD ["/server"]

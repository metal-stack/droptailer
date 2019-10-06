# Droptailer

Droptailer gathers packet drops from different machines, enriches them with data from kubernetes api resources and makes them accessible by kubernetes means.

## Client

- reads the systemd journal for kernel log messages about packet drops
- pushes them with gRPC to the `droptail` server

environment variables:
- `DROPTAILER_SERVER_ADDRESS`: endpoint for the server
- `DROPTAILER_PREFIXES_OF_DROPS`: prefixes that identify drop messages in the journal


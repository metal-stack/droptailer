# Droptailer

Droptailer gathers packet drops from different machines, enriches them with data from kubernetes api resources and makes them accessible by kubernetes means.

## Client

- reads the systemd journal for kernel log messages about packet drops
- pushes them with gRPC to the `droptail` server

environment variables:

- `DROPTAILER_SERVER_ADDRESS`: endpoint for the server
- `DROPTAILER_PREFIXES_OF_DROPS`: prefixes that identify drop messages in the journal

## Generating certificates

```bash
# Install cfssl tool
curl -s -L -o ~/bin/cfssl https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
curl -s -L -o ~/bin/cfssljson https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
chmod +x ~/bin/{cfssl,cfssljson}

# Create certificates for client and server
echo '{"CN":"CA","key":{"algo":"rsa","size":2048}}' | cfssl gencert -initca - | cfssljson -bare ca -
echo '{"signing":{"default":{"expiry":"43800h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json
export ADDRESS=droptailer
export NAME=droptailer-server
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' \
    | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - \
    | cfssljson -bare $NAME

export ADDRESS=
export NAME=droptailer-client
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' \
    | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - \
    | cfssljson -bare $NAME
```

## Testing droptailer

```bash
# install kind 0.6.0 or higher !
KIND_VERSION=v0.7.0
wget https://github.com/kubernetes-sigs/kind/releases/download/${KIND_VERSION}/kind-linux-amd64
mv kind-linux-amd64 ~/bin/kind
chmod +x ~/bin/kind

# Create a k8s cluster
kind create cluster

# Deploy droptailer-server
kubectl apply -f ./test/manifests/droptailer.yaml

# Expose droptailer-server port to host
podName=$(kubectl get pods -n firewall -o=jsonpath='{.items[0].metadata.name}')
echo $podName
kubectl port-forward -n firewall --address 0.0.0.0 pod/$podName 50051:50051 &

# Run droptailer-client
docker run -it \
  --privileged \
  --add-host droptailer:172.17.0.1 \
  --env DROPTAILER_SERVER_ADDRESS=droptailer:50051 \
  --volume $(pwd)/test/certs:/etc/droptailer-client:ro \
  --volume /run/systemd/private:/run/systemd/private \
  --volume /var/log/journal:/var/log/journal \
  --volume /run/log/journal:/run/log/journal \
  --volume /etc/machine-id:/etc/machine-id \
metalstack/droptailer-client

# Watch for drops
stern -n firewall drop

# Generate a sample message for the systemd journal that gets catched by the droptailer-client
sudo logger -t kernel "nftables-metal-dropped: IN=vrf09 OUT= MAC=12:99:fd:3b:ce:f8:1a:ae:e9:a7:95:50:08:00 SRC=1.2.3.4 DST=4.3.2.1 LEN=40 TOS=0x00 PREC=0x00 TTL=238 ID=46474 PROTO=TCP SPT=59265 DPT=445 WINDOW=1024 RES=0x00 SYN URGP=0"
```

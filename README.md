# Droptailer

Droptailer gathers packet drops from different machines, enriches them with data from kubernetes api resources and makes them accessible by kubernetes means.

## Client

- reads the systemd journal for kernel log messages about packet drops
- pushes them with gRPC to the `droptail` server

environment variables:

- `DROPTAILER_SERVER_ADDRESS`: endpoint for the server
- `DROPTAILER_PREFIXES_OF_DROPS`: prefixes that identify drop messages in the journal

## TLS between client and server

### Create self signed certificates for client and server

~~~bash
curl -s -L -o ~/bin/cfssl https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
curl -s -L -o ~/bin/cfssljson https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
chmod +x ~/bin/{cfssl,cfssljson}
~~~


~~~bash
echo '{"CN":"CA","key":{"algo":"rsa","size":2048}}' | cfssl gencert -initca - | cfssljson -bare ca -
echo '{"signing":{"default":{"expiry":"43800h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json
export ADDRESS=10.244.0.14,droptailer
export NAME=server
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' \
    | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - \
    | cfssljson -bare $NAME
export ADDRESS=
export NAME=client
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":2048}}' \
    | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - \
    | cfssljson -bare $NAME
~~~


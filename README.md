```
cfssl gencert -initca ./tls/ca-csr.json | cfssljson -bare /tmp/ca

cfssl gencert \
-ca=/tmp/ca.pem \
-ca-key=/tmp/ca-key.pem \
-config=./tls/ca-config.json \
-hostname="tcpdump-webhook,tcpdump-webhook.default.svc.cluster.local,tcpdump-webhook.default.svc,localhost,127.0.0.1" \
-profile=default \
./tls/ca-csr.json | cfssljson -bare /tmp/tcpdump-webhook

mv /tmp/tcpdump-webhook.pem ./ssl/tcpdump.pem
mv /tmp/tcpdump-webhook-key.pem ./ssl/tcpdump.key

Update ConfigMap data in the manifest/webhook-deployment.yaml file with your key and certificate.
cat ./ssl/tcpdump.key | base64 | tr -d '\n'
cat ./ssl/tcpdump.pem | base64 | tr -d '\n'

Update caBundle value in the manifest/webhook-configuration.yaml file with your base64 encoded CA certificate.
cat /tmp/ca.pem | base64 | tr -d '\n'

docker build -t dyslexicat/tcpdump-webhook .
docker build -t dyslexicat/tcpdump-alpine -f docker/Dockerfile .


kubectl apply -f manifests/webhook-deployment.yaml
kubectl apply -f manifests/webhook-configuration.yaml

kubectl apply -f manifests/test-pod.yaml
```
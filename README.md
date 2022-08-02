# k8s-container-integrity-mutator
## Architecture
### Statechart diagram
![File location: docs/diagrams/mutatorStatechartDiagram.png](/docs/diagrams/mutatorStatechartDiagram.png?raw=true "Statechart diagram")
### Sequence diagram
![File location: docs/diagrams/mutatorSequenceDiagram.png](/docs/diagrams/mutatorSequenceDiagram.png?raw=true "Sequence diagram")
## Quick start
Generate CA in /tmp :
```
cfssl gencert -initca ./certificates/tls/ca-csr.json | cfssljson -bare /tmp/ca
```

Generate private key and certificate for SSL connection:
```
cfssl gencert \
-ca=/tmp/ca.pem \
-ca-key=/tmp/ca-key.pem \
-config=./certificates/tls/ca-config.json \
-hostname="k8s-webhook-injector,k8s-webhook-injector.default.svc.cluster.local,k8s-webhook-injector.default.svc,localhost,127.0.0.1" \
-profile=default \
./certificates/tls/ca-csr.json | cfssljson -bare /tmp/k8s-webhook-injector
```

Move your SSL key and certificate to the ssl directory:
```
mkdir webhook/ssl
mv /tmp/k8s-webhook-injector.pem ./certificates/ssl/k8s-webhook-injector.pem
mv /tmp/k8s-webhook-injector-key.pem ./certificates/ssl/k8s-webhook-injector.key
```

Update configuration data in the manifests/webhook/webhook-configMap.yaml file with your key in the appropriate field `data:server.key` and certificate in the appropriate field `data:server.crt:`:
```
cat ./certificates/ssl/k8s-webhook-injector.key | base64 | tr -d '\n'
cat ./certificates/ssl/k8s-webhook-injector.pem | base64 | tr -d '\n'
```

Update field `caBundle` value in the manifests/webhook/webhook-configuration.yaml file with your base64 encoded CA certificate:
```
cat /tmp/ca.pem | base64 | tr -d '\n'
```

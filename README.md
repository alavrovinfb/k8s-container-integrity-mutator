![GitHub contributors](https://img.shields.io/github/contributors/ScienceSoft-Inc/k8s-container-integrity-mutator)
![GitHub last commit](https://img.shields.io/github/last-commit/ScienceSoft-Inc/k8s-container-integrity-mutator)
![GitHub issues](https://img.shields.io/github/issues/ScienceSoft-Inc/k8s-container-integrity-mutator)
![GitHub forks](https://img.shields.io/github/forks/ScienceSoft-Inc/k8s-container-integrity-mutator)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)

# k8s-container-integrity-mutator
This application provides the injection of any patch inside any k8s schemas like sidecar.

When applying a new scheme to a cluster, the application monitors the presence of a "hasher-certificates-injector-sidecar" label and, if available, makes a patch.

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

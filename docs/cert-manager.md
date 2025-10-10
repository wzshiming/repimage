# Cert-Manager Integration

This document explains how to use cert-manager with repimage for automatic certificate management.

## Prerequisites

- Kubernetes cluster (1.16+)
- [cert-manager](https://cert-manager.io/) installed in your cluster

## Installing cert-manager

If you don't have cert-manager installed, you can install it using:

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
```

Wait for cert-manager to be ready:

```bash
kubectl wait --for=condition=Available --timeout=300s \
  -n cert-manager deployment/cert-manager \
  deployment/cert-manager-webhook \
  deployment/cert-manager-cainjector
```

## How it works

When using cert-manager with repimage:

1. **Self-signed Issuer**: A self-signed Issuer is created in the `kube-system` namespace
2. **Certificate Resource**: A Certificate resource is created that requests a TLS certificate for the webhook service
3. **Secret Creation**: cert-manager automatically creates a Secret (`repimage-tls`) containing the TLS certificate and key
4. **CA Injection**: cert-manager automatically injects the CA bundle into the MutatingWebhookConfiguration
5. **Auto-renewal**: cert-manager automatically renews certificates before they expire (15 days before expiry by default)

## Installation with cert-manager

Deploy repimage with cert-manager support:

```bash
kubectl apply -k https://github.com/wzshiming/repimage/yaml
```

Or using kustomize locally:

```bash
kubectl apply -k yaml/
```

## Verify Installation

Check that the certificate was created successfully:

```bash
kubectl get certificate -n kube-system repimage-cert
```

You should see output similar to:

```
NAME            READY   SECRET         AGE
repimage-cert   True    repimage-tls   1m
```

Check the webhook configuration has the CA bundle injected:

```bash
kubectl get mutatingwebhookconfiguration repimage -o jsonpath='{.webhooks[0].clientConfig.caBundle}' | base64 -d | openssl x509 -text -noout | head -10
```

## Using a Different Issuer

If you want to use a different issuer (e.g., Let's Encrypt or your own CA), modify the `certificate.yaml` file:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: repimage-cert
  namespace: kube-system
spec:
  # ... other settings ...
  issuerRef:
    name: your-issuer-name  # Change this
    kind: ClusterIssuer      # Or use ClusterIssuer
    group: cert-manager.io
```

## Troubleshooting

### Certificate not ready

Check the certificate status:

```bash
kubectl describe certificate -n kube-system repimage-cert
```

### Webhook failing to start

Check the pod logs:

```bash
kubectl logs -n kube-system -l app=repimage
```

Verify the secret exists and contains the certificate:

```bash
kubectl get secret -n kube-system repimage-tls
kubectl get secret -n kube-system repimage-tls -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -text -noout
```

### CA bundle not injected

Make sure cert-manager's CA injector is running:

```bash
kubectl get pods -n cert-manager -l app.kubernetes.io/name=cainjector
```

Check the annotation is present on the webhook configuration:

```bash
kubectl get mutatingwebhookconfiguration repimage -o yaml | grep cert-manager
```

## Migrating from Static Certificates

If you're currently using static certificates and want to migrate to cert-manager:

1. Install cert-manager in your cluster
2. Apply the cert-manager configuration:
   ```bash
   kubectl apply -k https://github.com/wzshiming/repimage/yaml
   ```
3. The deployment will be updated to use the cert-manager generated certificates
4. The old static certificates in the container image will be ignored

## Certificate Renewal

Certificates are automatically renewed by cert-manager. The default configuration:

- Certificate valid for: 90 days (2160h)
- Renewal before: 15 days (360h)

You can modify these values in the `certificate.yaml` file if needed.

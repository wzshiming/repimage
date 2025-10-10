# Quick Start Guide: Allowlist Configuration

## Default Configuration (Zero Config)

The application works out-of-the-box with embedded defaults:

```bash
kubectl apply -f yaml/webhook.yaml
kubectl apply -f yaml/admission.yaml
```

The image includes:
- ✅ Pre-downloaded allowlist at build time
- ✅ Hourly updates from DaoCloud mirror
- ✅ Embedded fallback for 6 common registries

## Using Custom Allowlist File

### Step 1: Create ConfigMap with your allowlist

```bash
kubectl create configmap repimage-allowlist \
  --from-file=allowlist.txt \
  -n kube-system
```

Your `allowlist.txt` format:
```
docker.io=my-mirror.com/docker.io
gcr.io=my-mirror.com/gcr.io
```

### Step 2: Deploy with ConfigMap

```bash
kubectl apply -f yaml/webhook-with-configmap.yaml
kubectl apply -f yaml/admission.yaml
```

## Using Custom Mirror URL

Edit `yaml/webhook.yaml` and add:

```yaml
env:
  - name: ALLOWLIST_URL
    value: "https://my-company.com/mirror-allowlist.txt"
  - name: ALLOWLIST_UPDATE_INTERVAL
    value: "15m"  # Update every 15 minutes
```

Then deploy:
```bash
kubectl apply -f yaml/webhook.yaml
kubectl apply -f yaml/admission.yaml
```

## Disable Periodic Updates

For stable, build-time only configuration:

```yaml
env:
  - name: ALLOWLIST_FILE
    value: "/etc/repimage/allowlist.txt"
  - name: ALLOWLIST_UPDATE_INTERVAL
    value: "0"  # Disable updates
```

## Verify Configuration

Check that the allowlist is loaded correctly:

```bash
kubectl logs -n kube-system deployment/repimage | grep -i "domain map"
```

Expected output:
```
I1010 09:12:24.779406 parse.go:92] Loaded domain map from file: /etc/repimage/allowlist.txt
I1010 09:12:24.779419 parse.go:123] Domain map updated with 6 entries at 2025-10-10 09:12:24
```

## Test Image Replacement

Create a test pod:

```bash
kubectl run test-nginx --image=nginx --dry-run=client -o yaml | kubectl apply -f -
```

Check the actual image used:

```bash
kubectl get pod test-nginx -o jsonpath='{.spec.containers[0].image}'
```

Should show: `m.daocloud.io/docker.io/library/nginx`

## Troubleshooting

### Check Environment Variables
```bash
kubectl get deployment repimage -n kube-system -o yaml | grep -A 10 env:
```

### Check Logs
```bash
kubectl logs -n kube-system deployment/repimage -f
```

### Common Issues

**Issue**: Allowlist not loading from file
- Check file path is correct in ALLOWLIST_FILE
- Verify ConfigMap is mounted correctly
- Check file permissions

**Issue**: Updates not happening
- Verify ALLOWLIST_UPDATE_INTERVAL is set and > 0
- Check network connectivity from pod
- Review logs for error messages

**Issue**: Using wrong mirror
- Verify allowlist file format (key=value)
- Check for typos in domain names
- Ensure no trailing spaces in allowlist

## Advanced: Build Custom Image

To build with a different default allowlist source:

```dockerfile
FROM alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && apk add --no-cache ca-certificates wget

ADD ./certs /certs
ADD ./bin/repimage /repimage
ADD ./scripts/download-allowlist.sh /tmp/download-allowlist.sh

# Use your custom allowlist URL
ENV ALLOWLIST_URL=https://my-company.com/allowlist.txt
RUN /tmp/download-allowlist.sh && rm /tmp/download-allowlist.sh

ENV ALLOWLIST_FILE=/etc/repimage/allowlist.txt
ENV ALLOWLIST_UPDATE_INTERVAL=1h

CMD ["/repimage"]
```

Build and push:
```bash
make build
docker build -t my-registry/repimage:latest .
docker push my-registry/repimage:latest
```

## Summary

| Scenario | Configuration | Update Interval |
|----------|--------------|-----------------|
| Default (recommended) | Build-time + URL updates | 1h |
| Static/Offline | ConfigMap file only | 0 (disabled) |
| Dynamic | Custom URL | 15m-1h |
| Custom Build | Modified Dockerfile | Custom |

For detailed configuration options, see [ALLOWLIST_CONFIG.md](ALLOWLIST_CONFIG.md).

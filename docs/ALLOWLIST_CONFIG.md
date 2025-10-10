# Allowlist Configuration Guide

## Overview

The repimage webhook now supports multiple ways to configure the domain allowlist for image replacement. The allowlist determines which domains will be replaced with mirror addresses.

## Configuration Options

### 1. Default Embedded Allowlist

The application includes a default embedded allowlist with common container registries:

- `docker.io` → `m.daocloud.io/docker.io`
- `gcr.io` → `m.daocloud.io/gcr.io`
- `k8s.gcr.io` → `m.daocloud.io/k8s.gcr.io`
- `registry.k8s.io` → `m.daocloud.io/registry.k8s.io`
- `ghcr.io` → `m.daocloud.io/ghcr.io`
- `quay.io` → `m.daocloud.io/quay.io`

This embedded allowlist is used as a fallback when other sources are unavailable.

### 2. Build-time Allowlist

The Docker image automatically downloads the latest allowlist during build time and stores it at `/etc/repimage/allowlist.txt`. This ensures the image has an up-to-date allowlist without requiring network access at runtime.

### 3. Environment Variables

Configure the allowlist behavior using environment variables:

#### `ALLOWLIST_FILE`

Specifies the path to a file containing the domain allowlist.

```yaml
env:
  - name: ALLOWLIST_FILE
    value: "/etc/repimage/allowlist.txt"
```

#### `ALLOWLIST_URL`

Specifies a URL to fetch the allowlist from. Defaults to the DaoCloud public-image-mirror repository.

```yaml
env:
  - name: ALLOWLIST_URL
    value: "https://your-custom-url.com/allowlist.txt"
```

#### `ALLOWLIST_UPDATE_INTERVAL`

Controls how often the allowlist is refreshed from the URL. Accepts duration strings like `30m`, `1h`, `24h`. Set to `0` to disable periodic updates.

```yaml
env:
  - name: ALLOWLIST_UPDATE_INTERVAL
    value: "1h"
```

## Allowlist File Format

The allowlist file uses a simple key-value format:

```
# Comments start with #
docker.io=m.daocloud.io/docker.io
gcr.io=m.daocloud.io/gcr.io
k8s.gcr.io=m.daocloud.io/k8s.gcr.io

# Blank lines are ignored
registry.k8s.io=m.daocloud.io/registry.k8s.io
```

## Loading Priority

The application loads the allowlist in the following priority order:

1. **File** (if `ALLOWLIST_FILE` is set and the file exists)
2. **URL** (if `ALLOWLIST_URL` is set and accessible)
3. **Default embedded allowlist** (fallback)

## Examples

### Using Custom Allowlist File

Create a ConfigMap with your custom allowlist:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: repimage-allowlist
  namespace: kube-system
data:
  allowlist.txt: |
    docker.io=my.mirror.io/docker.io
    gcr.io=my.mirror.io/gcr.io
```

Mount it in the webhook deployment:

```yaml
spec:
  containers:
    - name: webhook
      image: registry.cn-beijing.aliyuncs.com/laoq/repimage:latest
      env:
        - name: ALLOWLIST_FILE
          value: "/config/allowlist.txt"
        - name: ALLOWLIST_UPDATE_INTERVAL
          value: "0"  # Disable updates, use file only
      volumeMounts:
        - name: allowlist
          mountPath: /config
  volumes:
    - name: allowlist
      configMap:
        name: repimage-allowlist
```

### Using Custom URL with Periodic Updates

```yaml
spec:
  containers:
    - name: webhook
      image: registry.cn-beijing.aliyuncs.com/laoq/repimage:latest
      env:
        - name: ALLOWLIST_URL
          value: "https://my-company.com/mirror-config/allowlist.txt"
        - name: ALLOWLIST_UPDATE_INTERVAL
          value: "30m"  # Update every 30 minutes
```

### Using Build-time Allowlist Only

```yaml
spec:
  containers:
    - name: webhook
      image: registry.cn-beijing.aliyuncs.com/laoq/repimage:latest
      env:
        - name: ALLOWLIST_FILE
          value: "/etc/repimage/allowlist.txt"  # Use build-time allowlist
        - name: ALLOWLIST_UPDATE_INTERVAL
          value: "0"  # Disable updates
```

## Building Custom Image with Different Allowlist

You can customize the allowlist during image build:

```dockerfile
FROM alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && apk add --no-cache ca-certificates wget

ADD ./certs /certs
ADD ./bin/repimage /repimage
ADD ./scripts/download-allowlist.sh /tmp/download-allowlist.sh

# Use custom URL for build-time download
ENV ALLOWLIST_URL=https://my-company.com/allowlist.txt
RUN /tmp/download-allowlist.sh && rm /tmp/download-allowlist.sh

ENV ALLOWLIST_FILE=/etc/repimage/allowlist.txt
ENV ALLOWLIST_UPDATE_INTERVAL=1h
```

## Monitoring

The application logs information about allowlist loading:

```
I1010 09:12:24.779406 parse.go:92] Loaded domain map from file: /etc/repimage/allowlist.txt
I1010 09:12:24.779419 parse.go:123] Domain map updated with 6 entries at 2025-10-10 09:12:24
```

Check logs to verify the allowlist is loaded correctly:

```bash
kubectl logs -n kube-system deployment/repimage
```

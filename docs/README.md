# Documentation Index

## Overview

This directory contains comprehensive documentation for the allowlist configuration feature in repimage.

## Documents

### 🚀 [QUICKSTART.md](QUICKSTART.md)
**Start here!** Quick guide to get up and running with different configuration scenarios.

Topics covered:
- Default configuration (zero config)
- Using custom allowlist files with ConfigMap
- Using custom mirror URLs
- Disabling periodic updates
- Testing and troubleshooting

### 📖 [ALLOWLIST_CONFIG.md](ALLOWLIST_CONFIG.md)
Complete configuration reference guide.

Topics covered:
- All configuration options explained
- Environment variables reference
- Allowlist file format
- Loading priority order
- Detailed examples for each scenario
- Monitoring and verification

### 🏗️ [ARCHITECTURE.md](ARCHITECTURE.md)
Technical architecture and design documentation.

Topics covered:
- System flow diagrams
- Component architecture
- Loading mechanisms
- Caching strategy
- Error handling
- Performance considerations

### 📝 [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)
Implementation details and changes summary.

Topics covered:
- Problem statement
- Solution approach
- Files changed
- Key features
- Testing results
- Usage examples
- Backward compatibility

## Quick Reference

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ALLOWLIST_FILE` | Path to allowlist file | `/etc/repimage/allowlist.txt` |
| `ALLOWLIST_URL` | URL to fetch allowlist from | DaoCloud mirror URL |
| `ALLOWLIST_UPDATE_INTERVAL` | Refresh interval | `1h` |

### Loading Priority

1. **File** (if `ALLOWLIST_FILE` exists)
2. **URL** (if `ALLOWLIST_URL` is accessible)
3. **Embedded defaults** (fallback)

### Default Allowlist

```
docker.io       → m.daocloud.io/docker.io
gcr.io          → m.daocloud.io/gcr.io
k8s.gcr.io      → m.daocloud.io/k8s.gcr.io
registry.k8s.io → m.daocloud.io/registry.k8s.io
ghcr.io         → m.daocloud.io/ghcr.io
quay.io         → m.daocloud.io/quay.io
```

## Common Use Cases

### 1. Air-gapped Environment
Use ConfigMap with static allowlist, disable updates:
```yaml
env:
  - name: ALLOWLIST_FILE
    value: "/config/allowlist.txt"
  - name: ALLOWLIST_UPDATE_INTERVAL
    value: "0"
```

### 2. Corporate Mirror
Use custom URL with periodic updates:
```yaml
env:
  - name: ALLOWLIST_URL
    value: "https://mirror.company.com/allowlist.txt"
  - name: ALLOWLIST_UPDATE_INTERVAL
    value: "30m"
```

### 3. Default (Recommended)
Use build-time allowlist with hourly updates:
```yaml
env:
  - name: ALLOWLIST_FILE
    value: "/etc/repimage/allowlist.txt"
  - name: ALLOWLIST_UPDATE_INTERVAL
    value: "1h"
```

## Support

For issues or questions:
1. Check the [QUICKSTART.md](QUICKSTART.md) troubleshooting section
2. Review [ALLOWLIST_CONFIG.md](ALLOWLIST_CONFIG.md) for configuration details
3. Examine logs: `kubectl logs -n kube-system deployment/repimage`

## Contributing

When updating the allowlist feature:
1. Update relevant documentation files
2. Add tests to `pkg/utils/parse_config_test.go`
3. Update examples in `../yaml/` directory
4. Verify all documentation links work

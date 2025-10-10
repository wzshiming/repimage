# Allowlist Loading Architecture

## Flow Diagram

```
Application Startup
        |
        v
   InitDomainMap()
        |
        v
   Read Environment Variables
   - ALLOWLIST_FILE
   - ALLOWLIST_URL  
   - ALLOWLIST_UPDATE_INTERVAL
        |
        v
   Try Load from File
   (Priority 1: ALLOWLIST_FILE)
        |
        +--> Success? --> Cache & Use
        |
        +--> Failed/Not Set
                |
                v
        Try Load from URL
        (Priority 2: ALLOWLIST_URL)
                |
                +--> Success? --> Cache & Use
                |
                +--> Failed/Not Set
                        |
                        v
                Use Default Embedded Allowlist
                (Priority 3: Fallback)
                        |
                        v
                    Cache & Use

        |
        v
Start Periodic Update Goroutine
(if ALLOWLIST_UPDATE_INTERVAL > 0)
        |
        v
Application Running
(ReplaceImageName uses cached domain map)
```

## Components

### 1. Initialization (InitDomainMap)
- Called once at application startup (sync.Once)
- Reads configuration from environment variables
- Loads allowlist from configured sources
- Starts periodic update goroutine

### 2. Loading Sources (Priority Order)
1. **File**: `/etc/repimage/allowlist.txt` or custom path
2. **URL**: DaoCloud mirror or custom URL
3. **Embedded**: Built-in default allowlist

### 3. Caching
- Thread-safe cache (sync.RWMutex)
- In-memory storage
- Periodic refresh based on update interval
- No network calls during image replacement

### 4. Runtime Behavior
- ReplaceImageName() uses cached map (fast)
- No blocking on network I/O
- No panics on network errors
- Graceful fallback to embedded defaults

## Configuration Examples

### Build-time Only
```yaml
env:
  - name: ALLOWLIST_FILE
    value: "/etc/repimage/allowlist.txt"
  - name: ALLOWLIST_UPDATE_INTERVAL
    value: "0"
```

### Dynamic with Updates
```yaml
env:
  - name: ALLOWLIST_URL
    value: "https://my-mirror.com/allowlist.txt"
  - name: ALLOWLIST_UPDATE_INTERVAL
    value: "30m"
```

### ConfigMap Based
```yaml
env:
  - name: ALLOWLIST_FILE
    value: "/config/allowlist.txt"
volumes:
  - name: config
    configMap:
      name: repimage-allowlist
```

## Error Handling

1. **File not found** → Try URL
2. **URL unreachable** → Use embedded defaults
3. **Invalid format** → Skip invalid lines, use valid entries
4. **Network timeout** → Log warning, use cached/embedded data
5. **No valid sources** → Always have embedded defaults

## Benefits

✅ **Reliability**: Always operational, even without network
✅ **Performance**: No repeated network calls
✅ **Flexibility**: Multiple configuration methods
✅ **Maintainability**: Clear separation of concerns
✅ **Testability**: Network-independent testing
✅ **Security**: Support for private mirrors

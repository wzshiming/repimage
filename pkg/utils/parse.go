package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"k8s.io/klog/v2"
)

const (
	mirrorProxyUrL = "https://mirror.ghproxy.com/https://raw.githubusercontent.com/DaoCloud/public-image-mirror/main/domain.txt"
)

var (
	legacyDefaultDomain = "index.docker.io"
	defaultDomain       = "docker.io"
	officialRepoName    = "library"
	defaultTag          = "latest"
)

// Default embedded allowlist as fallback
var defaultAllowlist = map[string]string{
	"docker.io":       "m.daocloud.io/docker.io",
	"gcr.io":          "m.daocloud.io/gcr.io",
	"k8s.gcr.io":      "m.daocloud.io/k8s.gcr.io",
	"registry.k8s.io": "m.daocloud.io/registry.k8s.io",
	"ghcr.io":         "m.daocloud.io/ghcr.io",
	"quay.io":         "m.daocloud.io/quay.io",
}

// Global domain map with cache and mutex
var (
	cachedDomainMap     map[string]string
	domainMapMutex      sync.RWMutex
	lastUpdateTime      time.Time
	updateInterval      = 1 * time.Hour // Default update interval
	allowlistFilePath   string
	allowlistSourceURL  string
	initOnce            sync.Once
)

// InitDomainMap initializes the domain map with configuration
// This should be called at application startup
func InitDomainMap() {
	initOnce.Do(func() {
		// Read configuration from environment variables
		allowlistFilePath = os.Getenv("ALLOWLIST_FILE")
		allowlistSourceURL = os.Getenv("ALLOWLIST_URL")
		if allowlistSourceURL == "" {
			allowlistSourceURL = mirrorProxyUrL
		}

		// Parse update interval from environment
		if intervalStr := os.Getenv("ALLOWLIST_UPDATE_INTERVAL"); intervalStr != "" {
			if duration, err := time.ParseDuration(intervalStr); err == nil {
				updateInterval = duration
			} else {
				klog.Warningf("Invalid ALLOWLIST_UPDATE_INTERVAL: %v, using default 1h", err)
			}
		}

		// Initialize domain map
		loadDomainMap()

		// Start periodic update goroutine if interval > 0
		if updateInterval > 0 {
			go func() {
				ticker := time.NewTicker(updateInterval)
				defer ticker.Stop()
				for range ticker.C {
					loadDomainMap()
				}
			}()
		}
	})
}

// loadDomainMap loads the domain map from various sources
func loadDomainMap() {
	var newMap map[string]string

	// Try to load from file first (highest priority)
	if allowlistFilePath != "" {
		if dm, err := loadFromFile(allowlistFilePath); err == nil {
			newMap = dm
			klog.Infof("Loaded domain map from file: %s", allowlistFilePath)
		} else {
			klog.Warningf("Failed to load domain map from file %s: %v", allowlistFilePath, err)
		}
	}

	// If file loading failed, try to load from URL
	if newMap == nil && allowlistSourceURL != "" {
		if dm, err := loadFromURL(allowlistSourceURL); err == nil {
			newMap = dm
			klog.Infof("Loaded domain map from URL: %s", allowlistSourceURL)
		} else {
			klog.Warningf("Failed to load domain map from URL %s: %v", allowlistSourceURL, err)
		}
	}

	// If both failed, use default allowlist
	if newMap == nil {
		newMap = make(map[string]string)
		for k, v := range defaultAllowlist {
			newMap[k] = v
		}
		klog.Info("Using default embedded allowlist")
	}

	// Update cached domain map
	domainMapMutex.Lock()
	cachedDomainMap = newMap
	lastUpdateTime = time.Now()
	domainMapMutex.Unlock()

	klog.Infof("Domain map updated with %d entries at %v", len(newMap), lastUpdateTime)
}

// loadFromFile loads domain map from a file
func loadFromFile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseDomainMap(file)
}

// loadFromURL loads domain map from a URL
func loadFromURL(url string) (map[string]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	return parseDomainMap(resp.Body)
}

// parseDomainMap parses domain map from an io.Reader
func parseDomainMap(r io.Reader) (map[string]string, error) {
	domainMap := make(map[string]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		oldDomain := strings.TrimSpace(parts[0])
		newDomain := strings.TrimSpace(parts[1])
		domainMap[oldDomain] = newDomain
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return domainMap, nil
}

func ReplaceImageName(name string) string {
	domainMap := GetDomainMap()
	parts := strings.SplitN(name, "/", 3)
	switch len(parts) {
	case 1:
		if matchDomain(domainMap, defaultDomain) {
			return fmt.Sprintf("%s/%s/%s", domainMap[defaultDomain], officialRepoName, parts[0])
		}
	case 2:
		if matchDomain(domainMap, defaultDomain) {
			return fmt.Sprintf("%s/%s/%s", domainMap[defaultDomain], parts[0], parts[1])
		}
	case 3:
		if matchDomain(domainMap, parts[0]) {
			return fmt.Sprintf("%s/%s/%s", domainMap[parts[0]], parts[1], parts[2])
		}
	}
	return name
}

// GetDomainMap returns the cached domain map
// Initializes if not already initialized
func GetDomainMap() map[string]string {
	InitDomainMap()

	domainMapMutex.RLock()
	defer domainMapMutex.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]string)
	for k, v := range cachedDomainMap {
		result[k] = v
	}
	return result
}

func matchDomain(domainMap map[string]string, domain string) bool {
	_, ok := domainMap[domain]
	return ok
}

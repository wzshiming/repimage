package utils

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestInitDomainMap(t *testing.T) {
	// Reset for testing
	cachedDomainMap = nil
	initOnce = sync.Once{}
	
	// Set environment to disable periodic updates
	os.Setenv("ALLOWLIST_UPDATE_INTERVAL", "0")
	defer os.Unsetenv("ALLOWLIST_UPDATE_INTERVAL")

	// Call InitDomainMap
	InitDomainMap()

	// Verify domain map is initialized
	dm := GetDomainMap()
	if len(dm) == 0 {
		t.Error("Expected domain map to be initialized with default values")
	}

	// Verify default values are present
	if dm["docker.io"] != "m.daocloud.io/docker.io" {
		t.Errorf("Expected docker.io mapping, got: %v", dm["docker.io"])
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-allowlist.txt")

	content := `# Test allowlist
docker.io=test.mirror.io/docker.io
gcr.io=test.mirror.io/gcr.io

# Comment line
k8s.gcr.io=test.mirror.io/k8s.gcr.io
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load from file
	dm, err := loadFromFile(testFile)
	if err != nil {
		t.Fatalf("loadFromFile failed: %v", err)
	}

	// Verify mappings
	expected := map[string]string{
		"docker.io":  "test.mirror.io/docker.io",
		"gcr.io":     "test.mirror.io/gcr.io",
		"k8s.gcr.io": "test.mirror.io/k8s.gcr.io",
	}

	if len(dm) != len(expected) {
		t.Errorf("Expected %d mappings, got %d", len(expected), len(dm))
	}

	for k, v := range expected {
		if dm[k] != v {
			t.Errorf("Expected %s=%s, got %s=%s", k, v, k, dm[k])
		}
	}
}

func TestLoadFromFileWithConfiguration(t *testing.T) {
	// Reset for testing
	cachedDomainMap = nil
	initOnce = sync.Once{}
	
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "config-allowlist.txt")

	content := `docker.io=custom.mirror.io/docker.io
gcr.io=custom.mirror.io/gcr.io
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Set environment variable
	os.Setenv("ALLOWLIST_FILE", testFile)
	os.Setenv("ALLOWLIST_UPDATE_INTERVAL", "0")
	defer os.Unsetenv("ALLOWLIST_FILE")
	defer os.Unsetenv("ALLOWLIST_UPDATE_INTERVAL")

	// Initialize
	InitDomainMap()

	// Get domain map
	dm := GetDomainMap()

	// Verify custom mapping is used
	if dm["docker.io"] != "custom.mirror.io/docker.io" {
		t.Errorf("Expected custom mapping, got: %v", dm["docker.io"])
	}
}

func TestParseDomainMapWithVariousFormats(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected map[string]string
	}{
		{
			name: "standard format",
			content: `docker.io=m.daocloud.io/docker.io
gcr.io=m.daocloud.io/gcr.io`,
			expected: map[string]string{
				"docker.io": "m.daocloud.io/docker.io",
				"gcr.io":    "m.daocloud.io/gcr.io",
			},
		},
		{
			name: "with spaces",
			content: `docker.io = m.daocloud.io/docker.io
gcr.io = m.daocloud.io/gcr.io`,
			expected: map[string]string{
				"docker.io": "m.daocloud.io/docker.io",
				"gcr.io":    "m.daocloud.io/gcr.io",
			},
		},
		{
			name: "with comments and blank lines",
			content: `# Comment line
docker.io=m.daocloud.io/docker.io

# Another comment
gcr.io=m.daocloud.io/gcr.io

`,
			expected: map[string]string{
				"docker.io": "m.daocloud.io/docker.io",
				"gcr.io":    "m.daocloud.io/gcr.io",
			},
		},
		{
			name:    "invalid format - ignored",
			content: `invalid line without equals
docker.io=m.daocloud.io/docker.io
another invalid line`,
			expected: map[string]string{
				"docker.io": "m.daocloud.io/docker.io",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.txt")
			if err := os.WriteFile(testFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			dm, err := loadFromFile(testFile)
			if err != nil {
				t.Fatalf("loadFromFile failed: %v", err)
			}

			if len(dm) != len(tt.expected) {
				t.Errorf("Expected %d entries, got %d", len(tt.expected), len(dm))
			}

			for k, v := range tt.expected {
				if dm[k] != v {
					t.Errorf("Expected %s=%s, got %s=%s", k, v, k, dm[k])
				}
			}
		})
	}
}

func TestDefaultAllowlist(t *testing.T) {
	// Reset for testing
	cachedDomainMap = nil
	initOnce = sync.Once{}
	
	// Set environment to use only defaults (no file, no URL)
	os.Setenv("ALLOWLIST_UPDATE_INTERVAL", "0")
	os.Setenv("ALLOWLIST_FILE", "/nonexistent/file.txt")
	os.Setenv("ALLOWLIST_URL", "http://invalid.url.that.does.not.exist/allowlist.txt")
	defer os.Unsetenv("ALLOWLIST_UPDATE_INTERVAL")
	defer os.Unsetenv("ALLOWLIST_FILE")
	defer os.Unsetenv("ALLOWLIST_URL")

	// Initialize - should fall back to default
	InitDomainMap()

	dm := GetDomainMap()

	// Verify default mappings exist
	defaultMappings := []string{"docker.io", "gcr.io", "k8s.gcr.io", "registry.k8s.io", "ghcr.io", "quay.io"}
	for _, domain := range defaultMappings {
		if _, exists := dm[domain]; !exists {
			t.Errorf("Expected default mapping for %s", domain)
		}
	}
}

func TestUpdateIntervalConfiguration(t *testing.T) {
	// Reset for testing
	cachedDomainMap = nil
	initOnce = sync.Once{}
	updateInterval = 1 * time.Hour // reset to default
	
	// Set custom update interval
	os.Setenv("ALLOWLIST_UPDATE_INTERVAL", "30m")
	os.Setenv("ALLOWLIST_FILE", "/nonexistent/file.txt")
	os.Setenv("ALLOWLIST_URL", "http://invalid.url/allowlist.txt")
	defer os.Unsetenv("ALLOWLIST_UPDATE_INTERVAL")
	defer os.Unsetenv("ALLOWLIST_FILE")
	defer os.Unsetenv("ALLOWLIST_URL")

	// Initialize
	InitDomainMap()

	// Verify interval was set
	if updateInterval != 30*time.Minute {
		t.Errorf("Expected update interval 30m, got %v", updateInterval)
	}
}

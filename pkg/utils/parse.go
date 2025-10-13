package utils

import (
	"strings"
)

const (
	defaultDomain    = "docker.io"
	officialRepoName = "library"
)

// ReplaceImageName adds a mirror prefix to container image names
func ReplaceImageName(prefix, name string) string {
	parts := strings.SplitN(name, "/", 3)
	if parts[0] == prefix {
		return name
	}

	switch len(parts) {
	case 1:
		return strings.Join([]string{prefix, defaultDomain, officialRepoName, parts[0]}, "/")
	case 2:
		if !isDomain(parts[0]) {
			return strings.Join([]string{prefix, defaultDomain, parts[0], parts[1]}, "/")
		}

		if isLegacyDefaultDomain(parts[0]) {
			parts[0] = defaultDomain
		}

		return strings.Join([]string{prefix, parts[0], parts[1]}, "/")
	case 3:
		if !isDomain(parts[0]) {
			return strings.Join([]string{prefix, defaultDomain, parts[0], parts[1], parts[2]}, "/")
		}

		if isLegacyDefaultDomain(parts[0]) {
			parts[0] = defaultDomain
		}
		return strings.Join([]string{prefix, parts[0], parts[1], parts[2]}, "/")
	}
	return name
}

func isDomain(name string) bool {
	return strings.Contains(name, ".")
}

var (
	legacyDefaultDomain = map[string]struct{}{
		"index.docker.io":      {},
		"registry-1.docker.io": {},
	}
)

func isLegacyDefaultDomain(name string) bool {
	_, ok := legacyDefaultDomain[name]
	return ok
}

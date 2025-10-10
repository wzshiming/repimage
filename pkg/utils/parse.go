package utils

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

const (
	mirrorProxyUrL = "https://files.m.daocloud.io/raw.githubusercontent.com/DaoCloud/public-image-mirror/main/domain.txt"
)

var (
	legacyDefaultDomain = "index.docker.io"
	defaultDomain       = "docker.io"
	officialRepoName    = "library"
	defaultTag          = "latest"
)

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

func GetDomainMap() map[string]string {
	domainMap := make(map[string]string)
	res, err := http.Get(mirrorProxyUrL)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "=")
		if len(parts) != 2 {
			continue
		}
		oldDomain := parts[0][:len(parts[0])-1]
		NewDomain := parts[1][:len(parts[1])-1]
		domainMap[oldDomain] = NewDomain
	}
	return domainMap
}

func matchDomain(domainMap map[string]string, domain string) bool {
	//domainMap := GetDomainMap()
	_, ok := domainMap[domain]
	return ok
}

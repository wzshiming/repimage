package utils

import (
	"testing"
)

func TestReplaceImageNameWithDomainMap(t *testing.T) {
	// Mock domain map for testing
	domainMap := map[string]string{
		"docker.io":   "m.daocloud.io/docker.io",
		"k8s.gcr.io":  "m.daocloud.io/k8s.gcr.io",
		"gcr.io":      "m.daocloud.io/gcr.io",
		"ghcr.io":     "m.daocloud.io/ghcr.io",
		"quay.io":     "m.daocloud.io/quay.io",
		"registry.k8s.io": "m.daocloud.io/registry.k8s.io",
	}

	type args struct {
		name      string
		domainMap map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case1_simple_image_name",
			args: args{
				name:      "nginx",
				domainMap: domainMap,
			},
			want: "m.daocloud.io/docker.io/library/nginx",
		},
		{
			name: "case2_image_with_tag",
			args: args{
				name:      "nginx:v1.1.1",
				domainMap: domainMap,
			},
			want: "m.daocloud.io/docker.io/library/nginx:v1.1.1",
		},
		{
			name: "case3_user_repo_with_tag",
			args: args{
				name:      "hongshixing/nginx:v1.1.1",
				domainMap: domainMap,
			},
			want: "m.daocloud.io/docker.io/hongshixing/nginx:v1.1.1",
		},
		{
			name: "case4_gcr_image",
			args: args{
				name:      "k8s.gcr.io/hongshixing/nginx:v1.1.1",
				domainMap: domainMap,
			},
			want: "m.daocloud.io/k8s.gcr.io/hongshixing/nginx:v1.1.1",
		},
		{
			name: "case5_custom_domain_not_in_map",
			args: args{
				name:      "myit.fun/hongshixing/nginx:v1.1.1",
				domainMap: domainMap,
			},
			want: "myit.fun/hongshixing/nginx:v1.1.1",
		},
		{
			name: "case6_registry_k8s_io",
			args: args{
				name:      "registry.k8s.io/coredns/coredns:v1.8.4",
				domainMap: domainMap,
			},
			want: "m.daocloud.io/registry.k8s.io/coredns/coredns:v1.8.4",
		},
		{
			name: "case7_ghcr_io",
			args: args{
				name:      "ghcr.io/owner/repo:latest",
				domainMap: domainMap,
			},
			want: "m.daocloud.io/ghcr.io/owner/repo:latest",
		},
		{
			name: "case8_quay_io",
			args: args{
				name:      "quay.io/prometheus/node-exporter:latest",
				domainMap: domainMap,
			},
			want: "m.daocloud.io/quay.io/prometheus/node-exporter:latest",
		},
		{
			name: "case9_empty_domain_map",
			args: args{
				name:      "nginx",
				domainMap: map[string]string{},
			},
			want: "nginx",
		},
		{
			name: "case10_image_with_digest",
			args: args{
				name:      "nginx@sha256:abcd1234",
				domainMap: domainMap,
			},
			want: "m.daocloud.io/docker.io/library/nginx@sha256:abcd1234",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceImageNameWithDomainMap(tt.args.name, tt.args.domainMap); got != tt.want {
				t.Errorf("ReplaceImageNameWithDomainMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchDomain(t *testing.T) {
	domainMap := map[string]string{
		"docker.io":  "m.daocloud.io/docker.io",
		"k8s.gcr.io": "m.daocloud.io/k8s.gcr.io",
	}

	tests := []struct {
		name      string
		domainMap map[string]string
		domain    string
		want      bool
	}{
		{
			name:      "match_docker_io",
			domainMap: domainMap,
			domain:    "docker.io",
			want:      true,
		},
		{
			name:      "match_k8s_gcr_io",
			domainMap: domainMap,
			domain:    "k8s.gcr.io",
			want:      true,
		},
		{
			name:      "no_match_custom_domain",
			domainMap: domainMap,
			domain:    "custom.io",
			want:      false,
		},
		{
			name:      "empty_domain_map",
			domainMap: map[string]string{},
			domain:    "docker.io",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchDomain(tt.domainMap, tt.domain); got != tt.want {
				t.Errorf("matchDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

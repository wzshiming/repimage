package utils

import (
	"testing"
)

func TestReplaceImageName(t *testing.T) {
	type args struct {
		prefix        string
		ignoreDomains []string
		name          string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ignore docker.io - single name",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "nginx",
				ignoreDomains: []string{"docker.io"},
			},
			want: "docker.io/library/nginx",
		},
		{
			name: "ignore docker.io - with user",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "user/nginx:v1.1.1",
				ignoreDomains: []string{"docker.io"},
			},
			want: "docker.io/user/nginx:v1.1.1",
		},
		{
			name: "ignore k8s.gcr.io",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "k8s.gcr.io/user/nginx:v1.1.1",
				ignoreDomains: []string{"k8s.gcr.io"},
			},
			want: "k8s.gcr.io/user/nginx:v1.1.1",
		},
		{
			name: "ignore gcr.io but not k8s.gcr.io",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "k8s.gcr.io/user/nginx:v1.1.1",
				ignoreDomains: []string{"gcr.io"},
			},
			want: "m.daocloud.io/k8s.gcr.io/user/nginx:v1.1.1",
		},
		{
			name: "multiple ignore domains - match first",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "nginx",
				ignoreDomains: []string{"docker.io", "k8s.gcr.io"},
			},
			want: "docker.io/library/nginx",
		},
		{
			name: "multiple ignore domains - match second",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "k8s.gcr.io/nginx",
				ignoreDomains: []string{"docker.io", "k8s.gcr.io"},
			},
			want: "k8s.gcr.io/nginx",
		},
		{
			name: "multiple ignore domains - no match",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "quay.io/nginx",
				ignoreDomains: []string{"docker.io", "k8s.gcr.io"},
			},
			want: "m.daocloud.io/quay.io/nginx",
		},
		{
			name: "no ignore domains",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "nginx",
				ignoreDomains: nil,
			},
			want: "m.daocloud.io/docker.io/library/nginx",
		},
		{
			name: "empty ignore domains",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "nginx",
				ignoreDomains: []string{},
			},
			want: "m.daocloud.io/docker.io/library/nginx",
		},
		{
			name: "ignore custom domain",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "myregistry.example.com/myimage:latest",
				ignoreDomains: []string{"myregistry.example.com"},
			},
			want: "myregistry.example.com/myimage:latest",
		},
		{
			name: "don't ignore when domain doesn't match",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "nginx",
				ignoreDomains: []string{"quay.io"},
			},
			want: "m.daocloud.io/docker.io/library/nginx",
		},
		{
			name: "legacy default domain ignored via docker.io in ignore list",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "registry-1.docker.io/myuser/myimage:tag",
				ignoreDomains: []string{"docker.io"},
			},
			want: "docker.io/myuser/myimage:tag",
		},
		{
			name: "ignore docker.io - single name (no prefixing)",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "nginx",
				ignoreDomains: []string{"docker.io"},
			},
			want: "docker.io/library/nginx",
		},
		{
			name: "ignore docker.io - with user (no prefixing)",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "user/nginx:v1.1.1",
				ignoreDomains: []string{"docker.io"},
			},
			want: "docker.io/user/nginx:v1.1.1",
		},
		{
			name: "multiple ignore domains - match first (no prefixing)",
			args: args{
				prefix:        "m.daocloud.io",
				name:          "nginx",
				ignoreDomains: []string{"docker.io", "k8s.gcr.io"},
			},
			want: "docker.io/library/nginx",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceImageName(tt.args.prefix, tt.args.ignoreDomains, tt.args.name); got != tt.want {
				t.Errorf("ReplaceImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

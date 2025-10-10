package utils

import (
	"os"
	"testing"
)

func TestReplaceImageName(t *testing.T) {
	// Set up test environment - use default embedded allowlist
	os.Setenv("ALLOWLIST_UPDATE_INTERVAL", "0") // Disable periodic updates in tests
	defer os.Unsetenv("ALLOWLIST_UPDATE_INTERVAL")
	
	// Initialize domain map for testing
	InitDomainMap()
	
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case1",
			args: args{
				name: "nginx",
			},
			want: "m.daocloud.io/docker.io/library/nginx",
		},
		{
			name: "case2",
			args: args{
				name: "nginx:v1.1.1",
			},
			want: "m.daocloud.io/docker.io/library/nginx:v1.1.1",
		},
		{
			name: "case3",
			args: args{
				name: "hongshixing/nginx:v1.1.1",
			},
			want: "m.daocloud.io/docker.io/hongshixing/nginx:v1.1.1",
		},
		{
			name: "case4",
			args: args{
				name: "k8s.gcr.io/hongshixing/nginx:v1.1.1",
			},
			want: "m.daocloud.io/k8s.gcr.io/hongshixing/nginx:v1.1.1",
		},
		{
			name: "case5",
			args: args{
				name: "myit.fun/hongshixing/nginx:v1.1.1",
			},
			want: "myit.fun/hongshixing/nginx:v1.1.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceImageName(tt.args.name); got != tt.want {
				t.Errorf("ReplaceImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

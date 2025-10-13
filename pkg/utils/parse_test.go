package utils

import (
	"testing"
)

func TestReplaceImageName(t *testing.T) {
	type args struct {
		prefix string
		name   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case1",
			args: args{
				prefix: "m.daocloud.io",
				name:   "nginx",
			},
			want: "m.daocloud.io/docker.io/library/nginx",
		},
		{
			name: "case2",
			args: args{
				prefix: "m.daocloud.io",
				name:   "nginx:v1.1.1",
			},
			want: "m.daocloud.io/docker.io/library/nginx:v1.1.1",
		},
		{
			name: "case3",
			args: args{
				prefix: "m.daocloud.io",
				name:   "hongshixing/nginx:v1.1.1",
			},
			want: "m.daocloud.io/docker.io/hongshixing/nginx:v1.1.1",
		},
		{
			name: "case4",
			args: args{
				prefix: "m.daocloud.io",
				name:   "k8s.gcr.io/hongshixing/nginx:v1.1.1",
			},
			want: "m.daocloud.io/k8s.gcr.io/hongshixing/nginx:v1.1.1",
		},
		{
			name: "case5 - different prefix",
			args: args{
				prefix: "mirror.example.com",
				name:   "nginx",
			},
			want: "mirror.example.com/docker.io/library/nginx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceImageName(tt.args.prefix, tt.args.name); got != tt.want {
				t.Errorf("ReplaceImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

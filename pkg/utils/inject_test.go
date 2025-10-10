package utils

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInjectDockerInDockerContainer(t *testing.T) {
	tests := []struct {
		name                string
		deployment          *appsv1.Deployment
		wantContainerCount  int
		wantFirstContainer  string
		wantDindImage       string
	}{
		{
			name: "inject_into_deployment_with_one_container",
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "default",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "app",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			},
			wantContainerCount: 2,
			wantFirstContainer: "dind",
			wantDindImage:      "docker.io/library/docker:20.10.12-dind",
		},
		{
			name: "inject_into_deployment_with_multiple_containers",
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "default",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "app",
									Image: "nginx:latest",
								},
								{
									Name:  "sidecar",
									Image: "busybox:latest",
								},
							},
						},
					},
				},
			},
			wantContainerCount: 3,
			wantFirstContainer: "dind",
			wantDindImage:      "docker.io/library/docker:20.10.12-dind",
		},
		{
			name: "inject_into_empty_deployment",
			deployment: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-deployment",
					Namespace: "default",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{},
						},
					},
				},
			},
			wantContainerCount: 1,
			wantFirstContainer: "dind",
			wantDindImage:      "docker.io/library/docker:20.10.12-dind",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InjectDockerInDockerContainer(tt.deployment)
			
			if result == nil {
				t.Error("InjectDockerInDockerContainer() returned nil")
				return
			}

			containers := result.Spec.Template.Spec.Containers
			if len(containers) != tt.wantContainerCount {
				t.Errorf("InjectDockerInDockerContainer() container count = %v, want %v", len(containers), tt.wantContainerCount)
			}

			if len(containers) > 0 {
				if containers[0].Name != tt.wantFirstContainer {
					t.Errorf("InjectDockerInDockerContainer() first container name = %v, want %v", containers[0].Name, tt.wantFirstContainer)
				}
				if containers[0].Image != tt.wantDindImage {
					t.Errorf("InjectDockerInDockerContainer() dind image = %v, want %v", containers[0].Image, tt.wantDindImage)
				}
				if containers[0].SecurityContext == nil {
					t.Error("InjectDockerInDockerContainer() dind container should have SecurityContext")
				} else if containers[0].SecurityContext.Privileged == nil || !*containers[0].SecurityContext.Privileged {
					t.Error("InjectDockerInDockerContainer() dind container should be privileged")
				}
			}
		})
	}
}

func TestDindContainerConfiguration(t *testing.T) {
	// Test that the dind variable is properly configured
	if dind.Name != "dind" {
		t.Errorf("dind.Name = %v, want %v", dind.Name, "dind")
	}
	if dind.Image != "docker.io/library/docker:20.10.12-dind" {
		t.Errorf("dind.Image = %v, want %v", dind.Image, "docker.io/library/docker:20.10.12-dind")
	}
	if dind.SecurityContext == nil {
		t.Error("dind.SecurityContext is nil")
	} else if dind.SecurityContext.Privileged == nil {
		t.Error("dind.SecurityContext.Privileged is nil")
	} else if !*dind.SecurityContext.Privileged {
		t.Error("dind.SecurityContext.Privileged should be true")
	}
}

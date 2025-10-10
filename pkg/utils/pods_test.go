package utils

import (
	"encoding/json"
	"testing"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestAdmitPods_WrongResource(t *testing.T) {
	ar := admissionv1.AdmissionReview{
		Request: &admissionv1.AdmissionRequest{
			Resource: metav1.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "services", // Wrong resource type
			},
		},
	}

	resp := AdmitPods(ar)
	if resp == nil {
		t.Fatal("AdmitPods() returned nil")
	}
	if resp.Result == nil {
		t.Fatal("AdmitPods().Result is nil")
	}
	if resp.Result.Message == "" {
		t.Error("AdmitPods() should return error message for wrong resource type")
	}
}

func TestAdmitPods_ValidPod(t *testing.T) {
	// Create a simple pod
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:latest",
				},
				{
					Name:  "sidecar",
					Image: "k8s.gcr.io/coredns/coredns:v1.8.4",
				},
			},
		},
	}

	// Marshal pod to JSON
	podJSON, err := json.Marshal(pod)
	if err != nil {
		t.Fatalf("Failed to marshal pod: %v", err)
	}

	ar := admissionv1.AdmissionReview{
		Request: &admissionv1.AdmissionRequest{
			UID: "test-uid",
			Resource: metav1.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "pods",
			},
			Object: runtime.RawExtension{
				Raw: podJSON,
			},
		},
	}

	// Mock domain map for testing
	domainMap := map[string]string{
		"docker.io":  "m.daocloud.io/docker.io",
		"k8s.gcr.io": "m.daocloud.io/k8s.gcr.io",
	}

	resp := AdmitPodsWithDomainMap(ar, domainMap)
	if resp == nil {
		t.Fatal("AdmitPods() returned nil")
	}
	if !resp.Allowed {
		t.Error("AdmitPods() should allow valid pod")
	}
	if resp.Patch == nil {
		t.Error("AdmitPods() should return patch")
	}
	if resp.PatchType == nil {
		t.Error("AdmitPods() should return patch type")
	}
	if *resp.PatchType != admissionv1.PatchTypeJSONPatch {
		t.Errorf("AdmitPods() patch type = %v, want %v", *resp.PatchType, admissionv1.PatchTypeJSONPatch)
	}

	// Verify the patch contains the modified spec
	var patches []patchSpec
	if err := json.Unmarshal(resp.Patch, &patches); err != nil {
		t.Fatalf("Failed to unmarshal patch: %v", err)
	}
	if len(patches) != 1 {
		t.Errorf("AdmitPods() patch count = %v, want %v", len(patches), 1)
	}
	if len(patches) > 0 {
		if patches[0].Option != "replace" {
			t.Errorf("AdmitPods() patch option = %v, want %v", patches[0].Option, "replace")
		}
		if patches[0].Path != "/spec" {
			t.Errorf("AdmitPods() patch path = %v, want %v", patches[0].Path, "/spec")
		}
		// Verify images were replaced
		if len(patches[0].Value.Containers) != 2 {
			t.Errorf("AdmitPods() container count = %v, want %v", len(patches[0].Value.Containers), 2)
		} else {
			if patches[0].Value.Containers[0].Image != "m.daocloud.io/docker.io/library/nginx:latest" {
				t.Errorf("AdmitPods() first container image = %v, want %v", patches[0].Value.Containers[0].Image, "m.daocloud.io/docker.io/library/nginx:latest")
			}
			if patches[0].Value.Containers[1].Image != "m.daocloud.io/k8s.gcr.io/coredns/coredns:v1.8.4" {
				t.Errorf("AdmitPods() second container image = %v, want %v", patches[0].Value.Containers[1].Image, "m.daocloud.io/k8s.gcr.io/coredns/coredns:v1.8.4")
			}
		}
	}
}

func TestAdmitPods_InvalidJSON(t *testing.T) {
	ar := admissionv1.AdmissionReview{
		Request: &admissionv1.AdmissionRequest{
			Resource: metav1.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "pods",
			},
			Object: runtime.RawExtension{
				Raw: []byte(`{invalid json}`),
			},
		},
	}

	resp := AdmitPods(ar)
	if resp == nil {
		t.Fatal("AdmitPods() returned nil")
	}
	if resp.Result == nil {
		t.Fatal("AdmitPods().Result is nil")
	}
	if resp.Result.Message == "" {
		t.Error("AdmitPods() should return error message for invalid JSON")
	}
}

func TestPatchSpec_Marshaling(t *testing.T) {
	spec := patchSpec{
		Option: "replace",
		Path:   "/spec",
		Value: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "test",
					Image: "test:latest",
				},
			},
		},
	}

	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("Failed to marshal patchSpec: %v", err)
	}

	var unmarshaled patchSpec
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal patchSpec: %v", err)
	}

	if unmarshaled.Option != spec.Option {
		t.Errorf("Unmarshaled option = %v, want %v", unmarshaled.Option, spec.Option)
	}
	if unmarshaled.Path != spec.Path {
		t.Errorf("Unmarshaled path = %v, want %v", unmarshaled.Path, spec.Path)
	}
	if len(unmarshaled.Value.Containers) != len(spec.Value.Containers) {
		t.Errorf("Unmarshaled container count = %v, want %v", len(unmarshaled.Value.Containers), len(spec.Value.Containers))
	}
}

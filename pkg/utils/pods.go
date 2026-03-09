package utils

import (
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

const (
	// SkipImageRewriteAnnotation is the annotation key to skip image rewriting
	SkipImageRewriteAnnotation = "repimage.kubernetes.io/skip"
)

type patchSpec struct {
	Option string         `json:"op"`
	Path   string         `json:"path"`
	Value  corev1.PodSpec `json:"value"`
}

// AdmitPods processes admission review requests for pods and replaces container images, skipping ignored domains
func AdmitPods(prefix string, ignoreDomains []string, ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.Info("admitting pods...")
	podResource := metav1.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	if ar.Request.Resource != podResource {
		err := fmt.Errorf("expect resource to be %s", podResource)
		klog.Error(err)
		return ToAdmissionResponse(err)
	}
	raw := ar.Request.Object.Raw
	pod := &corev1.Pod{}
	deserializer := Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, pod); err != nil {
		klog.Error(err)
		return ToAdmissionResponse(err)
	}

	// Check if the pod has the annotation to skip image rewriting
	if pod.Annotations != nil {
		if skip, exists := pod.Annotations[SkipImageRewriteAnnotation]; exists && skip == "true" {
			klog.Info("pod has skip annotation, skipping image rewriting")
			reviewResponse := admissionv1.AdmissionResponse{}
			reviewResponse.Allowed = true
			return &reviewResponse
		}
	}

	reviewResponse := admissionv1.AdmissionResponse{}

	var updated bool

	for i, c := range pod.Spec.InitContainers {
		newImage := ReplaceImageName(prefix, ignoreDomains, c.Image)
		if newImage != c.Image {
			pod.Spec.InitContainers[i].Image = newImage
			updated = true
		}
	}

	for i, c := range pod.Spec.Containers {
		newImage := ReplaceImageName(prefix, ignoreDomains, c.Image)
		if newImage != c.Image {
			pod.Spec.Containers[i].Image = newImage
			updated = true
		}
	}

	reviewResponse.Allowed = true

	if !updated {
		return &reviewResponse
	}

	podSpec := []patchSpec{
		{
			Option: "replace",
			Path:   "/spec",
			Value:  pod.Spec,
		},
	}
	podSpecJSON, err := json.Marshal(podSpec)
	if err != nil {
		klog.Error(err)
		return ToAdmissionResponse(err)
	}

	reviewResponse.Patch = podSpecJSON
	jsonPatchType := admissionv1.PatchTypeJSONPatch
	reviewResponse.PatchType = &jsonPatchType

	return &reviewResponse
}

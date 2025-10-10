package utils

import (
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
)

type patchSpec struct {
	Option string         `json:"op"`
	Path   string         `json:"path"`
	Value  corev1.PodSpec `json:"value"`
}

func AdmitPods(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	return AdmitPodsWithDomainMap(ar, nil)
}

func AdmitPodsWithDomainMap(ar admissionv1.AdmissionReview, domainMap map[string]string) *admissionv1.AdmissionResponse {
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

	klog.Info(pod)

	// 定义准入规则
	reviewResponse := admissionv1.AdmissionResponse{}
	containers := pod.Spec.Containers
	for i, container := range containers {
		if domainMap == nil {
			containers[i].Image = ReplaceImageName(container.Image)
		} else {
			containers[i].Image = ReplaceImageNameWithDomainMap(container.Image, domainMap)
		}
	}

	klog.Info(pod.Spec.Containers)
	reviewResponse.Allowed = true
	podSpec := []patchSpec{
		{
			Option: "replace",
			Path:   "/spec",
			Value:  pod.Spec,
		},
	}
	podSpecJson, err := json.Marshal(podSpec)
	if err != nil {
		klog.Error(err)
		return ToAdmissionResponse(err)
	}

	klog.Info(string(podSpecJson))

	reviewResponse.Patch = podSpecJson
	jsonPatchType := admissionv1.PatchTypeJSONPatch
	reviewResponse.PatchType = &jsonPatchType

	return &reviewResponse
}

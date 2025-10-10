package utils

import (
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AdmitFunc 处理不同种类的请求
type AdmitFunc func(v1.AdmissionReview) *v1.AdmissionResponse

func ToAdmissionResponse(err error) *v1.AdmissionResponse {
	return &v1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

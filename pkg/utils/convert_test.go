package utils

import (
	"errors"
	"testing"

	"k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestToAdmissionResponse(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantMsg string
	}{
		{
			name:    "simple_error",
			err:     errors.New("test error"),
			wantMsg: "test error",
		},
		{
			name:    "empty_error_message",
			err:     errors.New(""),
			wantMsg: "",
		},
		{
			name:    "complex_error_message",
			err:     errors.New("failed to process: invalid resource type"),
			wantMsg: "failed to process: invalid resource type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToAdmissionResponse(tt.err)
			if got == nil {
				t.Error("ToAdmissionResponse() returned nil")
				return
			}
			if got.Result == nil {
				t.Error("ToAdmissionResponse().Result is nil")
				return
			}
			if got.Result.Message != tt.wantMsg {
				t.Errorf("ToAdmissionResponse().Result.Message = %v, want %v", got.Result.Message, tt.wantMsg)
			}
		})
	}
}

func TestAdmitFunc(t *testing.T) {
	// Test that AdmitFunc type works as expected
	var admitFunc AdmitFunc = func(ar v1.AdmissionReview) *v1.AdmissionResponse {
		return &v1.AdmissionResponse{
			Allowed: true,
			Result: &metav1.Status{
				Message: "test",
			},
		}
	}

	ar := v1.AdmissionReview{}
	resp := admitFunc(ar)
	
	if resp == nil {
		t.Error("AdmitFunc returned nil")
		return
	}
	if !resp.Allowed {
		t.Error("AdmitFunc response should be allowed")
	}
	if resp.Result.Message != "test" {
		t.Errorf("AdmitFunc response message = %v, want %v", resp.Result.Message, "test")
	}
}

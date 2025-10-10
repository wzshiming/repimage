package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/wzshiming/repimage/pkg/utils"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

var (
	cert = getEnvOrDefault("TLS_CERT_FILE", "./certs/serverCert.pem")
	key  = getEnvOrDefault("TLS_KEY_FILE", "./certs/serverKey.pem")
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func serve(w http.ResponseWriter, r *http.Request, admit utils.AdmitFunc) {
	klog.Infof("request URI: %s", r.RequestURI)
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		} else {
			klog.Errorf("failed to read request body: %v", err)
		}
	}

	klog.Infof("handling request: %s", string(body))

	reqAdmissionReview := v1.AdmissionReview{}
	resAdmissionReview := v1.AdmissionReview{TypeMeta: metav1.TypeMeta{
		Kind:       "AdmissionReview",
		APIVersion: "admission.k8s.io/v1",
	}}

	deserializer := utils.Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &reqAdmissionReview); err != nil {
		klog.Error(err)
		resAdmissionReview.Response = utils.ToAdmissionResponse(err)
	} else {
		resAdmissionReview.Response = admit(reqAdmissionReview)
	}

	resAdmissionReview.Response.UID = reqAdmissionReview.Request.UID

	klog.Infof("sending response: %v", resAdmissionReview.Response)

	respBytes, err := json.Marshal(resAdmissionReview)
	if err != nil {
		klog.Error(err)
		return
	}
	if _, err := w.Write(respBytes); err != nil {
		klog.Error(err)
	}
}

func servePods(w http.ResponseWriter, r *http.Request) {
	serve(w, r, utils.AdmitPods)
}

func main() {
	http.HandleFunc("/pods", servePods)
	klog.Infof("server starting with TLS cert: %s, key: %s", cert, key)
	if err := http.ListenAndServeTLS(":443", cert, key, nil); err != nil {
		klog.Exit(err)
	}
}

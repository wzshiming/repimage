package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/shixinghong/repimage/pkg/utils"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

var (
	cert = "./certs/serverCert.pem"
	key  = "./certs/serverKey.pem"
)

func serve(w http.ResponseWriter, r *http.Request, admit utils.AdmitFunc) {
	klog.Info(r.RequestURI)
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	klog.Info(fmt.Sprintf("handling request: %s", string(body)))

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

	klog.Info(fmt.Sprintf("sending response: %v", resAdmissionReview.Response))

	respBytes, err := json.Marshal(resAdmissionReview)
	if err != nil {
		klog.Error(err)
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
	klog.Info("server start")
	if err := http.ListenAndServeTLS(":443", cert, key, nil); err != nil {
		klog.Exit(err)
	}
}

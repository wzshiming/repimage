package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	"github.com/shixinghong/repimage/pkg/utils"
)

var (
	cert = "./certs/serverCert.pem"
	key  = "./certs/serverKey.pem"
)

func serve(w http.ResponseWriter, r *http.Request, admit utils.AdmitFunc) {
	klog.Info(r.RequestURI)
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	klog.Info(fmt.Sprintf("handling request: %s", string(body)))

	reqAdmissionReview := v1.AdmissionReview{}                          // 请求
	resAdmissionReview := v1.AdmissionReview{TypeMeta: metav1.TypeMeta{ // 响应
		Kind:       "AdmissionReview",
		APIVersion: "admission.k8s.io/v1",
	}}

	deserializer := utils.Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &reqAdmissionReview); err != nil {
		klog.Error(err)
		resAdmissionReview.Response = utils.ToAdmissionResponse(err)
	} else {
		// pass to admitFunc
		resAdmissionReview.Response = admit(reqAdmissionReview) // 业务逻辑
	}

	// 以下是固定写法
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
	// Initialize domain map at startup
	utils.InitDomainMap()
	
	http.HandleFunc("/pods", servePods)
	klog.Info("server start")
	if err := http.ListenAndServeTLS(":8080", cert, key, nil); err != nil {
		klog.Exit(err)
	}
}

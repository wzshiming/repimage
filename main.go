package main

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"strings"

	"github.com/wzshiming/repimage/pkg/utils"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

var (
	cert          = flag.String("cert", "./certs/serverCert.pem", "Path to TLS certificate file")
	key           = flag.String("key", "./certs/serverKey.pem", "Path to TLS key file")
	prefix        = flag.String("prefix", "m.daocloud.io", "Image mirror prefix")
	ignoreDomains = flag.String("ignore-domains", "", "Comma-separated list of domains to ignore (not replace)")
)

func serve(w http.ResponseWriter, r *http.Request, prefix string, ignoreDomains []string) {
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
		resAdmissionReview.Response = utils.AdmitPods(prefix, ignoreDomains, reqAdmissionReview)
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
	var domains []string
	if *ignoreDomains != "" {
		domains = strings.Split(*ignoreDomains, ",")
		// Trim whitespace from each domain
		for i := range domains {
			domains[i] = strings.TrimSpace(domains[i])
		}
	}
	serve(w, r, *prefix, domains)
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	http.HandleFunc("/pods", servePods)
	klog.Info("server start")
	if err := http.ListenAndServeTLS(":443", *cert, *key, nil); err != nil {
		klog.Exit(err)
	}
}

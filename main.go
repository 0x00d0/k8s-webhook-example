package main

import (
	"io/ioutil"
	"k8s-webhook-example/webhook"
	"k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/pods", func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.RequestURI)
		var body []byte
		if request.Body != nil {
			if data, err := ioutil.ReadAll(request.Body); err == nil {
				body = data
			}
		}

		requestAdmissionReview := v1.AdmissionReview{}
		responseAdmissionReview := v1.AdmissionReview{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AdmissionReview",
				APIVersion: "admission.k8s.io/v1",
			},
		}

		deserializer := webhook.Codecs.UniversalDeserializer()
		if _, _, err := deserializer.Decode(body, nil, &requestAdmissionReview); err != nil {
			klog.Error(err)
			responseAdmissionReview.Response = webhook.ToV1AdmissionResponse(err)
		} else {
			responseAdmissionReview.Response = webhook.AdmitPods(requestAdmissionReview)
		}
		responseAdmissionReview.Response.UID = requestAdmissionReview.Request.UID
		respBytes, _ := json.Marshal(responseAdmissionReview)

		writer.Write(respBytes)
	})

	//http.ListenAndServe(":8080", nil)
	tlsConfig := webhook.Config{
		CertFile: "/etc/webhook/certs/tls.crt",
		KeyFile:  "/etc/webhook/certs/tls.key",
	}
	server := &http.Server{
		Addr:      ":443",
		TLSConfig: webhook.ConfigTLS(tlsConfig),
	}

	server.ListenAndServeTLS("", "")
}

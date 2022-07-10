package main

import (
	"context"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"strings"
	"testing"
)

func TestPODWebHook(t *testing.T) {
	restConfig := &rest.Config{
		Host: "http://127.0.0.1:8080",
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatalln(err)
	}

	//body :=&v1.AdmissionReview{
	//	TypeMeta: metav1.TypeMeta{
	//		APIVersion: "admission.k8s.io/v1",
	//		Kind: "AdmissionReview",
	//	},
	//	Request: &v1.AdmissionRequest{
	//		UID: "705ab4f5-6393-11e8-b7cc-42010a800002",
	//		Kind: metav1.GroupVersionKind{
	//			Group: "",
	//			Version: "v1",
	//			Kind: "pods",
	//		},
	//		Resource: metav1.GroupVersionResource{
	//			Group: "",
	//			Version: "v1",
	//			Resource: "pods",
	//		},
	//		Name: "pod-example",
	//		Namespace: "default",
	//		Operation: "CREATE",
	//	},
	//}
	var body string = `
{
  "apiVersion": "admission.k8s.io/v1",
  "kind": "AdmissionReview",
  "request": {
    "uid": "705ab4f5-6393-11e8-b7cc-42010a800002",
    "kind": {"group":"","version":"v1","kind":"pods"},
    "resource": {"group":"","version":"v1","resource":"pods"},
    "name": "pod-example",
    "namespace": "default",
    "operation": "CREATE",
    "object": {"apiVersion":"v1","kind":"Pod","metadata":{"name":"pod-example","namespace":"default"}},
    "userInfo": {
      "username": "admin",
      "uid": "014fbff9a07c",
      "groups": ["system:authenticated","my-admin-group"],
      "extra": {
        "some-key":["some-value1", "some-value2"]
      }
    },
    "dryRun": false
  }
}
`
	result := clientSet.AdmissionregistrationV1().RESTClient().Post().Body(strings.NewReader(body)).Do(context.Background())
	b, _ := result.Raw()
	fmt.Println(string(b))
}

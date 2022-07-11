package webhook

import (
	"fmt"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

// only allow pods to pull images from specific registry.
func AdmitPods(ar v1.AdmissionReview) *v1.AdmissionResponse {
	klog.V(2).Info("admitting pods")
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

	if ar.Request.Resource != podResource {
		err := fmt.Errorf("expect resource to be %s", podResource)
		klog.Error(err)
		return ToV1AdmissionResponse(err)
	}
	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := Codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		klog.Error(err)
		return ToV1AdmissionResponse(err)
	}

	reviewResponse := v1.AdmissionResponse{}
	if pod.Name == "pod-example" {
		reviewResponse.Allowed = false
		reviewResponse.Result = &metav1.Status{Code: 403, Message: "pod name cannot be pod-example"}
	} else {
		reviewResponse.Allowed = true
		reviewResponse.Patch = patchImage()
		jsonPatchType := v1.PatchTypeJSONPatch
		reviewResponse.PatchType = &jsonPatchType
	}

	return &reviewResponse
}

func patchImage() []byte {
	str := `[
   {
		"op" : "replace" ,
		"path" : "/spec/containers/0/image" ,
		"value" : "nginx:1.19-alpine"
	},
   {
		"op" : "add" ,
		"path" : "/spec/initContainers" ,
		"value" : [{
						"name" : "myinit",
						"image" : "busybox:1.28",
 						"command" : ["sh", "-c", "echo The app is running!"]
 					 }]
	}

    
   
]`
	return []byte(str)
}

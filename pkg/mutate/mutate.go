package mutate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const AnnotationIntegrityMonitorInject = "integrity-monitor/inject"

func InjectIntegrityMonitor(logger *logrus.Logger, admReq *admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error) {
	// check if valid pod resource
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if admReq.Resource != podResource {
		return nil, fmt.Errorf("Receive unexpected resource type: %s", admReq.Resource.Resource)
	}

	admissionResponse := admissionv1.AdmissionResponse{
		Allowed: true,
	}

	// Decode the pod from the AdmissionReview.
	var pod corev1.Pod
	err := json.NewDecoder(bytes.NewReader(admReq.Object.Raw)).Decode(&pod)
	if err != nil {
		return nil, fmt.Errorf("error decoding raw pod: %w", err)
	}

	logEntry := logger.WithField("Pod", pod.Name)
	logEntry.WithField("Annotations", pod.Annotations).Debug("Process Pod")
	if value, ok := pod.Annotations[AnnotationIntegrityMonitorInject]; ok {
		inject, err := strconv.ParseBool(value)
		if err != nil {
			logEntry.WithError(err).Error("failed parse inject annotation value")
		}
		if inject {
			err := patchPod(&admissionResponse)
			if err != nil {
				logEntry.WithError(err).Error("failed patch pod")
			}
		}
	}
	return &admissionResponse, nil
}

func patchPod(admissionResponse *admissionv1.AdmissionResponse) error {
	file, err := os.Open("/app/patch-json-command.json")
	if err != nil {
		return fmt.Errorf("failed open patch file: %w", err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed read patch file: %w", err)
	}

	patchType := admissionv1.PatchTypeJSONPatch
	admissionResponse.PatchType = &patchType
	admissionResponse.Patch = data
	return nil
}

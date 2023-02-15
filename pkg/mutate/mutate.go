package mutate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	AnnotationIntegrityMonitorInject = "integrity-monitor.scnsoft.com/inject"
	AnnotationProcessName            = "integrity-monitor.scnsoft.com/process"
	AnnotationMonitoringPath         = "integrity-monitor.scnsoft.com/monitoring-path"
)

func InjectIntegrityMonitor(logger *logrus.Logger, admReq *admissionv1.AdmissionRequest) (*admissionv1.AdmissionResponse, error) {
	logger.Debug("processing request", admReq)
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
	if err := checkAnnotations(pod.GetAnnotations()); err != nil {
		logEntry.WithError(err).Error("failed parse inject annotation value")
		return &admissionResponse, nil
	}
	value := pod.GetAnnotations()[AnnotationIntegrityMonitorInject]
	inject, err := strconv.ParseBool(value)
	if err != nil {
		logEntry.WithError(err).Error("failed parse inject annotation value")
		return &admissionResponse, nil
	}
	if inject {
		sidecar := &SidecarConfig{}
		if err := sidecar.Load(viper.GetString("sidecar.cfg.file"), pod.GetAnnotations()); err != nil {
			logEntry.WithError(err).Error("failed loading sidecar config")
			return &admissionResponse, nil
		}
		patch, err := sidecar.CreatePatch(pod)
		if err != nil {
			logEntry.WithError(err).Error("failed creating patch")
			return &admissionResponse, nil
		}
		logEntry.Debugf("sidecar patches being applied for %v: patches: %v", pod.GetName(), patch)
		if err := patchPod(&admissionResponse, patch); err != nil {
			logEntry.WithError(err).Error("failed patching pod")
			return &admissionResponse, nil
		}
	}

	return &admissionResponse, nil
}

func patchPod(admissionResponse *admissionv1.AdmissionResponse, operation []PatchOperation) error {
	data, err := json.Marshal(&operation)
	if err != nil {
		return err
	}

	patchType := admissionv1.PatchTypeJSONPatch
	admissionResponse.PatchType = &patchType
	admissionResponse.Patch = data

	return nil
}

func checkAnnotations(annotations map[string]string) error {
	missedAnnotations := make([]string, 0)
	if _, ok := annotations[AnnotationIntegrityMonitorInject]; !ok {
		missedAnnotations = append(missedAnnotations, AnnotationIntegrityMonitorInject)
	}
	if _, ok := annotations[AnnotationMonitoringPath]; !ok {
		missedAnnotations = append(missedAnnotations, AnnotationMonitoringPath)
	}
	if _, ok := annotations[AnnotationProcessName]; !ok {
		missedAnnotations = append(missedAnnotations, AnnotationProcessName)
	}

	if len(missedAnnotations) > 0 {
		return fmt.Errorf("one ore more required annotations are missed %q", strings.Join(missedAnnotations, ","))
	}

	return nil
}

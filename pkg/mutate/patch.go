package mutate

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
)

// SidecarConfig for sidecar parameters.
type SidecarConfig struct {
	Volumes        []v1.Volume    `json:"volumes"`
	Containers     []v1.Container `json:"containers"`
	InitContainers []v1.Container `json:"initcontainers"`
}

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// Load loads common sidecar config from file, application specific config like monitoring process name and monitoring path
// from application pod annotations
func (sc *SidecarConfig) Load(configFile string, annotations map[string]string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, sc); err != nil {
		return err
	}

	sc.ConfigFromAnnotations(annotations)

	return nil
}

// CreatePatch creates mutation patch for pod
func (sc *SidecarConfig) CreatePatch(pod v1.Pod) ([]PatchOperation, error) {
	var patches []PatchOperation
	if sc != nil {
		patches = append(patches, addPatches(sc.InitContainers, pod.Spec.InitContainers, "/spec/initContainers")...)
		patches = append(patches, addPatches(sc.Containers, pod.Spec.Containers, "/spec/containers")...)
		patches = append(patches, addPatches(sc.Volumes, pod.Spec.Volumes, "/spec/volumes")...)
	}

	return patches, nil
}

func addPatches[T any](newCollection []T, existingCollection []T, path string) []PatchOperation {
	var patches []PatchOperation
	for index, item := range newCollection {
		indexPath := path
		var value interface{}
		first := index == 0 && len(existingCollection) == 0
		if !first {
			indexPath = indexPath + "/-"
			value = item

		} else {
			value = []T{item}
		}
		patches = append(patches, PatchOperation{
			Op:    "add",
			Path:  indexPath,
			Value: value,
		})
	}
	return patches
}

// ConfigFromAnnotations creates config from pod annotations
func (sc *SidecarConfig) ConfigFromAnnotations(annotations map[string]string) {
	for i := range sc.Containers {
		for k, v := range annotations {
			if k != AnnotationIntegrityMonitorInject {
				list := strings.Split(k, "/")
				sc.Containers[i].Args = append(sc.Containers[i].Args, fmt.Sprintf("--%s=%s", list[len(list)-1], v))
			}
		}
	}
}

package main

import (
	"encoding/json"
	"io"

	"github.com/imdario/mergo"

	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/yaml"
)

// KubeResource represents a resource kind and raw data to be used
// later for injecting EnvVars
type KubeResource struct {
	Kind string
	Data []byte
}

// ParseInlineDocs iterates through YAML or JSON docs and discovers
// their type returning a list of KubeResources
func ParseInlineDocs(reader io.Reader) ([]KubeResource, error) {
	resources := []KubeResource{}
	decoder := yaml.NewYAMLOrJSONDecoder(reader, 4096)

	for {
		rawExtension := runtime.RawExtension{}
		err := decoder.Decode(&rawExtension)
		if err == io.EOF {
			break
		} else if err != nil {
			return resources, err
		}

		kind, err := getResourceKind(rawExtension.Raw)
		if err != nil {
			return resources, err
		}

		resources = append(resources, KubeResource{
			Kind: kind,
			Data: rawExtension.Raw,
		})
	}

	return resources, nil
}

// ParseDoc reads a YAML or JSON doc and discovers the Kubernetes
// type returning a KubeResource
func ParseDoc(reader io.Reader) (KubeResource, error) {
	resource := KubeResource{}

	decoder := yaml.NewYAMLOrJSONDecoder(reader, 4096)
	rawExtension := runtime.RawExtension{}
	err := decoder.Decode(&rawExtension)
	if err != nil {
		return resource, err
	}

	resource.Kind, err = getResourceKind(rawExtension.Raw)
	if err != nil {
		return resource, err
	}

	resource.Data = rawExtension.Raw
	return resource, nil
}

// MergeDeployment ...
func (k *KubeResource) MergeDeployment(data []byte) (*v1beta1.Deployment, error) {
	target := &v1beta1.Deployment{}
	source := &v1beta1.Deployment{}
	err := mergeTargetSource(target, source, k.Data, data)
	return target, err
}

// MergeDaemonSet ...
func (k *KubeResource) MergeDaemonSet(data []byte) (*v1beta1.DaemonSet, error) {
	target := &v1beta1.DaemonSet{}
	source := &v1beta1.DaemonSet{}
	err := mergeTargetSource(target, source, k.Data, data)
	return target, err
}

// MergeReplicaSet ...
func (k *KubeResource) MergeReplicaSet(data []byte) (*v1beta1.ReplicaSet, error) {
	target := &v1beta1.ReplicaSet{}
	source := &v1beta1.ReplicaSet{}
	err := mergeTargetSource(target, source, k.Data, data)
	return target, err
}

// MergeRC ...
func (k *KubeResource) MergeRC(data []byte) (*v1.ReplicationController, error) {
	target := &v1.ReplicationController{}
	source := &v1.ReplicationController{}
	err := mergeTargetSource(target, source, k.Data, data)
	return target, err
}

// UnmarshalGeneric does not attempt to unmarshal to a known type,
// instead returns a generic interface object for displaying to the user
func (k *KubeResource) UnmarshalGeneric() (interface{}, error) {
	var generic interface{}
	err := json.Unmarshal(k.Data, &generic)
	return generic, err
}

// getResourceKind unmarshalls a file and returns the kind of resource doc
func getResourceKind(data []byte) (string, error) {
	typeMeta := unversioned.TypeMeta{}
	if err := json.Unmarshal(data, &typeMeta); err != nil {
		return "", err
	}
	return typeMeta.Kind, nil
}

func mergeTargetSource(t1 interface{}, t2 interface{}, d1 []byte, d2 []byte) error {
	if err := json.Unmarshal(d1, t1); err != nil {
		return err
	}

	if err := json.Unmarshal(d2, t2); err != nil {
		return err
	}

	err := mergo.Merge(t1, t2)
	return err
}

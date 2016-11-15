package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseInlineDocs(t *testing.T) {
	file, err := os.Open("fixtures/deployment-service.yml")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := ParseInlineDocs(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(resources) != 2 {
		t.Fatalf("Resource count is not 2: %+v", resources)
	}
}

func TestParseDoc(t *testing.T) {
	file, err := os.Open("fixtures/deployment.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	resource, err := ParseDoc(file)
	if err != nil {
		t.Fatal(err)
	}

	if resource.Kind != "Deployment" {
		t.Fatalf("Resource type is not Deployment: %+v", resource)
	}
}

func TestMergeDeployment(t *testing.T) {
	r, data := setupTestResources(t, "fixtures/deployment.json")
	result, err := r.MergeDeployment(data)
	if err != nil {
		t.Fatal(err)
	}
	if result.Spec.Template.Spec.ImagePullSecrets[0].Name != "registry" {
		t.Fatalf("result imagePullSecrets not set properly: %+v", result)
	}
}

func TestMergeDaemonSet(t *testing.T) {
	r, data := setupTestResources(t, "fixtures/daemonset.json")
	result, err := r.MergeDaemonSet(data)
	if err != nil {
		t.Fatal(err)
	}
	if result.Spec.Template.Spec.ImagePullSecrets[0].Name != "registry" {
		t.Fatalf("result imagePullSecrets not set properly: %+v", result)
	}
}

func TestMergeReplicaSet(t *testing.T) {
	r, data := setupTestResources(t, "fixtures/replicaset.json")
	result, err := r.MergeReplicaSet(data)
	if err != nil {
		t.Fatal(err)
	}
	if result.Spec.Template.Spec.ImagePullSecrets[0].Name != "registry" {
		t.Fatalf("result imagePullSecrets not set properly: %+v", result)
	}
}

func TestMergeReplicationController(t *testing.T) {
	r, data := setupTestResources(t, "fixtures/replicationcontroller.json")
	result, err := r.MergeRC(data)
	if err != nil {
		t.Fatal(err)
	}
	if result.Spec.Template.Spec.ImagePullSecrets[0].Name != "registry" {
		t.Fatalf("result imagePullSecrets not set properly: %+v", result)
	}
}

func TestUnmarshalGeneric(t *testing.T) {
	file, err := os.Open("fixtures/deployment.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	resources, err := ParseInlineDocs(file)
	if err != nil {
		t.Fatal(err)
	}

	_, err = resources[0].UnmarshalGeneric()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetResourceKind(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/deployment.json")
	if err != nil {
		t.Fatal(err)
	}

	kind, err := getResourceKind(data)
	if err != nil {
		t.Fatal(err)
	}

	if kind != "Deployment" {
		t.Fatalf("kind not equal")
	}
}

func setupTestResources(t *testing.T, path string) (KubeResource, []byte) {
	target, err := os.Open(path)
	defer target.Close()
	if err != nil {
		t.Fatal(err)
	}

	source, err := os.Open("fixtures/source.json")
	defer source.Close()
	if err != nil {
		t.Fatal(err)
	}

	targetResource, err := ParseDoc(target)
	if err != nil {
		t.Fatal(err)
	}

	sourceResource, err := ParseDoc(source)
	if err != nil {
		t.Fatal(err)
	}

	return targetResource, sourceResource.Data
}

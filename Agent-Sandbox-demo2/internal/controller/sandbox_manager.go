package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/dynamic"
)

type SandboxManager struct {
	DynamicClient dynamic.Interface
}

type InvestigationRequest struct {
	PodName   string `json:"podName"`
	Namespace string `json:"namespace"`
	Reason    string `json:"reason"`
	Describe  string `json:"describe"`
	Logs      string `json:"logs"`
	Events    string `json:"events"`
}

type InvestigationResponse struct {
	Diagnosis string `json:"diagnosis"`
}

func NewSandboxManager(
	client dynamic.Interface,
) *SandboxManager {

	return &SandboxManager{
		DynamicClient: client,
	}
}

func (s *SandboxManager) EnsureSandbox(
	name string,
	namespace string,
) error {

	gvr := schema.GroupVersionResource{
		Group:    "agents.x-k8s.io",
		Version:  "v1alpha1",
		Resource: "sandboxes",
	}

	_, err := s.DynamicClient.
		Resource(gvr).
		Namespace(namespace).
		Get(context.Background(), name, metav1.GetOptions{})

	if err == nil {
		return nil
	}

	sandbox := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "agents.x-k8s.io/v1alpha1",
			"kind":       "Sandbox",
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				"templateRef": map[string]interface{}{
					"name": "agent-template",
				},
			},
		},
	}

	_, err = s.DynamicClient.
		Resource(gvr).
		Namespace(namespace).
		Create(
			context.Background(),
			sandbox,
			metav1.CreateOptions{},
		)

	if err != nil {
		return fmt.Errorf(
			"failed creating sandbox: %w",
			err,
		)
	}

	return nil
}

func (s *SandboxManager) SendInvestigation(
	req InvestigationRequest,
) (string, error) {

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(
		"http://my-active-agent:8080/investigation",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var result InvestigationResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Diagnosis, nil
}
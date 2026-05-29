package controller

import (
	"context"
	"fmt"
	"log"

	"Agent-Sandbox-demo2/internal/diagnostic"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Reconciler struct {
	Clientset     *kubernetes.Clientset
	Collector     *diagnostic.Collector
	SandboxManager *SandboxManager
}

func NewReconciler(
	clientset *kubernetes.Clientset,
	sandboxManager *SandboxManager,
) *Reconciler {

	return &Reconciler{
		Clientset: clientset,
		Collector: diagnostic.NewCollector(clientset),
		SandboxManager: sandboxManager,
	}
}

func (r *Reconciler) ProcessPod(
	namespace string,
	podName string,
	reason string,
) error {

	log.Printf("Collecting logs for pod %s", podName)

	logs, err := r.Collector.CollectLogs(namespace, podName)

	if err != nil {
		log.Printf("Failed collecting logs: %v", err)
	}

	log.Printf("Finished collecting logs")

	log.Printf("Collecting events")

	events, err := r.Collector.CollectEvents(namespace, podName)

	if err != nil {
		log.Printf("Failed collecting events: %v", err)
	}

	log.Printf("Finished collecting events")

	log.Printf("Collecting describe")

	describe, err := r.Collector.CollectDescribe(namespace, podName)

	if err != nil {
		log.Printf("Failed collecting describe: %v", err)
	}

	log.Printf("Finished collecting describe")

	// adding optimization to limit the size of the data sent to the AI agent
	if len(logs) > 4000 {
		logs = logs[:4000]
	}

	if len(events) > 2000 {
		events = events[:2000]
	}

	if len(describe) > 4000 {
		describe = describe[:4000]
	}

	log.Printf("Sending investigation to AI agent")

	diagnosis, err := r.SandboxManager.SendInvestigation(
		InvestigationRequest{
			PodName:   podName,
			Namespace: namespace,
			Reason:    reason,
			Describe:  describe,
			Logs:      logs,
			Events:    events,
		},
	)

	if err != nil {
		return err
	}

	log.Printf("Received diagnosis from AI")

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-insight", podName),
		},
		Data: map[string]string{
			"failureReason": reason,
			"describe":      describe,
			"logs":          logs,
			"events":        events,
			"diagnosis":     diagnosis,
		},
	}

	log.Printf("Creating or updating ConfigMap")

	existingCM, err := r.Clientset.CoreV1().
		ConfigMaps(namespace).
		Get(
			context.Background(),
			cm.Name,
			metav1.GetOptions{},
		)

	if err == nil {

		existingCM.Data = cm.Data

		_, err = r.Clientset.CoreV1().
			ConfigMaps(namespace).
			Update(
				context.Background(),
				existingCM,
				metav1.UpdateOptions{},
			)

		if err != nil {
			return err
		}

		log.Printf("Updated ConfigMap %s", cm.Name)

	} else {

		_, err = r.Clientset.CoreV1().
			ConfigMaps(namespace).
			Create(
				context.Background(),
				cm,
				metav1.CreateOptions{},
			)

		if err != nil {
			return err
		}

		log.Printf("Created ConfigMap %s", cm.Name)
	}

	log.Printf("Finished processing pod %s", podName)

	return nil
}
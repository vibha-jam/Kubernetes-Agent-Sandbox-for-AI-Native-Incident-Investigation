package diagnostic

import (
	"bytes"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Collector struct {
	Clientset *kubernetes.Clientset
}

func NewCollector(clientset *kubernetes.Clientset) *Collector {
	return &Collector{
		Clientset: clientset,
	}
}

func (c *Collector) CollectLogs(namespace, podName string) (string, error) {
	req := c.Clientset.CoreV1().
		Pods(namespace).
		GetLogs(podName, &corev1.PodLogOptions{})

	stream, err := req.Stream(context.Background())
	if err != nil {
		return "", err
	}
	defer stream.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(stream)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c *Collector) CollectEvents(namespace, podName string) (string, error) {
	events, err := c.Clientset.CoreV1().
		Events(namespace).
		List(context.Background(), metav1.ListOptions{})

	if err != nil {
		return "", err
	}

	var output string

	for _, event := range events.Items {
		if event.InvolvedObject.Name == podName {
			output += fmt.Sprintf(
				"%s %s %s\n",
				event.Type,
				event.Reason,
				event.Message,
			)
		}
	}

	return output, nil
}

func (c *Collector) CollectDescribe(namespace, podName string) (string, error) {

	pod, err := c.Clientset.CoreV1().
		Pods(namespace).
		Get(context.Background(), podName, metav1.GetOptions{})

	if err != nil {
		return "", err
	}

	var output string

	output += fmt.Sprintf("Pod Name: %s\n", pod.Name)
	output += fmt.Sprintf("Namespace: %s\n", pod.Namespace)
	output += fmt.Sprintf("Node: %s\n", pod.Spec.NodeName)
	output += fmt.Sprintf("Phase: %s\n", pod.Status.Phase)

	output += "\n=== Container Statuses ===\n"

	for _, cs := range pod.Status.ContainerStatuses {

		output += fmt.Sprintf(
			"Container: %s\n",
			cs.Name,
		)

		output += fmt.Sprintf(
			"Ready: %v\n",
			cs.Ready,
		)

		output += fmt.Sprintf(
			"Restart Count: %d\n",
			cs.RestartCount,
		)

		if cs.State.Waiting != nil {

			output += fmt.Sprintf(
				"Current State: Waiting\n",
			)

			output += fmt.Sprintf(
				"Reason: %s\n",
				cs.State.Waiting.Reason,
			)

			output += fmt.Sprintf(
				"Message: %s\n",
				cs.State.Waiting.Message,
			)
		}

		if cs.State.Running != nil {

			output += "Current State: Running\n"
		}

		if cs.State.Terminated != nil {

			output += fmt.Sprintf(
				"Current State: Terminated\n",
			)

			output += fmt.Sprintf(
				"Exit Code: %d\n",
				cs.State.Terminated.ExitCode,
			)

			output += fmt.Sprintf(
				"Termination Reason: %s\n",
				cs.State.Terminated.Reason,
			)
		}

		if cs.LastTerminationState.Terminated != nil {

			last := cs.LastTerminationState.Terminated

			output += "\n--- Last Termination ---\n"

			output += fmt.Sprintf(
				"Exit Code: %d\n",
				last.ExitCode,
			)

			output += fmt.Sprintf(
				"Reason: %s\n",
				last.Reason,
			)

			output += fmt.Sprintf(
				"Finished At: %s\n",
				last.FinishedAt.String(),
			)
		}

		output += "\n"
	}

	output += "\n=== Conditions ===\n"

	for _, cond := range pod.Status.Conditions {

		output += fmt.Sprintf(
			"%s = %s\n",
			cond.Type,
			cond.Status,
		)
	}

	return output, nil
}
package controller

import (
	"context"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type PodWatcher struct {
	Clientset     *kubernetes.Clientset
	Reconciler    *Reconciler
	ProcessedPods map[string]bool
}

func NewPodWatcher(
	clientset *kubernetes.Clientset,
	dynamicClient dynamic.Interface,
) *PodWatcher {

	sandboxManager := NewSandboxManager(dynamicClient)

	return &PodWatcher{
		Clientset: clientset,
		Reconciler: NewReconciler(
			clientset,
			sandboxManager,
		),
		ProcessedPods: make(map[string]bool),
	}
}

func (w *PodWatcher) Watch() {

	for {

		pods, err := w.Clientset.CoreV1().
			Pods("demo-pods").
			List(context.Background(), metav1.ListOptions{})

		if err != nil {
			log.Println(err)
			time.Sleep(30 * time.Second)
			continue
		}

		for _, pod := range pods.Items {

			key := pod.Namespace + "/" + pod.Name

			if w.ProcessedPods[key] {
				continue
			}

			for _, cs := range pod.Status.ContainerStatuses {

				if cs.State.Waiting != nil {

					reason := cs.State.Waiting.Reason

					if reason == "CrashLoopBackOff" {

						log.Printf(
							"Detected CrashLoopBackOff in pod %s",
							pod.Name,
						)

						err := w.Reconciler.ProcessPod(
							pod.Namespace,
							pod.Name,
							reason,
						)

						if err != nil {
							log.Println(err)
							continue
						}

						w.ProcessedPods[key] = true
					}
				}

				if cs.LastTerminationState.Terminated != nil {

					if cs.LastTerminationState.Terminated.Reason == "OOMKilled" {

						log.Printf(
							"Detected OOMKilled in pod %s",
							pod.Name,
						)

						err := w.Reconciler.ProcessPod(
							pod.Namespace,
							pod.Name,
							"OOMKilled",
						)

						if err != nil {
							log.Println(err)
							continue
						}

						w.ProcessedPods[key] = true
					}
				}
			}
		}

		time.Sleep(30 * time.Second)
	}
}
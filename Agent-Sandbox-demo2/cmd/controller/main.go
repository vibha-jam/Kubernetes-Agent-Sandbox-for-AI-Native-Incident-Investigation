package main

import (
	"log"

	"Agent-Sandbox-demo2/internal/controller"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	watcher := controller.NewPodWatcher(
		clientset,
		dynamicClient,
	)

	log.Println("Starting AI investigator controller...")

	watcher.Watch()
}
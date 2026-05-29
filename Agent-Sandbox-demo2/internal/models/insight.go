package models

type ClusterInsight struct {
	PodName       string
	Namespace     string
	FailureReason string
	Logs          string
	Describe      string
	Events        string
	Diagnosis     string
}
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

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

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func buildPrompt(req InvestigationRequest) string {

	return fmt.Sprintf(`
You are a Kubernetes SRE assistant.

Analyze the following pod failure.

Pod Name:
%s

Namespace:
%s

Failure Reason:
%s

Describe Output:
%s

Events:
%s

Logs:
%s

Return:
1. probable root cause
2. confidence level
3. suggested next debugging step

Keep the answer concise.
`,
		req.PodName,
		req.Namespace,
		req.Reason,
		req.Describe,
		req.Events,
		req.Logs,
	)
}

func analyze(prompt string) (string, error) {

	log.Println("Preparing Ollama request")

	reqBody := OllamaRequest{
		Model:  "qwen3.5:2b",
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	start := time.Now()

	log.Println("Sending request to Ollama")

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return "", fmt.Errorf(
			"ollama returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	log.Printf("Ollama responded in %s", time.Since(start))

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Println("Received Ollama response body")

	var result OllamaResponse

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Printf("Failed decoding response: %s", string(bodyBytes))
		return "", err
	}

	log.Println("Successfully parsed Ollama response")

	return result.Response, nil
}

func writeHistory(
	req InvestigationRequest,
	diagnosis string,
) error {

	log.Println("Writing investigation history")

	timestamp := time.Now().Format("20060102-150405")

	filename := fmt.Sprintf(
		"/workspace/history/%s-%s.json",
		timestamp,
		req.PodName,
	)

	payload := map[string]interface{}{
		"timestamp": time.Now().String(),
		"podName":   req.PodName,
		"namespace": req.Namespace,
		"reason":    req.Reason,
		"describe":  req.Describe,
		"logs":      req.Logs,
		"events":    req.Events,
		"diagnosis": diagnosis,
	}

	data, err := json.MarshalIndent(
		payload,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	err = os.WriteFile(
		filename,
		data,
		0644,
	)

	if err != nil {
		return err
	}

	log.Println("Finished writing history")

	return nil
}

func investigationHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	log.Println("Received investigation request")

	var req InvestigationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Building prompt")

	prompt := buildPrompt(req)

	log.Printf("Prompt size: %d characters", len(prompt))

	log.Println("Starting analysis")

	diagnosis, err := analyze(prompt)
	if err != nil {
		log.Printf("Analysis failed: %v", err)

		resp := InvestigationResponse{
			Diagnosis: fmt.Sprintf(
				"analysis unavailable: %v",
				err,
			),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(resp)
		
		return
	}

	log.Println("Analysis completed")

	err = writeHistory(req, diagnosis)
	if err != nil {
		log.Printf("History write failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := InvestigationResponse{
		Diagnosis: diagnosis,
	}

	w.Header().Set("Content-Type", "application/json")

	log.Println("Returning response to controller")

	json.NewEncoder(w).Encode(resp)
}

func main() {

	http.HandleFunc(
		"/investigation",
		investigationHandler,
	)

	log.Println("Starting sandbox AI agent on :8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
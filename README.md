# AI-Powered Kubernetes Failure Investigator with Agent Sandboxes + gVisor

**An experimental Kubernetes-native AI investigation system that automatically detects failing Pods, gathers diagnostics, sends them to an isolated AI sandbox for analysis, and stores generated insights back into the cluster.**

# This project demonstrates:

- Automated pod failure detection
- AI-assisted Kubernetes troubleshooting
- Agent Sandbox warm pools
- Runtime persistence using PVCs
- gVisor sandbox isolation
- In-cluster inference using Ollama
- ConfigMap-based insight storage

## Architecture
```
Failing Pod
    ↓
AI Controller
    ↓
Collect Logs / Events / Describe
    ↓
SandboxWarmPool (gVisor isolated)
    ↓
Ollama + Local LLM
    ↓
AI Diagnosis Generated
    ↓
Stored back into Kubernetes ConfigMap
```

## Features 
### Automatic Failure Detection

For this MVP, the controller continuously watches Pods and detects:

- CrashLoopBackOff
- OOMKilled

## AI Investigation Pipeline

For every failing Pod, the system automatically collects:

- Pod logs
- Kubernetes events
- Pod describe output

## Failure metadata

- These diagnostics are sent to an isolated AI sandbox running Ollama.

- Sandbox Isolation with gVisor: the AI agent runs inside a gVisor sandbox using:

```
runtimeClassName: gvisor
```

### This provides:

- syscall interception
- userspace kernel isolation
- reduced kernel surface exposure
- stronger runtime isolation for AI agents

### Validation

Inside the gVisor sandbox:
```
ls /proc
```
showed:
```
sentry-meminfo
```
while multiple kernel-level /proc entries available in regular containers were absent.

This demonstrates gVisor’s virtualized kernel environment.

### Warm Pool Support

The project uses SandboxWarmPool to maintain pre-initialized AI agent sandboxes.

***Benefits:**

- reduced cold-start latency
- preloaded inference environment
- faster response during incidents

### Persistent Investigation Memory

- Investigation history is persisted using a PVC mounted into the sandbox:
```
mountPath: /workspace/history
```
- Generated investigation reports survive sandbox restarts and recreation.

## Repository Structure

```
.
├── cmd
│   ├── agent
│   └── controller
├── deploy
│   ├── demo
│   ├── kind
│   └── sandbox
├── internal
│   ├── controller
│   ├── diagnostic
│   └── models
├── history
├── Dockerfile.agent
├── Dockerfile.controller
├── Makefile
└── Readme.md
```
### Tech Stack
```
Go
Kubernetes
Ollama
gVisor
Kind
Agent Sandboxes
Persistent Volumes
Kubernetes Platform Engineering
Operatord/Controllers
```
## Demo Scenarios
### CrashLoopBackOff Detection
The controller detects a crashing Pod & Then generates a diagnosis.
```
CrashLoopBackOff
```
The AI sandbox analyzes:
- logs
- events
- container exit status

### OOMKilled Analysis

The controller detects:
```
OOMKilled
```
The AI agent identifies:

- memory exhaustion
- restart loops
- probable remediation steps

## Running the Project
### Prerequisites
- Docker
- Kind
- kubectl
- gVisor installed
- Agent Sandbox CRDs installed


**NOTE:** use Makefile to build, load, deploy

### Verify Sandbox Warm Pool
You should see this running before investigations begin.
```
agent-pool-xxxxx
```

### Verify Investigation Persistence
Generated investigation reports should persist across sandbox restarts.
```
kubectl exec -it <sandbox-pod> -- sh
ls /workspace/history
```

### Verify gVisor Isolation

**Inside the sandbox:**
```
ls /proc
```
**Observe:**
Indicates gVisor Sentry virtualization.
```
sentry-meminfo
```

## Future Improvements
- Informer watch based controller
- Multi-agent investigation workflows
- RAG-based cluster memory
- Distributed inference
- Policy-driven sandbox orchestration
- GPU-backed inference sandboxes
- Cross-cluster investigations

## Why This Project Exists

Modern Kubernetes incidents generate large volumes of operational signals.

**This project explores a future where:**

- AI agents autonomously investigate failures
- sandboxed runtimes safely isolate AI execution
- Kubernetes becomes self-diagnosing and self-healing, while maintaining strong runtime isolation boundaries.

## Demo Screenshots
### All pods
<img width="1603" height="382" alt="pods-all" src="https://github.com/user-attachments/assets/bbd5b41f-e1fe-4a47-8539-9275f6558388" />

### Containerized Ollama model qwen3.5:2b loading 
<img width="1904" height="346" alt="model-loading" src="https://github.com/user-attachments/assets/4183c6f5-6483-4a38-8c98-a7127a028bd7" />

### Agent Sandbox logs (agent-pool-xxxx)
<img width="1144" height="328" alt="agent-pool-logs" src="https://github.com/user-attachments/assets/e4e3a199-cc86-48cb-bfe3-74b522a12ce2" />
<img width="655" height="284" alt="agent-pool-logs2" src="https://github.com/user-attachments/assets/9eee297a-0684-4563-bd5a-799704433b6c" />
<img width="1915" height="765" alt="agent-pool-logs3" src="https://github.com/user-attachments/assets/19ea6a76-1b1f-4e6d-a90c-76d2abf9494b" />

### AI controller logs (ai-controller deployment)
<img width="782" height="359" alt="ai-controller-logs" src="https://github.com/user-attachments/assets/56a8eb35-8194-4599-97b2-eb26017e7611" />

### ConfigMaps (holds AI diagnositc insights)
<img width="656" height="143" alt="cm" src="https://github.com/user-attachments/assets/36c02ac2-343a-46f3-80e6-a76ae7e5fa3f" />

### AI CrashloopBackoff Pod Insight
<img width="1919" height="761" alt="crashloop-demo2" src="https://github.com/user-attachments/assets/b7d95564-544a-4dc8-a8fa-b2ab9ed57d88" />

### AI OOMKilled Pod Insight
<img width="1920" height="930" alt="oomkilled-cm" src="https://github.com/user-attachments/assets/fdd74c61-eabb-44e4-aeef-4e3bbcc20957" />

### Investigation Persistence with Agent Sandbox 
<img width="1762" height="247" alt="persistence" src="https://github.com/user-attachments/assets/de12af4b-14cb-4ea3-b526-bf1db69aee7f" />

### gVisor Isolation 
normal pod
<img width="1771" height="196" alt="normal-pod-proof" src="https://github.com/user-attachments/assets/82404c8c-8b0d-4600-9e18-450772068611" />
V/S
gVisor Isolated pos 
<img width="1298" height="178" alt="gvisor-proof" src="https://github.com/user-attachments/assets/978ca5b7-864e-4315-9183-aa587314af25" />

### Diagnosis Highlight
<img width="1916" height="468" alt="pod-cat-output" src="https://github.com/user-attachments/assets/3d56fdf7-6b01-4f50-9e8e-99825917e1f5" />


### License
MIT License

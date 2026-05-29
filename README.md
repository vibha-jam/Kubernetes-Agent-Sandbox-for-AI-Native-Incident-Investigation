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

ls /proc

showed:

sentry-meminfo

while multiple kernel-level /proc entries available in regular containers were absent.

This demonstrates gVisor’s virtualized kernel environment.

### Warm Pool Support

The project uses SandboxWarmPool to maintain pre-initialized AI agent sandboxes.

*** Benefits:**

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

Verify Sandbox Warm Pool
You should see this running before investigations begin.
```
agent-pool-xxxxx
```

Verify Investigation Persistence

Exec into sandbox:

kubectl exec -it <sandbox-pod> -- sh

Then:

ls /workspace/history

Generated investigation reports should persist across sandbox restarts.

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

### License
MIT License

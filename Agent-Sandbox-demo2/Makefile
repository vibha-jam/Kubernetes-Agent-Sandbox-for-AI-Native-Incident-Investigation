CLUSTER_NAME=agent-sandbox-dev

.PHONY: build-agent
build-agent:
	docker build -t ai-sandbox-agent:v1 \
	-f docker/agent/Dockerfile .

.PHONY: build-controller
build-controller:
	docker build -t ai-controller:v1 \
	-f docker/controller/Dockerfile .

.PHONY: load-agent
load-agent:
	kind load docker-image ai-sandbox-agent:v1 \
	--name $(CLUSTER_NAME)

.PHONY: load-controller
load-controller:
	kind load docker-image ai-controller:v1 \
	--name $(CLUSTER_NAME)

.PHONY: build
build: build-agent build-controller

.PHONY: load
load: load-agent load-controller

.PHONY: deploy-sandbox
deploy-sandbox:
	kubectl apply -f deploy/sandbox/

.PHONY: deploy-controller
deploy-controller:
	kubectl apply -f deploy/controller/

.PHONY: deploy-demo-pods
deploy-demo-pods:
	kubectl apply -f deploy/demo-pods/

.PHONY: deploy
deploy: deploy-sandbox deploy-controller deploy-demo-pods

.PHONY: restart-controller
restart-controller:
	kubectl rollout restart deploy/ai-controller

.PHONY: logs-controller
logs-controller:
	kubectl logs -f deploy/ai-controller

.PHONY: logs-agent
logs-agent:
	kubectl logs -f my-active-agent

.PHONY: clean-demo-pods
clean-demo-pods:
	kubectl delete -f deploy/demo-pods/ --ignore-not-found

.PHONY: clean-controller
clean-controller:
	kubectl delete -f deploy/controller/ --ignore-not-found

.PHONY: clean-sandbox
clean-sandbox:
	kubectl delete -f deploy/sandbox/ --ignore-not-found

.PHONY: clean
clean: clean-demo-pods clean-controller clean-sandbox
SHELL := /bin/bash

run:
	go run main.go

# build:
# 	go build -ldflags "-X main.build=local"

# Building containers

# $(shell git rev-parse --short HEAD)
VERSION := 1.0

all: service

service:
	docker build \
		-f zarf/docker/dockerfile \
		-t service-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.


# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all


# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := ardan-starter-cluster

# Upgrade to latest Kind (>=v0.11): e.g. brew upgrade kind
# For full Kind v0.11 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.0
# Kind release used for our project: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.1
# The image used below was copied by the above link and supports both amd64 and arm64.

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:fae9a58f17f18f06aeac9772ca8b5ac680ebbed985e266f711d936e91d113bad \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=service-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-service:
	kubectl get pods -o wide --watch --namespace=service-system

kind-load:
	kind load docker-image service-amd64:$(VERSION) --name $(KIND_CLUSTER)
	# cd zarf/k8s/kind/service-pod; kustomize edit set image service-api-image=service-api-amd64:$(VERSION)

kind-apply:
	kustomize build zarf/k8s/kind/service-pod | kubectl apply -f -

kind-restart:
	kubectl rollout restart deployment service-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-logs:
	kubectl logs -l app=service --all-containers=true -f --tail=100
	# kubectl logs -l app=service --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go

kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=service
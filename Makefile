# Image URL to use all building/pushing image targets
IMG ?= wozniakjan/cluster-cni-wh-hack:latest

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./main.go

# Run go fmt against code
fmt:
	go fmt ./main.go

# Run go vet against code
vet:
	go vet ./main.go

# Build the docker image
docker-build:
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

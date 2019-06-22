# Vars
STAGE_TAG=stage
PROD_TAG=v0.0.1
IMAGE_NAME=broker

# Accounts
DOCKERHUB_USER=n/a

# GKE
CLUSTER_STAGE=n/a
REGION=n/a
PROJECT=n/a

# Go
MAKE_CMD=make
# Go
GO_CMD=go

## Docker
DOCKER_CMD=docker

## Kubernetes
KUBECTL_CMD=kubectl

## Helm
HELM_CMD=helm

# Gcloud
GCLOUD_CMD=gcloud

# Misc
BINARY_NAME=broker
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	$(GO_BUILD) -o ./bin/$(BINARY_NAME) ./main.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(BINARY_UNIX) -v

test:
# Be sure to set up environment variables that apply for your case.
# PROVIDER_ID_KEY_1, PROVIDER_API_KEY_2, AWS_ACCESS_KEY_ID, AWS_SECRET_KEY
	$(GO_CMD) test -v ./...

clean:
	$(GO_CMD) clean
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)

deps:
	$(GO_CMD) get -u github.com/BurntSushi/toml
	$(GO_CMD) get -u github.com/cenkalti/backoff
  $(GO_CMD) get -u github.com/google/uuid
  $(GO_CMD) get -u github.com/mitchellh/mapstructure
  $(GO_CMD) get -u github.com/streadway/amqp
  $(GO_CMD) get -u gitlab.com/mikrowezel/backend/log

build-stage:
	$(MAKE_CMD) build
	$(DOCKER_CMD) login
  #$(DOCKER_CMD) build --no-cache -t $(DOCKERHUB_USER)/$(IMAGE_NAME):$(STAGE_TAG) .
	$(DOCKER_CMD) build --no-cache -t $(DOCKERHUB_USER)/$(IMAGE_NAME):$(STAGE_TAG) .
	$(DOCKER_CMD) push $(DOCKERHUB_USER)/$(IMAGE_NAME):$(STAGE_TAG)

connect-stage:
	$(GCLOUD_CMD) beta container clusters get-credentials $(CLUSTER_STAGE) --region $(REGION) --project $(PROJECT)

install-stage:
	$(MAKE_CMD) connect-stage
	$(HELM_INSTALL) --name $(IMAGE_NAME) -f ./deployments/helm/values-stage.yaml ./deployments/helm

delete-stage:
	$(MAKE_CMD) connect-stage
	$(HELM_DEL) --purge $(IMAGE_NAME)

deploy-stage:
	$(MAKE_CMD) build-stage
	$(MAKE_CMD) connect-stage
	$(MAKE_CMD) delete-stage
	$(HELM_INSTALL) --replace --name $(IMAGE_NAME) -f ./deployments/helm/values-stage.yaml ./deployments/helm


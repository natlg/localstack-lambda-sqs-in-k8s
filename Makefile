PROJECT_NAME := localstack-lambda-sqs-in-k8s
PROJECT_REPO := github.com/natlg/$(PROJECT_NAME)
GO111MODULE := on
PLATFORMS ?= linux_amd64
GO_STATIC_PACKAGES = $(GO_PROJECT)/cmd/analyzer $(GO_PROJECT)/cmd/publisher $(GO_PROJECT)/cmd/worker
GO_LDFLAGS += -X $(GO_PROJECT)/pkg/version.Version=$(VERSION)
DOCKER_REGISTRY = natlg
IMAGES = worker publisher analyzer provisioner

include build/makelib/common.mk
include build/makelib/output.mk
include build/makelib/golang.mk
include build/makelib/image.mk

GO_PROJECT = github.com/natlg/localstack-lambda-sqs-in-k8s
GO_BIN_DIR := $(abspath output/bin)
GO_IMG_DIR := $(abspath images)
GO_OUT_DIR := $(abspath output/images)

go.build:
	@echo start go build
	mkdir -p $(GO_OUT_DIR)
	 CGO_ENABLED=0 go build  -v -i -o $(GO_BIN_DIR)/publisher  $(GO_PROJECT)/cmd/publisher
	 CGO_ENABLED=0 go build  -v -i -o $(GO_BIN_DIR)/analyzer  $(GO_PROJECT)/cmd/analyzer
	 CGO_ENABLED=0 go build  -v -i -o $(GO_BIN_DIR)/worker  $(GO_PROJECT)/cmd/worker
	@echo go build is finished
	@echo start images build
	make -C $(GO_IMG_DIR)/worker PLATFORM=linux_amd64
	make -C $(GO_IMG_DIR)/publisher PLATFORM=linux_amd64
	make -C $(GO_IMG_DIR)/analyzer PLATFORM=linux_amd64
	make -C $(GO_IMG_DIR)/provisioner PLATFORM=linux_amd64
	@echo images build is finished



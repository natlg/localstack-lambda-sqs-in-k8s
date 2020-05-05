IMAGE = natlg/provision-localstack
GO_BIN_DIR := $(abspath ../../output)

img.build:
	echo docker build $(IMAGE)
	zip worker.zip -j $(GO_BIN_DIR)/bin/worker
	docker build \
		-t $(IMAGE) \
		-f Dockerfile \
		.
	echo done docker build $(IMAGE)




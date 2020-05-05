IMAGE = natlg/publisher
IMAGE_TEMP_DIR := $(abspath ../../output/images)
GO_BIN_DIR := $(abspath ../../output)

img.build:
	echo docker build $(IMAGE)
	cp $(GO_BIN_DIR)/bin/publisher $(IMAGE_TEMP_DIR)
	@docker build \
		-t $(IMAGE) \
		-f Dockerfile \
		$(IMAGE_TEMP_DIR)
	echo  done docker build $(IMAGE)




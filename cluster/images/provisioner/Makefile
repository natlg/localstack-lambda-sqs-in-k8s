PLATFORMS := linux_amd64
include ../../../build/makelib/common.mk

IMAGE = $(BUILD_REGISTRY)/provision-localstack-$(ARCH)
include ../../../build/makelib/image.mk

img.build:
	@$(INFO) docker build $(IMAGE)
	@cp Dockerfile $(IMAGE_TEMP_DIR) || $(FAIL)
	zip worker.zip -j $(OUTPUT_DIR)/bin/$(OS)_$(ARCH)/worker || $(FAIL)
	@cp worker.zip $(IMAGE_TEMP_DIR) || $(FAIL)
	@cp -r scripts $(IMAGE_TEMP_DIR) || $(FAIL)
	@cd $(IMAGE_TEMP_DIR) && $(SED_CMD) 's|BASEIMAGE|$(OSBASEIMAGE)|g' Dockerfile || $(FAIL)
	@docker build $(BUILD_ARGS) \
		--build-arg ARCH=$(ARCH) \
		--build-arg TINI_VERSION=$(TINI_VERSION) \
		-t $(IMAGE) \
		$(IMAGE_TEMP_DIR) || $(FAIL)
	@$(OK) docker build $(IMAGE)


.PHONY: ffmpeg
ffmpeg:
	@echo "Setting up FFmpeg..."
	@chmod +x ./hack/setup.sh
	@./hack/setup.sh darwin/universal


.PHONY: cleanup
cleanup:
	@echo "Teardown..."
	@rm -rf ./resources/

LAYERED_DIR := layered-architecture

.PHONY: help layered-build layered-run layered-clean

help:
	@echo Available targets:
	@echo   layered-build  Build the layered-architecture module
	@echo   layered-run    Run the layered-architecture demo
	@echo   layered-clean  Clean the layered-architecture Go build cache

layered-build:
	go -C $(LAYERED_DIR) build ./...

layered-run:
	go -C $(LAYERED_DIR) run ./cmd/quote-demo

layered-clean:
	go -C $(LAYERED_DIR) clean

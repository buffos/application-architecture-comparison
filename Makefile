LAYERED_DIR := layered-architecture
HEXAGONAL_DIR := hexagonal-architecture

.PHONY: help layered-build layered-run layered-clean

help:
	@echo Available targets:
	@echo   layered-build  Build the layered-architecture module
	@echo   layered-run    Run the layered-architecture demo
	@echo   layered-clean  Clean the layered-architecture Go build cache

layered-build:
	go -C $(LAYERED_DIR) build ./...

layered-test:
	go -C $(LAYERED_DIR) test ./...

layered-run:
	go -C $(LAYERED_DIR) run ./cmd/quote-demo

layered-clean:
	go -C $(LAYERED_DIR) clean

hex-build:
	go -C $(HEXAGONAL_DIR) build ./...

hex-test:
	go -C $(HEXAGONAL_DIR) test ./...

hex-run:
	go -C $(HEXAGONAL_DIR) run ./cmd/quote-demo

hex-clean:
	go -C $(HEXAGONAL_DIR) clean

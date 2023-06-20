.PHONY: build
build:
	go build ./cmd/publisher
	go build ./cmd/l0

.PHONY: run
run: build
	start ./lo.exe
	start ./publisher.exe


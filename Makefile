all: photo

.PHONY: photo
photo:
	go build -o bin/photo ./cmd/photo

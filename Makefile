all: photo

generate:
	go generate

.PHONY: photo
photo: generate
	go build -o bin/photo ./cmd/photo

all: photo posts

generate:
	go generate

clean:
	find . -name "*_easyjson.go" -delete

.PHONY: photo
photo: generate
	go build -o bin/photo ./cmd/photo

.PHONY: posts
posts: generate
	go build -o bin/posts ./cmd/posts

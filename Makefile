all: photo posts friends

GENERATE_FILES=$(shell fgrep -l -r "go:generate easyjson" $(PWD) | grep ".*\.go")

generate: clean
	easyjson -stubs -all $(GENERATE_FILES)
	easyjson -all $(GENERATE_FILES)

clean:
	find . -name "*_easyjson.go" -delete

.PHONY: photo
photo: generate
	go build -o bin/photo ./cmd/photo

.PHONY: posts
posts: generate
	go build -o bin/posts ./cmd/posts

.PHONY: friends
friends: generate
	go build -o bin/friends ./cmd/friends

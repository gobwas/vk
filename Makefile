all: photo posts friends

GENERATE_FILES=$(shell fgrep -l -r "go:generate easyjson" $(PWD) | grep ".*\.go")

generate: clean
	easyjson -stubs -all $(GENERATE_FILES)
	easyjson -all $(GENERATE_FILES)

clean:
	find . -name "*_easyjson.go" -delete

.PHONY: command
command:
ifeq ($(CMD),)
	$(error empty CMD variable)
else
	cp -r command/stub command/$(CMD)
	mv command/$(CMD)/stub.go command/$(CMD)/$(CMD).go
	find command/$(CMD) -type f -exec sed -i '' 's/stub/$(CMD)/g' {} \;
endif 

.PHONY: vk
vk: generate
	go build -o bin/vk ./cmd/vk

.PHONY: photo
photo: generate
	go build -o bin/photo ./cmd/photo

.PHONY: posts
posts: generate
	go build -o bin/posts ./cmd/posts

.PHONY: friends
friends: generate
	go build -o bin/friends ./cmd/friends

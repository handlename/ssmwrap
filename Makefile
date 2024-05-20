VERSION=$(shell git describe --tags --always --dirty="-dev")
PROJECT_USERNAME=handlename
PROJECT_REPONAME=ssmwrap
DIST_DIR=dist

export GO111MODULE := on

cmd/ssmwrap/ssmwrap: *.go */**/*.go
	CGO_ENABLED=0 go build -v -o $@ cmd/ssmwrap/main.go

test:
	go test -v ./...

.PHONY: build-docker-image
build-docker-image:
	docker build \
	  --rm \
	  --tag $(PROJECT_USERNAME)/$(PROJECT_REPONAME):$(VERSION) \
	  --tag ghcr.io/$(PROJECT_USERNAME)/$(PROJECT_REPONAME):$(VERSION) \
	  .

.PHONY: push-docker-image
push-docker-image:
	docker push ghcr.io/$(PROJECT_USERNAME)/$(PROJECT_REPONAME):$(VERSION)

clean:
	rm -rf cmd/ssmwrap/ssmwrap $(DIST_DIR)/*

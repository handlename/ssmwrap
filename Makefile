VERSION=${shell cat ./VERSION}
PROJECT_USERNAME=handlename
PROJECT_REPONAME=ssmwrap
DIST_DIR=dist

export GO111MODULE := on

cmd/ssmwrap/ssmwrap: *.go */**/*.go
	CGO_ENABLED=0 go build -v -ldflags '-X main.version=$(VERSION)' -o $@ cmd/ssmwrap/main.go

test:
	go test -v ./...

.PHONY: tag
tag:
	-git tag v$(VERSION)
	git push
	git push --tags

.PHONY: dist
dist: clean
	CGO_ENABLED=0 goxz \
	  -pv 'v$(VERSION)' \
	  -n ssmwrap \
	  -build-ldflags '-X main.version=$(VERSION)' \
	  -os='linux,darwin,windows' \
	  -arch='amd64,arm64' \
	  -d $(DIST_DIR) \
	  ./cmd/ssmwrap

.PHONY: upload
upload: dist
	mkdir -p $(DIST_DIR)
	ghr \
	  -u '$(PROJECT_USERNAME)' \
	  -r '$(PROJECT_REPONAME)' \
	  -prerelease \
	  -replace \
	  'v$(VERSION)' \
	  $(DIST_DIR)

.PHONY: build-docker-image
build-docker-image:
	docker build \
	  --rm \
	  --build-arg VERSION=$(VERSION) \
	  --tag $(PROJECT_USERNAME)/$(PROJECT_REPONAME):$(VERSION) \
	  --tag ghcr.io/$(PROJECT_USERNAME)/$(PROJECT_REPONAME):$(VERSION) \
	  .

.PHONY: push-docker-image
push-docker-image:
	docker push ghcr.io/$(PROJECT_USERNAME)/$(PROJECT_REPONAME):$(VERSION)

clean:
	rm -rf cmd/ssmwrap/ssmwrap $(DIST_DIR)/*

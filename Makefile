VERSION=${shell cat ./VERSION}
PROJECT_USERNAME=handlename
PROJECT_REPONAME=ssmwrap
DIST_DIR=dist

export GO111MODULE := on

cmd/ssmwrap/ssmwrap: *.go */**/*.go
	CGO_ENABLED=0 go build -v -ldflags '-X main.version=$(VERSION)' -o $@ cmd/ssmwrap/main.go

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
	  -arch='386,amd64' \
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

clean:
	rm -rf cmd/ssmwrap/ssmwrap $(DIST_DIR)/*

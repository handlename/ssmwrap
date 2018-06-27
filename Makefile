VERSION=${shell cat ./VERSION}
PROJECT_USERNAME=handlename
PROJECT_REPONAME=ssmwrap

cmd/ssmwrap/ssmwrap: *.go */**/*.go
	go build -v -ldflags '-X main.version=$(VERSION)' -o $@ cmd/ssmwrap/main.go

.PHONY: tag
tag:
	-git tag v$(VERSION)
	git push
	git push --tags

.PHONY: dist
dist:
	goxz \
	  -pv 'v$(VERSION)' \
	  -n ssmwrap \
	  -build-ldflags '-X main.version=$(VERSION)' \
	  -os='linux,darwin,windows' \
	  -arch='386,amd64' \
	  -d dist \
	  ./cmd/ssmwrap

.PHONY: upload
upload: dist
	ghr \
	  -u '$(PROJECT_USERNAME)' \
	  -r '$(PROJECT_REPONAME)' \
	  -prerelease \
	  -replace \
	  'v$(VERSION)' \
	  dist

clean:
	rm -rf cmd/ssmwrap/ssmwrap dist/*

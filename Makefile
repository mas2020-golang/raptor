export GIT_COMMIT=$(shell git rev-list -1 --abbrev-commit HEAD)

testing:
	@echo "==> cryptex test..."
	@go test ./... -v

goreleaser:
	@echo "start building..."
	@goreleaser  --rm-dist --snapshot --skip-publish
	@echo "done!"

install-on-mac: build test
	@echo "start install..."
	@echo "copying into $(GOPATH)/bin..."
	@cp bin/cryptex-darwin-amd64 $(GOPATH)/bin/cryptex
	@echo "done!"

run:
	clear
	go run main.go

build:
	# compile Go-AL for several platform
	@echo "oompiling for every OS and Platform..."
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" -o bin/cryptex-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" -o bin/cryptex-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" -o bin/cryptex-windows-amd64.exe main.go
	@echo "done!"

clean:
	@rm -rf bin
	@rm -rf dist

look_update_pkgs:
	# take a look at the newer versions of dependency modules
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null
export GIT_COMMIT=$(shell git rev-list -1 --abbrev-commit HEAD)

testing:
	@echo "==> raptor test..."
	@go test ./... -v

goreleaser:
	@echo "start building..."
	@goreleaser  --rm-dist --snapshot --skip-publish
	@echo "done!"

install-on-mac: build testing
	@echo "start installing..."
	@echo "copying into $(GOPATH)/bin..."
	@cp bin/raptor-darwin-amd64 $(GOPATH)/bin/raptor
	@echo "done!"

install-on-linux: testing
	@echo "start installing..."
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" -o ./bin/raptor-linux-amd64 main.go
	@echo "copying into $(GOPATH)/bin..."
	@cp ./bin/raptor-linux-amd64 $(GOPATH)/bin/raptor
	@echo "done!"

run:
	clear
	go run main.go

build:
	# compile Go-AL for several platform
	@echo "compiling for every OS and Platform..."
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" -o bin/raptor-darwin-amd64 main.go
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" -o bin/raptor-linux-amd64 main.go
	# GOOS=windows GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GIT_COMMIT}" -o bin/raptor-windows-amd64.exe main.go
	@echo "done!"

clean:
	@rm -rf bin
	@rm -rf dist

look_update_pkgs:
	# take a look at the newer versions of dependency modules
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null
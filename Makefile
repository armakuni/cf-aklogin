build:
	@go build

deps:
	@glide i

test:
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic

install:
	@cf install-plugin -f cf-aklogin

run:
	@cf aklogin -l

release:
	@${PWD}/bin/create-release

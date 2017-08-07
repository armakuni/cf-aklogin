build:
	@go build

deps:
	@go get -v github.com/Masterminds/glide
	@glide i

test:
	@go vet
	@go test -v -race -coverprofile=coverage.out -covermode=atomic

install:
	@cf install-plugin -f cf-aklogin

run:
	@cf aklogin -l

release:
	@${PWD}/bin/create-release

coverage:
	@go tool cover -html=coverage.out

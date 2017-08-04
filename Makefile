build:
	@go build

test:
	@go test

install:
	@cf install-plugin -f cf-aklogin

run:
	@cf aklogin -l

release:
	@${PWD}/bin/create-release

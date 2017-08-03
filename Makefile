PLUGIN_NAME=cf-aklogin

build:
	@go build

test:
	@go test

install:
	@cf install-plugin -f ${PLUGIN_NAME}

run:
	@cf aklogin -l

release:
	@GOOS=linux GOARCH=amd64 go build -o ${PLUGIN_NAME}.linux64
	@GOOS=linux GOARCH=386 go build -o ${PLUGIN_NAME}.linux32
	@GOOS=windows GOARCH=amd64 go build -o ${PLUGIN_NAME}.win64.exe
	@GOOS=windows GOARCH=386 go build -o ${PLUGIN_NAME}.win32.exe
	@GOOS=darwin GOARCH=amd64 go build -o ${PLUGIN_NAME}.darwin
	@shasum -a 1 ${PLUGIN_NAME}.linux64
	@shasum -a 1 ${PLUGIN_NAME}.linux32
	@shasum -a 1 ${PLUGIN_NAME}.win64.exe
	@shasum -a 1 ${PLUGIN_NAME}.win32.exe
	@shasum -a 1 ${PLUGIN_NAME}.darwin

CURDIR := "$(shell pwd)"

all: service randomservice isprimeservice concatenateservice serviceuser registryserver signalHandler menu

service:
	export GOPATH=${CURDIR}; \
	go build github.com/jzipfler/HTW-SwArchitektur/service

randomservice:
	export GOPATH=${CURDIR}; \
	go install github.com/jzipfler/HTW-SwArchitektur/randomservice

isprimeservice:
	export GOPATH=${CURDIR}; \
	go install github.com/jzipfler/HTW-SwArchitektur/isprimeservice

concatenateservice:
	export GOPATH=${CURDIR}; \
	go install github.com/jzipfler/HTW-SwArchitektur/concatenateservice
	
serviceuser:
	export GOPATH=${CURDIR}; \
	go install github.com/jzipfler/HTW-SwArchitektur/serviceuser

registryserver:
	export GOPATH=${CURDIR}; \
	go install github.com/jzipfler/HTW-SwArchitektur/registryserver

signalHandler:
	export GOPATH=${CURDIR}; \
	go build github.com/jzipfler/HTW-SwArchitektur/signalHandler

menu:
	export GOPATH=${CURDIR}; \
	go install github.com/jzipfler/HTW-SwArchitektur/menu
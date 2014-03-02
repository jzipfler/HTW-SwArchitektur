CURDIR := "$(shell pwd)"

all: service randomservice serviceuser registryserver

service:
	export GOPATH=${CURDIR}; \
	go build github.com/jzipfler/HTW-SwArchitektur/service

randomservice:
	export GOPATH=${CURDIR}; \
	go install github.com/jzipfler/HTW-SwArchitektur/randomservice

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
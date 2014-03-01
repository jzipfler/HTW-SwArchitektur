all: service randomservice serviceuser registryserver

service:
	go build github.com/jzipfler/HTW-SwArchitektur/service

randomservice:
	go install github.com/jzipfler/HTW-SwArchitektur/randomservice

serviceuser:
	go install github.com/jzipfler/HTW-SwArchitektur/serviceuser

registryserver:
	go install github.com/jzipfler/HTW-SwArchitektur/registryserver

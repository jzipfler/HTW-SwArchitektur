@echo off
echo buildung
set GOPATH=%CD%
go build github.com/jzipfler/HTW-SwArchitektur/service
go install github.com/jzipfler/HTW-SwArchitektur/randomservice
go install github.com/jzipfler/HTW-SwArchitektur/serviceuser
go install github.com/jzipfler/HTW-SwArchitektur/registryserver
echo done
pause

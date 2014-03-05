@echo off
echo buildung
set GOPATH=%CD%
go build github.com/jzipfler/HTW-SwArchitektur/service
go install github.com/jzipfler/HTW-SwArchitektur/randomservice
go install github.com/jzipfler/HTW-SwArchitektur/isprimeservice
go install github.com/jzipfler/HTW-SwArchitektur/concatenateservice
go install github.com/jzipfler/HTW-SwArchitektur/serviceuser
go install github.com/jzipfler/HTW-SwArchitektur/registryserver
go build github.com/jzipfler/HTW-SwArchitektur/signalHandler
go install github.com/jzipfler/HTW-SwArchitektur/menu
echo done

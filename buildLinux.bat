set GOOS=linux
set GOARCH=amd64
cd cmd\smallgopher-server
go build
cd ..\..
cd cmd\smallgopher-worker
go build
cd ..\..
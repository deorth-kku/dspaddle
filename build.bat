go build -o dspaddle.exe -trimpath -ldflags "-s -w -buildid= "
go build -buildmode=c-shared -tags dll -o dspaddle.dll -trimpath -ldflags "-s -w -buildid= "
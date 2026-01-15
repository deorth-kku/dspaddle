go1.26rc1 build -o dspaddle.exe -trimpath -ldflags "-s -w -buildid= "
go1.26rc1 build -buildmode=c-shared -tags dll -o dspaddle.dll -trimpath -ldflags "-s -w -buildid= "
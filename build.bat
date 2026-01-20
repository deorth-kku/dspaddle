go1.26rc2 build -o dspaddle.exe -trimpath -ldflags "-s -w -buildid= "
go1.26rc2 build -buildmode=c-shared -tags dll -o dspaddle.dll -trimpath -ldflags "-s -w -buildid= "
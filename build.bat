call env.bat

set GOOS=linux

go build -o o2api cmd/main.go

set GOOS=windows

go build -o o2api.exe cmd/main.go

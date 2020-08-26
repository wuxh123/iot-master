call env.bat

set GOOS=linux

go build -o dtu-admin main.go

set GOOS=windows

go build -o dtu-admin.exe main.go

@echo off

:: 编译前端（编译结束后会退出中止批处理，待解）
::cd portal
::ng build --prod
::cd ..

:: 把前端转化成go
:: go get -u github.com/UnnoTed/fileb0x
fileb0x b0x.yaml

:: 整体编译
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=*.gitlab.com,*.gitee.com
go env -w GOSUMDB=off


set GOOS=linux
go build -o dtu-admin main.go

set GOOS=windows
go build -o dtu-admin.exe main.go

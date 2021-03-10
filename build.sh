# 把前端转化成go
# go get -u github.com/UnnoTed/fileb0x
# fileb0x b0x.yaml

# 整体编译
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=*.gitlab.com,*.gitee.com
go env -w GOSUMDB=off


export GOOS=linux
#go build -o MyDTU main.go
go build -ldflags "-X 'args.goVersion=$(go version)' -X 'args.gitHash=$(git show -s --format=%H)' -X 'args.buildTime=$(git show -s --format=%cd)'" -o mydtu main.go

export GOOS=windows
#go build -o MyDTU.exe main.go
go build -ldflags "-X 'args.goVersion=$(go version)' -X 'args.gitHash=$(git show -s --format=%H)' -X 'args.buildTime=$(git show -s --format=%cd)'" -o mydtu.exe main.go

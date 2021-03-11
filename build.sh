# 把前端转化成go
# go get -u github.com/UnnoTed/fileb0x
# fileb0x b0x.yaml

# 整体编译
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOPRIVATE=*.gitlab.com,*.gitee.com
go env -w GOSUMDB=off

version="1.0.0"
read -t 5 -p "please input version(default:$version)" ver
if [ -n "${ver}" ];then
	version=$ver
fi


goVersion=$(go version | awk '{print $3}')
gitHash=$(git show -s --format=%H)
buildTime=$(date -d today +"%Y-%m-%d %H:%M:%S")

ldflags="-X 'mydtu/args.Version=$version' \
-X 'mydtu/args.goVersion=$goVersion' \
-X 'mydtu/args.gitHash=$gitHash' \
-X 'mydtu/args.buildTime=$buildTime'"

export GOOS=linux
#go build -o MyDTU main.go
go build -ldflags "$ldflags" -o mydtu-linux main.go

export GOOS=windows
#go build -o MyDTU.exe main.go
go build -ldflags "$ldflags" -o mydtu-win64.exe main.go

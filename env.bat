

go env -w GOPROXY=https://goproxy.cn,direct

go env -w GOPRIVATE=*.gitlab.com,*.gitee.com

go env -w GOSUMDB=off
:: go env -w GOSUMDB="sum.golang.google.cn"
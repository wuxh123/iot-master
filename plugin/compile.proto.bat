
:: 1 下载protoc https://github.com/google/protobuf/releases
:: 2 获取编译插件go get -u github.com/golang/protobuf/protoc-gen-go


:: 3 编译protoc --go_out=. *.proto

protoc  --go_out=plugins=grpc:. *.proto

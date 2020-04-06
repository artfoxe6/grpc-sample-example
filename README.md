##### 要求
1. go版本大于1.13
2. 已经下载 protoc 编译器并放到环境变量中
    >[protoc](https://github.com/protocolbuffers/protobuf/releases)
3. 已经下载 protoc-gen-go 并发到环境变量中
    >go get -u github.com/golang/protobuf/protoc-gen-go
    >
    >运行上面命令后 protoc-gen-go会下载到 GOBIN环境变量中，默认也就是$GOPATH/bin/
4. 在server目录中运行 
    > go run server.go
5. 在client中运行
    > go run client.go

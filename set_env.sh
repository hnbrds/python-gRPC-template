#!/bin/bash
# set go ENVPATH
export GOPATH=$(go env GOPATH)
export GOBIN=$GOPATH/bin
export PATH=$GOROOT/bin:$GOBIN:$PATH
mkdir -p $GOPATH/pkg
mkdir -p $GOBIN

# Place four binaries in your $GOBIN
[ ! -f $GOBIN/protoc-gen-grpc-gateway ] && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
GATEWAY_FILE=$GOBIN/protoc-gen-grpc-gateway
if test -f "$GATEWAY_FILE"; then
    echo "$GATEWAY_FILE already exists, skip making soft link"
else
    ln -s $GOBIN/protoc-gen-grpc-gateway $GATEWAY_FILE
fi
[ ! -f $GOBIN/protoc-gen-openapiv2 ] && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
[ ! -f $GOBIN/protoc-gen-go ] && go install google.golang.org/protobuf/cmd/protoc-gen-go
[ ! -f $GOBIN/protoc-gen-go-grpc ] && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

# resolve go package
go mod init your_service
go mod tidy

# compile proto
buf generate
python -m grpc_tools.protoc -I./proto --python_out=./proto --grpc_python_out=./proto ./proto/sample.proto

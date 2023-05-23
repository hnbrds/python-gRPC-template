#!/bin/bash
# Set your service name (= .proto file name)
SERVICE_NAME=sample

# check if GOROOT is set
if [[ -z "${GOROOT}" ]]; then
  echo "ENV GOROOT is not set, please set environment path GOROOT"
  exit 1
fi

# set go ENVPATH
export GOPATH=$(go env GOPATH)
export GOBIN=$GOPATH/bin
export PATH=$GOROOT/bin:$GOBIN:$PATH
mkdir -p $GOPATH/pkg
mkdir -p $GOBIN
# link project directory to GOPATH/src
ln -s $PWD $GOROOT/src
# link generated file into swagger folder
ln -s ../proto/gen/openapiv2/proto/$SERVICE_NAME.swagger.json ./swagger/swagger.json

# Place four binaries in your $GOBIN
[ ! -f $GOBIN/protoc-gen-grpc-gateway ] && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
GATEWAY_FILE=$GOBIN/protoc-gen-grpc_gateway
if test -f "$GATEWAY_FILE"; then
    echo "$GATEWAY_FILE already exists, skip making soft link"
else
    ln -s $GOBIN/protoc-gen-grpc-gateway $GATEWAY_FILE
fi
[ ! -f $GOBIN/protoc-gen-openapiv2 ] && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
[ ! -f $GOBIN/protoc-gen-go ] && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
[ ! -f $GOBIN/protoc-gen-go-grpc ] && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# resolve go package
rm -rf go.*
go mod init $SERVICE_NAME/v2
go mod tidy

# compile proto
buf generate
python -m grpc_tools.protoc -I./proto --python_out=./proto --grpc_python_out=./proto ./proto/$SERVICE_NAME.proto

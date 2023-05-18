package main

import (
    "context"
    "flag"
    "net/http"

    "github.com/golang/glog"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    // TODO : python-gRPC-template -> {your_project_folder_name}
    gw "python-gRPC-template/proto/gen/go/proto"
)

var (
    // command-line options:
    // gRPC server endpoint
    // TODO : change port to your python gRPC service
    grpcServerEndpoint = flag.String("grpc-server-endpoint",  "localhost:13271", "gRPC server endpoint")
)

func run() error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    // Register gRPC server endpoint
    // Note: Make sure the gRPC server is running properly and accessible
    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
    // TODO : change RegisterSampleServiceHandlerFromEndpoint -> Register{YourService}HandlerFromEndpoint
    err := gw.RegisterSampleServiceHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
    if err != nil {
        return err
    }

    // Start HTTP server (and proxy calls to gRPC server endpoint)
    return http.ListenAndServe(":13270", mux)
}

func main() {
    flag.Parse()
    defer glog.Flush()

    if err := run(); err != nil {
        glog.Fatal(err)
    }
}

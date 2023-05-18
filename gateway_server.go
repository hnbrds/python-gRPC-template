package main

import (
    "context"
    "flag"
    "net/http"
    "mime"
    "io"
    "strings"

    "github.com/philips/grpc-gateway-example/pkg/ui/data/swagger"
    "github.com/philips/go-bindata-assetfs"
    pb "github.com/philips/grpc-gateway-example/echopb"
    
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

func serveSwagger(mux *http.ServeMux) {
    mime.AddExtensionType(".svg", "image/svg+xml")

    // Expose files in third_party/swagger-ui/ on <host>/swagger-ui
    fileServer := http.FileServer(&assetfs.AssetFS{
        Asset:    swagger.Asset,
        AssetDir: swagger.AssetDir,
        Prefix:   "third_party/swagger-ui",
    })
    prefix := "/swagger-ui/"
    mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}


func run() error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()


    mux := http.NewServeMux()
    mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
        io.Copy(w, strings.NewReader(pb.Swagger))
    })
    // Register gRPC server endpoint
    // Note: Make sure the gRPC server is running properly and accessible
    gwmux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
    // TODO : change RegisterSampleServiceHandlerFromEndpoint -> Register{YourService}HandlerFromEndpoint
    err := gw.RegisterSampleServiceHandlerFromEndpoint(ctx, gwmux,  *grpcServerEndpoint, opts)
    if err != nil {
        return err
    }

    mux.Handle("/", gwmux)
    serveSwagger(mux)
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

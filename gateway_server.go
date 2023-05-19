package main

import (
    "log"
    "context"
    "flag"
    "net"
    "net/http"

    "golang.org/x/sync/errgroup"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    "github.com/soheilhy/cmux"
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
    prefix := "/swagger-ui/"
    // TODO : http.Dir({your swagger dist folder})
    fs := http.FileServer(http.Dir("./swagger"))
    mux.Handle(prefix, http.StripPrefix(prefix, fs))
}


func grpcServe(listener net.Listener) error {
    //not implemented
}


func httpServe(listener net.Listener) error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := http.NewServeMux()
    gwmux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
    err := gw.RegisterSampleServiceHandlerFromEndpoint(ctx, gwmux, *grpcServerEndpoint, opts)

    if err != nil {
        log.Fatal(err)
    }

    mux.Handle("/", gwmux)
    serveSwagger(mux)
    s := &http.Server{ Handler : mux }
    return s.Serve(listener)
}


func main() {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    // create listener for server
    listener, err := net.Listen("tcp", ":13270")
    if err != nil {
        log.Fatal(err)
    }
    mux := cmux.New(listener)

    grpcListener := mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
    httpListener := mux.Match(cmux.HTTP1Fast())
    
    g := new(errgroup.Group)
    g.Go(func() error { return grpcServe(grpcListener) })
    g.Go(func() error { return httpServe(httpListener) })
    g.Go(func() error { return mux.Serve() })

    log.Println("RUN SERVER : ", g.Wait())
}

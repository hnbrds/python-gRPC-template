package main

import (
    "log"
    "context"
    "flag"
    "net"
    "net/http"
    "io"

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
    // TODO : change port to match your service port
    gatewayServerAddr = "localhost:13270"
    grpcServerAddr = "localhost:13271"
    // TODO : change swagger dir to your swagger dist directory
    swaggerDir = "./swagger"
    grpcServerEndpoint = flag.String("grpc-server-endpoint",  grpcServerAddr, "gRPC server endpoint")
)

func serveSwagger(mux *http.ServeMux) {
    prefix := "/swagger-ui/"
    fs := http.FileServer(http.Dir(swaggerDir))
    mux.Handle(prefix, http.StripPrefix(prefix, fs))
}

 func forward(src, dest net.Conn) {
    defer src.Close()
    defer dest.Close()
    io.Copy(src, dest)
 }

func handleConnection(conn net.Conn) error {
    log.Println("Connection from gRPC client : ", conn.RemoteAddr())

    remote, err := net.Dial("tcp", grpcServerAddr)
    if err != nil {
        log.Fatal(err)
        return err
    }
    log.Println("Connected to gRPC Server : ", grpcServerAddr)

    // go routines to initiate bi-directional communication for local server with a
    // remote server
    go forward(conn, remote)
    go forward(remote, conn)
    return nil
 }


func grpcServe(listener net.Listener) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
            panic(err)
        }

        go handleConnection(conn)
    }
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
    s.Serve(listener)
    return nil
}


func main() {
    // create listener for server
    listener, err := net.Listen("tcp", gatewayServerAddr)
    if err != nil {
        log.Fatal(err)
    }
    mux := cmux.New(listener)

    //grpcListener := mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
    grpcListener := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
    httpListener := mux.Match(cmux.HTTP1Fast())


    g := errgroup.Group{}
    //g.Go(func() error { return grpcServe(grpcListener) })
    go grpcServe(grpcListener)
    g.Go(func() error { return httpServe(httpListener) })
    g.Go(func() error { return mux.Serve() })
    err = g.Wait()
}

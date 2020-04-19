package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    _ "net/http/pprof"

    gxnet "github.com/AlexStocks/goext/net"
    "github.com/panjf2000/gnet"
)

func initProfiling() {
    var (
        addr string
    )

    addr = gxnet.HostAddress("localhost", 10090)
    go func() {
        http.ListenAndServe(addr, nil)
    }()
}

type echoServer struct {
    *gnet.EventServer
}

func (es *echoServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
    log.Printf("Echo server is listening on %s (multi-cores: %t, loops: %d)\n",
        srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
    return
}
func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
    // Echo synchronously.
    out = frame
    return

    // Echo asynchronously.
    /*
    data := append([]byte{}, frame...)
    go func() {
    	time.Sleep(time.Second)
    	c.AsyncWrite(data)
    }()
    return
    */
}

func main() {
    var port int
    var multicore bool

    initProfiling()
    // Example command: go run echo.go --port 9000 --multicore=true
    flag.IntVar(&port, "port", 9000, "--port 9000")
    flag.BoolVar(&multicore, "multicore", false, "--multicore true")
    flag.Parse()
    echo := new(echoServer)
    log.Fatal(gnet.Serve(echo, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(multicore)))
}
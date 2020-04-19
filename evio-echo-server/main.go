// Copyright 2017 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    _ "net/http/pprof"
    "strings"

    gxnet "github.com/AlexStocks/goext/net"
    "github.com/tidwall/evio"
)

func initProfiling() {
    var (
        addr string
    )

    // addr = *host + ":" + "10000"
    addr = gxnet.HostAddress("localhost", 10088)
    go func() {
        fmt.Println(http.ListenAndServe(addr, nil))
    }()
}

func main() {
    var port int
    var loops int
    var udp bool
    var trace bool
    var reuseport bool
    var stdlib bool

    flag.IntVar(&port, "port", 5000, "server port")
    flag.BoolVar(&udp, "udp", false, "listen on udp")
    flag.BoolVar(&reuseport, "reuseport", false, "reuseport (SO_REUSEPORT)")
    flag.BoolVar(&trace, "trace", false, "print packets to console")
    flag.IntVar(&loops, "loops", 10, "num loops")
    flag.BoolVar(&stdlib, "stdlib", true, "use stdlib")
    flag.Parse()

    initProfiling()
    var events evio.Events
    events.NumLoops = loops
    events.Serving = func(srv evio.Server) (action evio.Action) {
        log.Printf("echo server started on port %d (loops: %d)", port, srv.NumLoops)
        if reuseport {
            log.Printf("reuseport")
        }
        if stdlib {
            log.Printf("stdlib")
        }
        return
    }
    events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
        if trace {
            log.Printf("%s", strings.TrimSpace(string(in)))
        }
        out = in
        return
    }
    scheme := "tcp"
    if udp {
        scheme = "udp"
    }
    if stdlib {
        scheme += "-net"
    }
    fmt.Println("stdlib: ", stdlib)
    log.Fatal(evio.Serve(events, fmt.Sprintf("%s://:%d?reuseport=%t", scheme, port, reuseport)))
}

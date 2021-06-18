package main

import (
	gxsync "github.com/dubbogo/gost/sync"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	gxnet "github.com/AlexStocks/goext/net"
	"github.com/apache/dubbo-getty"
	log "github.com/dubbogo/log4go"
)

func initProfiling() {
	var (
		addr string
	)

	addr = gxnet.HostAddress("localhost", 10089)
	go func() {
		println("listening...")
		http.ListenAndServe(addr, nil)
	}()
}

func initSignal() {
	// signal.Notify的ch信道是阻塞的(signal.Notify不会阻塞发送信号), 需要设置缓冲
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		log.Info("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
		// reload()
		default:
			log.Exit("app exit now...")
			log.Close()
			return
		}
	}
}

type PackageHandler struct{}

func (h *PackageHandler) Read(ss getty.Session, data []byte) (interface{}, int, error) {
	return data, len(data), nil
}

func (h *PackageHandler) Write(ss getty.Session, pkg interface{}) ([]byte, error) {
	return pkg.([]byte), nil
}

type MessageHandler struct {
}

func (h *MessageHandler) OnOpen(session getty.Session) error       { return nil }
func (h *MessageHandler) OnError(session getty.Session, err error) {}
func (h *MessageHandler) OnClose(session getty.Session)            {}
func (h *MessageHandler) OnMessage(session getty.Session, pkg interface{}) {
	time.Sleep(time.Millisecond * 3)
	session.WritePkg(pkg, 50*time.Microsecond)
}
func (h *MessageHandler) OnCron(session getty.Session) {
	println("xxxx")
}

func newSessionCallback(session getty.Session, handler *MessageHandler) error {
	pkgHandler := &PackageHandler{}
	session.SetName("hello-client-session")
	session.SetMaxMsgLen(65535)
	session.SetPkgHandler(pkgHandler)
	session.SetEventListener(handler)
	//session.SetWQLen(32)
	session.SetReadTimeout(3e9)
	session.SetWriteTimeout(3e9)
	session.SetCronPeriod((int)(30e9 / 1e6))
	session.SetWaitTime(3e9)
	return nil
}

func main() {
	_ = getty.SetLoggerLevel(2)
	initProfiling()
	var (
		serverMsgHandler MessageHandler
	)
	addr := "127.0.0.1:5003"
	server := getty.NewTCPServer(
		getty.WithLocalAddress(addr),
		getty.WithServerTaskPool(gxsync.NewTaskPoolSimple(100)),
	)
	newServerSession := func(session getty.Session) error {
		return newSessionCallback(session, &serverMsgHandler)
	}
	server.RunEventLoop(newServerSession)
	initSignal()
}

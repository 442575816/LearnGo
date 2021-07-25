package deploy

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

var (
	router = mux.NewRouter()
)

// HttpServer 处理HTTP请求的Server
type HttpServer interface {
	// 处理请求
	HandleRequest(http.ResponseWriter, *http.Request)

	// 初始化
	Init()
}

// 安装服务
func InstallServer(server HttpServer) {
	server.Init()
}

// AddHandler 增加handler
func AddHandler(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	router.HandleFunc(pattern, handler)
}

func loggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func recoverHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// Startup 启动服务
func Startup(port int) {
	chain := alice.New(loggerHandler, recoverHandler)
	var addr = fmt.Sprintf(":%d", port)
	http.Handle("/", chain.Then(router))
	http.ListenAndServe(addr, nil)
}

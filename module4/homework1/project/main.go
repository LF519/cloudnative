package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

type GracefulServer struct {
	Server           *http.Server
	shutdownFinished chan struct{}
}

func (s *GracefulServer) ListenAndServe() (err error) {
	if s.shutdownFinished == nil {
		s.shutdownFinished = make(chan struct{})
	}

	err = s.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		// expected error after calling Server.Shutdown().
		err = nil
	} else if err != nil {
		err = fmt.Errorf("unexpected error from ListenAndServe: %w", err)
		return
	}

	log.Println("waiting for shutdown finishing...")
	<-s.shutdownFinished
	log.Println("shutdown finished")

	return
}

func (s *GracefulServer) WaitForExitingSignal(timeout time.Duration) {
	var waiter = make(chan os.Signal, 1) // buffered channel
	signal.Notify(waiter, syscall.SIGTERM, syscall.SIGINT)

	// blocks here until there's a signal
	<-waiter

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := s.Server.Shutdown(ctx)
	if err != nil {
		log.Println("shutting down: " + err.Error())
	} else {
		log.Println("shutdown processed successfully")
		close(s.shutdownFinished)
	}
}

func ReturnHeader(w http.ResponseWriter, r *http.Request) {
	// 往response的header上加入requests的header
	for k, v := range r.Header {
		for _, value := range v {
			w.Header().Add(k, value)
		}
	}
	// 往response的header里面添加环境变量里面的version键值对, VERSION不存在时, 值为""
	version := os.Getenv("VERSION")
	w.Header().Add("VERSION", version)
	io.WriteString(w, fmt.Sprintf("%s", w.Header()))
}

// 返回200状态码
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	io.WriteString(w, "200")
}

// 包装函数, 打印请求的地址和返回的状态码
func wrapperFun(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
		vW := reflect.ValueOf(w)
		if vW.Kind() == reflect.Ptr {
			vW = vW.Elem()
		}
		status := vW.FieldByName("status")
		log.Printf("Url: %s, ClientIP: %s StatusCode: %d\n", r.URL, r.Host, status)
	}
}

// 读取yml文件
func InitConfig() {

}

func main() {
	flag.Parse()

	var err error
	defer func() {
		if err != nil {
			log.Println("exited with error: " + err.Error())
		}
	}()

	// 读取配置文件
	work, _ := os.Getwd()                 // 获取当前目录路径
	viper.SetConfigName("config")         // 设置文件名
	viper.SetConfigType("yml")            // 设置文件类型
	viper.AddConfigPath(work + "/config") // 设置配置文件路径
	e := viper.ReadInConfig()
	if e != nil {
		panic(e)
	}
	portStr := viper.GetString("server.port")
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	server := &GracefulServer{
		Server: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
	}

	go server.WaitForExitingSignal(10 * time.Second)

	http.HandleFunc("/", wrapperFun(ReturnHeader))
	http.HandleFunc("/healthz", wrapperFun(healthz))

	log.Printf("listening on port %d...", port)
	err = server.ListenAndServe()
	if err != nil {
		err = fmt.Errorf("unexpected error from ListenAndServe: %w", err)
	}
	log.Println("main goroutine exited.")
}

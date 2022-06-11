package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
)

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

func main() {
	http.HandleFunc("/", wrapperFun(ReturnHeader))
	http.HandleFunc("/healthz", wrapperFun(healthz))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

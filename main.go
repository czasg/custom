package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func Always204(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusNoContent)
}

func Custom(writer http.ResponseWriter, request *http.Request) {
	// method
	if request.Method == http.MethodOptions {
		writer.Header().Set("allow-method", "*")
	}
	// query
	query := request.URL.Query()
	// type
	typeStr := query.Get("type")
	switch typeStr {
	case "txt":
		writer.Header().Set("content-type", "text/plain; charset=utf-8")
	case "json":
		writer.Header().Set("content-type", "application/json; charset=utf-8")
	case "file":
		writer.Header().Set("content-type", "")
	case "zip":
		writer.Header().Set("content-type", "")
	case "jpg":
	case "png":
	case "css":
	case "js":
	}
	// header
	for _, header := range query["header"] {
		kv := strings.SplitN(header, ":", 2)
		if len(kv) != 2 {
			continue
		}
		writer.Header().Set(kv[0], kv[1])
	}
	// code
	codeStr := query.Get("code")
	if codeStr == "" {
		codeStr = "200"
	}
	code := 0
	for i, j := 1, len(codeStr)-1; j >= 0; i, j = i*10, j-1 {
		code += int(codeStr[j]-'0') * i
	}
	if code >= 100 && code < 1000 {
		writer.WriteHeader(code)
	}
	// body
	bodyStr := query.Get("body")
	if bodyStr != "" {
		_, _ = writer.Write([]byte(bodyStr))
	}
}

func main() {
	// handler
	handlers := map[string]func(http.ResponseWriter, *http.Request){
		"/liveness":  Always204,
		"/readiness": Always204,
		"/custom":    Custom,
	}
	for route, handler := range handlers {
		fmt.Printf("handle: %s\n", route)
		http.HandleFunc(route, handler)
	}
	// list & start
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("list on port: %s\n", port)
	_ = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

/* todo
1、是否能够监听 signal 以便优雅退出
2、所有的 content-type
3、head 的预处理
*/

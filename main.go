package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
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
	code, err := strconv.ParseInt(codeStr, 10, 0)
	if err != nil {
		return
	}
	writer.WriteHeader(int(code))
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

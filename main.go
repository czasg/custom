package main

import (
	"fmt"
	"github.com/czasg/custom/types"
	"net/http"
	"os"
	"strings"
)

var (
	contentType                   = "content-type"
	accessAllowControlHeaders     = "access-allow-control-headers"
	accessAllowControlMethods     = "access-allow-control-methods"
	accessAllowControlOrigin      = "access-allow-control-origin"
	accessAllowControlCredentials = "access-allow-control-credentials"
)

func Always204(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusNoContent)
}

func Custom(writer http.ResponseWriter, request *http.Request) {
	// method
	if request.Method == http.MethodOptions {
		writer.Header().Set(accessAllowControlHeaders, "*")
		writer.Header().Set(accessAllowControlMethods, "*")
		writer.Header().Set(accessAllowControlOrigin, "*")
		writer.Header().Set(accessAllowControlCredentials, "true")
	}
	// query
	query := request.URL.Query()
	// type
	typeStr := types.Types[query.Get("type")]
	if typeStr != "" {
		writer.Header().Set(contentType, typeStr)
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
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		_ = fmt.Errorf("server close: %v", err)
	}
}

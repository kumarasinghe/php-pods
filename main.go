package main

import (
	"net/http/httputil"
)

func onRequest(req *httputil.ProxyRequest) string {
	return "http://localhost:8080"
}

func main() {
	StartReverseProxy(onRequest)
}

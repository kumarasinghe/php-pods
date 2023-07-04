package main

import (
	"net/http/httputil"
	rp "php-pods/reverse-proxy"
)

func main() {
	rp.Listen(func(req *httputil.ProxyRequest) string {
		return "http://localhost:8080"
	})
}

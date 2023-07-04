package rp

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ResolveDstHost func(*httputil.ProxyRequest) string

func Listen(resolveDstHost ResolveDstHost) {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(req *httputil.ProxyRequest) {
			// log request
			log.Println(req.In.Method, req.In.URL.String())

			// use callback to get destination host
			dstHost := resolveDstHost(req)

			// set destination host
			dstHostUrl, _ := url.Parse(dstHost)
			req.SetURL(dstHostUrl)

			// set host in outbound request
			// this prevents browser's address bar from changing
			req.Out.Host = req.In.Host
		},
	}

	// Create a new HTTP handler function
	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		// forward the request to the target server
		proxy.ServeHTTP(writer, req)
	})

	// Start the server and listen on port 80
	log.Fatal(http.ListenAndServe(":80", nil))
}

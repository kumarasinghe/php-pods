package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ResolveDstHost func(*httputil.ProxyRequest) (string, error)

const retryTimeout = 5000

func StartReverseProxy(resolveDstHost ResolveDstHost) {
	reverseProxy := &httputil.ReverseProxy{}

	reverseProxy.Rewrite = func(req *httputil.ProxyRequest) {

		// use callback to get destination host
		dstHost, err := resolveDstHost(req)

		// cannot determine destination host
		if err != nil {
			log.Println(err.Error())
			notFoundUrl, _ := url.Parse("/404")
			req.SetURL(notFoundUrl)
			return
		}

		// log external requests
		if !strings.Contains(req.In.URL.String(), dstHost) {
			log.Println(req.In.Method, req.In.URL.String())
		}

		// change request URL to destination host
		dstHostUrl, _ := url.Parse(dstHost)
		req.SetURL(dstHostUrl)

		// set host in outbound request
		// this prevents browser's address bar from changing
		req.Out.Host = req.In.Host
	}

	reverseProxy.ErrorHandler = func(writer http.ResponseWriter, req *http.Request, err error) {
		/* when the original request fail keep trying for some time */

		retryStartTimeStr := req.Header.Get("X-Retry-Start")
		currentTime := time.Now()

		// stamp time on request on the first retry
		if retryStartTimeStr == "" {
			currentTime := strconv.FormatInt(currentTime.Unix(), 10)
			req.Header.Set("X-Retry-Start", currentTime)
			retryStartTimeStr = currentTime
		}

		// calculate the time spent retrying
		retryStartTimeInt, _ := strconv.ParseInt(retryStartTimeStr, 10, 64)
		retryStartTime := time.Unix(retryStartTimeInt, 0)
		retryDuration := currentTime.Sub(retryStartTime)

		// give up retrying if timeout has reached
		if retryDuration.Milliseconds() >= retryTimeout {
			log.Println("proxy request timed out:", req.URL.String())
			writer.WriteHeader(408)
			return
		}

		// resend the request
		reverseProxy.ServeHTTP(writer, req)
	}

	// Create a new HTTP handler function
	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		// forward the request to the target server
		reverseProxy.ServeHTTP(writer, req)
	})

	// Start the server and listen on port 80
	log.Fatal(http.ListenAndServe(":80", nil))
}

package main

import (
	"errors"
	"log"
	"net/http/httputil"
	"regexp"
)

type Orchestrator struct {
	subdomainMatcher *regexp.Regexp
}

func (o Orchestrator) bootstrap() {
	o.subdomainMatcher, _ = regexp.Compile("^.+\\.")
	StartReverseProxy(o.ResolveDstHostHandler)
}

func (o Orchestrator) ResolveDstHostHandler(req *httputil.ProxyRequest) (string, error) {

	// extract subdomain from host as the site name
	siteName := o.subdomainMatcher.FindString(req.In.Host)

	// handle sub domain errors
	if siteName == "" {
		return "", errors.New("could not determine site name")
	}

	siteName = siteName[0 : len(siteName)-1]

	log.Println("**** siteName", siteName)

	return "localhost:8080/", nil
}

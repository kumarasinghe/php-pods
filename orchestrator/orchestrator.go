package main

import (
	"errors"
	"net/http/httputil"
	"regexp"
	"strconv"
)

type Orchestrator struct {
	subdomainMatcher *regexp.Regexp
}

func (o Orchestrator) bootstrap() {
	o.subdomainMatcher, _ = regexp.Compile("^.+\\.")
	StartReverseProxy(o.ResolveDstHostHandler)
}

func (o Orchestrator) ResolveDstHostHandler(req *httputil.ProxyRequest) (string, error) {

	// determine site name by host subdomain
	siteName := o.subdomainMatcher.FindString(req.In.Host)

	if siteName == "" {
		return "", errors.New("could not determine site for " + req.In.Host)
	}

	// trim separators
	siteName = siteName[0 : len(siteName)-1]

	// start container
	port, err := SpawnContainer(siteName)

	if err != nil {
		return "", err
	}

	return "http://localhost:" + strconv.Itoa(port), nil
}

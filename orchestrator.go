package main

import (
	"errors"
	"log"
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

	// determine site name by host
	siteName := o.subdomainMatcher.FindString(req.In.Host)

	if siteName == "" {
		return "", errors.New("could not determine subdomain for " + req.In.Host)
	}

	// trim separators
	siteName = siteName[0 : len(siteName)-1]

	// start container
	port, err := SpawnContainer(siteName)

	if err != nil {
		return "", err
	}

	log.Println("*** Spawn Container", siteName, port)

	return "localhost" + strconv.Itoa(port), nil
}

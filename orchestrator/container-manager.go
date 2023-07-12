package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

var workingDir, _ = os.Getwd()
var projectRoot = filepath.Join(workingDir, "../")
var containerPorts = make(map[string]int)

func SpawnContainer(siteName string) (int, error) {

	// TODO: obtain an unused port for the container
	port := 8080

	log.Println("spawning container", siteName, "on port", port, "...")

	cmd := exec.Command(
		"docker",
		"run",
		"--rm",
		"--name",
		siteName,
		"--network",
		"ppp",
		"--publish",
		strconv.Itoa(port)+":80",
		"--volume",
		filepath.Join(projectRoot, "sites", siteName)+":/var/www/default",
		"-d",
		"ppp-wp",
	)

	err := cmd.Run()

	if err != nil {
		return 0, errors.New("could not start container for " + siteName + " " + err.Error())

	}

	containerPorts[siteName] = port

	return 0, nil

}

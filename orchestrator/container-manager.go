package main

import (
	"errors"
	"log"
	"os/exec"
	"strconv"
)

const PROJECT_ROOT = "C:\\Users\\NaveenKumarasinghe\\workspace\\php-pods\\sites"

func SpawnContainer(siteName string) (int, error) {

	// obtain a random port
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
		// TODO: remove absolute path
		PROJECT_ROOT+"\\"+siteName+":/var/www/default",
		"-d",
		"ppp-wp",
	)

	err := cmd.Run()

	if err != nil {
		return 0, errors.New("could not start container for " + siteName + " " + err.Error())
	}

	return 0, nil

}

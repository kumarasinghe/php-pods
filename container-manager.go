package main

import (
	"errors"
	"os/exec"
	"strconv"
)

func SpawnContainer(siteName string) (int, error) {

	// obtain a random port
	port := 8080

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
		"~/Workspace/php-pods/sites/"+siteName+":/var/www/default",
		"-d",
		"ppp-wordpress",
	)

	err := cmd.Run()

	if err != nil {
		return 0, errors.New("could not start container for " + siteName)
	}

	return 0, nil

}

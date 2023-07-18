package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type ContainerInfo struct {
	port     int
	lastUsed time.Time
}

var containerInfoMap = make(map[string]ContainerInfo)
var workingDir, _ = os.Getwd()
var projectRoot = filepath.Join(workingDir, "../")
var containerTimeout = time.Second * 30

func SpawnContainer(name string) (int, error) {

	containerInfo := containerInfoMap[name]

	// if container exists prolong it
	if containerInfo.port != 0 {
		containerInfo.lastUsed = time.Now()
		containerInfoMap[name] = containerInfo
		return containerInfo.port, nil
	}

	// TODO: obtain an unused port for the container
	port := 8080

	log.Println("spawning container", name, "on port", port, "...")

	cmd := exec.Command(
		"docker",
		"run",
		"--rm",
		"--name",
		name,
		"--network",
		"ppp",
		"--publish",
		strconv.Itoa(port)+":80",
		"--volume",
		filepath.Join(projectRoot, "sites", name)+":/var/www/default",
		"-d",
		"ppp-wp",
	)

	err := cmd.Run()

	if err != nil {
		return 0, errors.New("failed to start container: " + name)
	}

	containerInfoMap[name] = ContainerInfo{port: port, lastUsed: time.Now()}

	return port, nil

}

func containerTerminateRoutine() {

	for {
		// kill containers which met timeout
		for containerName, containerInfo := range containerInfoMap {

			// skip if timeout has not met
			if time.Since(containerInfo.lastUsed) < containerTimeout {
				continue
			}

			// execute kill command
			cmd := exec.Command(
				"docker",
				"kill",
				containerName,
			)

			err := cmd.Run()

			delete(containerInfoMap, containerName)

			if err != nil {
				log.Println("failed to terminate idle container", containerName, err.Error())
				continue
			}

			log.Println("terminated idle container:", containerName)

		}

		time.Sleep(time.Second)
	}
}

func init() {
	go containerTerminateRoutine()
}

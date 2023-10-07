package main

import (
	"context"
	"fmt"
	"github.com/avakhov/docker_clean_containers/util"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"os"
	"strings"
)

func doMain(prefix string) error {
	// initialize docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return util.WrapError(err)
	}

	// stop and remove containers with the specified prefix
	containers, err := dockerClient.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return util.WrapError(err)
	}
	var containersToStop []string
	var containersToRemove []string
	for _, container := range containers {
		if strings.HasPrefix(container.Names[0], "/"+prefix) {
			containersToStop = append(containersToStop, container.ID)
			containersToRemove = append(containersToRemove, container.ID)
		}
	}

	// stop containers
	fmt.Println("Stopping containers...")
	for _, containerID := range containersToStop {
		timeout := 3
		err := dockerClient.ContainerStop(context.Background(), containerID, container.StopOptions{
			Timeout: &timeout,
		})
		if err != nil {
			return util.WrapError(err)
		}
		fmt.Printf("Stopped container %s\n", containerID)
	}

	// remove containers
	fmt.Println("Removing containers...")
	for _, containerID := range containersToRemove {
		err := dockerClient.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{})
		if err != nil {
			return util.WrapError(err)
		}
		fmt.Printf("Removed container %s\n", containerID)
	}

	// out
	return nil
}

func main() {
	fmt.Printf("docker_clean_containers version=%s\n", util.GetVersion())
	if len(os.Args) < 2 {
		fmt.Println("Usage ./docker_clean_containers <prefix>")
		os.Exit(1)
	}
	err := doMain(os.Args[1])
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("done\n")
}

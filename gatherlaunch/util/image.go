package util

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func PullDockerImage(imageName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Println("Failed to create new docker client")
		return
	}
	defer cli.Close()
	reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		log.Printf("Failed to pull docker image %s\n", imageName)
		return
	}
	defer reader.Close()
}

package util

import (
	"context"
	"io"
	"log"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func GetAllImages() []string {
	return []string{
		"imarenf/osv.dev:1.3",
		"imarenf/toxic-repos:1.1",
		"imarenf/govulncheck:1.0",
		"imarenf/packj-static:1.0",
		"imarenf/packj-trace:1.1",
		"imarenf/decision-making:1.2",
		"imarenf/deps.dev:1.0",
		"imarenf/ossgadget:1.3",
		"imarenf/appinspector:1.2",
	}
}

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

	out, err := io.ReadAll(reader)
	if err != nil {
		log.Println("Error while parsing: ", err)
	}
	log.Println("Out from ImagePull: ", string(out))
}

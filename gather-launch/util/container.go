package util

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type PackageInfo struct {
	Registry string
	Name     string
	Type     string
	Version  string
	Purl     string
}

type ContainerOutput struct {
	ToolName string
	Output   string `json:"result"`
	ExitCode int64  `json:"exit_code"`
}

type ToolResponse struct {
	RespCh chan ContainerOutput
	ErrCh  chan error
}

func getContainerCmd(toolName string, pkgInfo PackageInfo) ([]string, error) {
	switch toolName {
	case "packj-static":
		pkg := pkgInfo.Registry + ":" + pkgInfo.Name + ":" + pkgInfo.Version
		return []string{pkg}, nil
	case "packj-trace":
		pkg := pkgInfo.Registry + ":" + pkgInfo.Name + ":" + pkgInfo.Version
		return []string{pkg}, nil
	case "deps.dev":
		return []string{pkgInfo.Purl}, nil
	case "osv.dev":
		return []string{pkgInfo.Purl}, nil
	case "toxic-repos":
		return []string{pkgInfo.Purl}, nil
	default:
		return nil, fmt.Errorf("Unexpected tool name: %s", toolName)
	}
}

func RunDockerContainer(toolName, toolImage string, pkgInfo PackageInfo, tr ToolResponse) {
	log.Printf("Started processing with the tool %s\n", toolName)
	ctx := context.Background()
	containerName := "faross-" + toolName
	containerCmd, err := getContainerCmd(toolName, pkgInfo)
	if err != nil {
		log.Println("Failed to get container cmd")
		tr.ErrCh <- err
		return
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Println("Failed to create new docker client")
		tr.ErrCh <- err
		return
	}
	defer cli.Close()

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: toolImage,
			Cmd:   containerCmd,
			Tty:   true,
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		&v1.Platform{},
		containerName,
	)
	if err != nil {
		log.Printf("Failed to create the container for the tool %s\n", toolName)
		tr.ErrCh <- err
		return
	}
	defer cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Printf("Failed to start the container for the tool %s\n", toolName)
		tr.ErrCh <- err
		return
	}

	statusCh, waitErrCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	var exitCode int64
	select {
	case err := <-waitErrCh:
		if err != nil {
			log.Printf("Container for the tool %s returned the error %s\n", toolName, err.Error())
			tr.ErrCh <- err
			return
		}
	case status := <-statusCh:
		log.Printf("Container for the tool %s finished successfully\n", toolName)
		exitCode = status.StatusCode
	}

	outRaw, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		log.Printf("Failed to get container logs for the tool %s\n", toolName)
		tr.ErrCh <- err
		return
	}

	out, err := io.ReadAll(outRaw)
	if err != nil {
		log.Printf("Failed to get container logs for the tool %s\n", toolName)
		tr.ErrCh <- err
		return
	}

	containerOutput := ContainerOutput{
		ToolName: toolName,
		Output:   string(out),
		ExitCode: exitCode,
	}

	tr.RespCh <- containerOutput
}

package gatherlaunch

import (
	"fmt"
	"log"
	"os"
	"sync"
	"text/template"

	"faross/gatherlaunch/util"

	"github.com/package-url/packageurl-go"
)

func InitGatherLaunch(instrumentsPath string) error {
	err := util.InitTools(instrumentsPath)
	if err != nil {
		return err
	}
	// Get list of all images
	// err = PullImages(...)
	return nil
}

func PullImages(images map[string]string) {
	var wg sync.WaitGroup
	for toolName, toolImage := range images {
		wg.Add(1)
		go func(toolName string, toolImage string) {
			defer wg.Done()
			util.PullDockerImage(toolImage)
		}(toolName, toolImage)
	}
	wg.Wait()
}

func Scan(purl packageurl.PackageURL) (map[string]any, error) {
	pkgInfo := util.ParsePurl(purl)

	toolsImageMapping, err := util.SelectTools(pkgInfo.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to select tools: %s", err)
	}

	ResultMapping := make(map[string]util.ToolResponse, len(toolsImageMapping))
	var wg sync.WaitGroup
	for toolName, toolImage := range toolsImageMapping {
		wg.Add(1)
		ResultMapping[toolName] = util.ToolResponse{
			RespCh: make(chan util.ContainerOutput, 1),
			ErrCh:  make(chan error, 1),
		}

		go func(toolName string, toolImage string) {
			defer wg.Done()
			util.RunDockerContainer(toolName, toolImage, pkgInfo, ResultMapping[toolName])
		}(toolName, toolImage)
	}
	wg.Wait()

	containerOutputs := []util.ContainerOutput{}
	for toolName, res := range ResultMapping {
		select {
		case err = <-res.ErrCh:
			log.Printf("Tool %s exited unsuccessfully:\n%s\n", toolName, err.Error())
		case respRaw := <-res.RespCh:
			parser, err := util.GetParser(toolName)
			if err != nil {
				log.Printf("Failed to get parser for the tool %s\n", toolName)
			}
			resp, err := parser.Parse(respRaw)
			if err != nil {
				log.Printf("Failed to parse container output for the tool %s\n", toolName)
			}
			containerOutputs = append(containerOutputs, respRaw)

			// TODO: replace raw output with writing into a file/storing into smwh and implement parsers for every tool
			log.Printf("Output for the tool %s:\n%s\n", toolName, string(resp))
		}
	}
	log.Println("All checks have finished successfully")

	tmpl, err := template.New("result").Parse(util.OutputTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse output template: %s", err)
	}
	// TODO: replace os.Stdout with file if needed
	tmpl.Execute(os.Stdout, containerOutputs)

	// TODO: replace placeholder with dedcision-making report
	return map[string]any{
		"final_score": 6.0,
	}, nil
}

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
	// Init function, must be called once before calling Scan
	err := util.InitTools(instrumentsPath)
	if err != nil {
		return err
	}
	images := util.GetAllImages()
	var wg sync.WaitGroup
	for _, toolImage := range images {
		wg.Add(1)
		go func(toolImage string) {
			defer wg.Done()
			util.PullDockerImage(toolImage)
		}(toolImage)
	}
	wg.Wait()
	return nil
}

func Scan(purl packageurl.PackageURL) (*util.Decision, error) {
	// Scan the package with all the tools in instruments
	// Returns map[string]any, a score for the proper package
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
			util.RunCheck(toolName, toolImage, pkgInfo, ResultMapping[toolName])
		}(toolName, toolImage)
	}
	wg.Wait()

	containerOutputs := []util.ContainerOutput{}
	var traceResp util.ContainerOutput
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
			for _, r := range resp {
				if r.ToolName == "packj-trace" {
					traceResp = r
				} else {
					containerOutputs = append(containerOutputs, r)
				}
			}

			log.Printf("Output for the tool %s:\n%s\n", toolName, util.RespToString(resp))
		}
	}
	containerOutputs = append(containerOutputs, traceResp)
	log.Println("All checks have finished successfully")

	tmpl, err := template.New("result").Parse(util.OutputTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse output template: %s", err)
	}

	dname, err := os.MkdirTemp("", "faross")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(dname)
	f, err := os.Create(dname + "/input.json")

	if err != nil {
		return nil, fmt.Errorf("failed to create file for containers output %s", err)
	}
	defer f.Close()

	tmpl.Execute(f, containerOutputs)

	decision, err := util.RunDecisionMaking(dname)
	if err != nil {
		return nil, fmt.Errorf("decision-making finished with an error %s", err)
	}
	return &decision, nil
}

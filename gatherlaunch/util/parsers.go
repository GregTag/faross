package util

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/package-url/packageurl-go"
)

func ParsePurl(purl packageurl.PackageURL) PackageInfo {
	return PackageInfo{
		Registry:  purl.Type,
		Name:      purl.Name,
		Type:      purl.Type,
		Version:   purl.Version,
		Namespace: purl.Namespace,
		Purl:      purl.ToString(),
	}
}

type Parser interface {
	Parse(ContainerOutput) ([]ContainerOutput, error)
}

type DefaultParser struct {
}

func (dp DefaultParser) Parse(respRaw ContainerOutput) ([]ContainerOutput, error) {
	return []ContainerOutput{respRaw}, nil
}

type OssgadgetOutputParser struct {
}

func (op OssgadgetOutputParser) Parse(respRaw ContainerOutput) ([]ContainerOutput, error) {
	output := "{" + strings.Split(respRaw.Output, "{")[1]
	output = strings.Split(output, "}")[0] + "}"
	r := ContainerOutput{
		ToolName: respRaw.ToolName,
		ExitCode: respRaw.ExitCode,
		Output:   output,
	}
	return []ContainerOutput{r}, nil
}

type AppInspectorParser struct {
}

func (ap AppInspectorParser) Parse(respRaw ContainerOutput) ([]ContainerOutput, error) {
	var unsafeOperations, fileTypes string
	if respRaw.ExitCode != 0 {
		unsafeOperations = "{\"checkName\": \"unsafe-operations\", \"score\": \"?\", \"risk\": \"Medium\", \"description\": \"Error while fetching data\"}"
		fileTypes = "{\"checkName\": \"unsafe-operations\", \"score\": \"?\", \"risk\": \"Medium\", \"description\": \"Error while fetching data\"}"
	} else {
		unsafeOperations, fileTypes = strings.Split(respRaw.Output, "\r\n")[0], strings.Split(respRaw.Output, "\r\n")[1]
	}
	r1 := ContainerOutput{
		ToolName: respRaw.ToolName + "-operations",
		ExitCode: respRaw.ExitCode,
		Output:   unsafeOperations,
	}
	r2 := ContainerOutput{
		ToolName: respRaw.ToolName + "-filetypes",
		ExitCode: respRaw.ExitCode,
		Output:   fileTypes,
	}
	res := []ContainerOutput{r1, r2}
	return res, nil
}

type ScorecardParser struct {
}

func (sc ScorecardParser) Parse(respRaw ContainerOutput) ([]ContainerOutput, error) {
	var output string
	if respRaw.ExitCode != 0 {
		output = "{\"checkName\": \"scorecard\", \"score\": \"?\", \"risk\": \"Low\", \"description\": \"Error while fetching github repo\"}"
	} else {
		output = respRaw.Output
	}
	r := ContainerOutput{
		ExitCode: respRaw.ExitCode,
		ToolName: respRaw.ToolName,
		Output:   output,
	}
	return []ContainerOutput{r}, nil
}

type SafeParser struct {
}

func (sp SafeParser) Parse(respRaw ContainerOutput) ([]ContainerOutput, error) {
	var output string
	if respRaw.ExitCode != 0 {
		output = "{\"checkName\": \"packj\", \"score\": \"?\", \"risk\": \"Low\", \"description\": \"Error while fetching data\"}"
	} else {
		output = respRaw.Output
	}
	r := ContainerOutput{
		ExitCode: respRaw.ExitCode,
		ToolName: respRaw.ToolName,
		Output:   output,
	}
	return []ContainerOutput{r}, nil
}

func RespToString(out []ContainerOutput) []byte {
	resp, _ := json.Marshal(out)
	return resp
}

func GetParser(toolName string) (Parser, error) {
	switch toolName {
	case "packj-static":
		return SafeParser{}, nil
	case "packj-trace":
		return SafeParser{}, nil
	case "deps.dev":
		return DefaultParser{}, nil
	case "osv.dev":
		return DefaultParser{}, nil
	case "toxic-repos":
		return DefaultParser{}, nil
	case "ossgadget":
		return OssgadgetOutputParser{}, nil
	case "application-inspector":
		return AppInspectorParser{}, nil
	case "scorecard":
		return ScorecardParser{}, nil
	default:
		return nil, fmt.Errorf("unexpected tool name: %s", toolName)
	}
}

type Decision struct {
	Score           float64 `json:"score"`
	IsQuarantined   bool    `json:"is_quarantined"`
	ImpactfulScores []any   `json:"impactful_scores"`
}

var FailDecision = Decision{
	Score:           6,
	IsQuarantined:   true,
	ImpactfulScores: []any{},
}

func ParseDecision(out []byte) (Decision, error) {
	desicion := Decision{}
	err := json.Unmarshal(out, &desicion)
	if err != nil {
		return FailDecision, err
	}
	return desicion, nil
}

func GetPurl(purlRaw string) (packageurl.PackageURL, error) {
	return packageurl.FromString(purlRaw)
}

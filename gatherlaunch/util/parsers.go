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
	unsafeOperations, fileTypes := strings.Split(respRaw.Output, "\r\n")[0], strings.Split(respRaw.Output, "\r\n")[1]

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

func RespToString(out []ContainerOutput) []byte {
	resp, _ := json.Marshal(out)
	return resp
}

func GetParser(toolName string) (Parser, error) {
	switch toolName {
	case "packj-static":
		return DefaultParser{}, nil
	case "packj-trace":
		return DefaultParser{}, nil
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
		return DefaultParser{}, nil
	default:
		return nil, fmt.Errorf("unexpected tool name: %s", toolName)
	}
}

type Decision struct {
	Score         float64 `json:"score"`
	IsQuarantined bool    `json:"is_quarantined"`
}

var FailDecision = Decision{
	Score:         6,
	IsQuarantined: true,
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

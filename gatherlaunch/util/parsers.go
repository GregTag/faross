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
	Parse(ContainerOutput) (ContainerOutput, error)
}

type DefaultParser struct {
}

func (dp DefaultParser) Parse(respRaw ContainerOutput) (ContainerOutput, error) {
	return respRaw, nil
}

type PowershellOutputParser struct {
}

func (dp PowershellOutputParser) Parse(respRaw ContainerOutput) (ContainerOutput, error) {
	output := "{" + strings.Split(respRaw.Output, "{")[1]
	output = strings.Split(output, "}")[0] + "}"
	r := ContainerOutput{
		ToolName: respRaw.ToolName,
		ExitCode: respRaw.ExitCode,
		Output:   output,
	}
	return r, nil
}

func RespToString(out ContainerOutput) []byte {
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
		return PowershellOutputParser{}, nil
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

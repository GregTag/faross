package util

import (
	"encoding/json"
	"fmt"
	"log"

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
	Parse(ContainerOutput) ([]byte, error)
}

type DefaultParser struct {
}

func (dp DefaultParser) Parse(respRaw ContainerOutput) ([]byte, error) {
	resp, err := json.Marshal(respRaw)
	if err != nil {
		log.Printf("Failed to parse container output: %s", err.Error())
		return nil, err
	}
	return resp, nil
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
	default:
		return nil, fmt.Errorf("unexpected tool name: %s", toolName)
	}
}

type Decision struct {
	Score any `json:"score"`
}

func ParseDecision(out []byte) (any, error) {
	desicion := Decision{}
	err := json.Unmarshal(out, &desicion)
	if err != nil {
		return nil, err
	}
	return desicion.Score, nil
}

func GetPurl(purlRaw string) (packageurl.PackageURL, error) {
	return packageurl.FromString(purlRaw)
}

package util

import (
	"fmt"
	"maps"

	koanfjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func SelectTools(packageType string) (map[string]string, error) {
	k := koanf.New(".")
	if err := k.Load(file.Provider("../instruments.json"), koanfjson.Parser()); err != nil {
		return nil, fmt.Errorf("Error while loading config file: %v", err)
	}
	toolMap := k.StringMap(packageType + ".static")
	dynamicToolMap := k.StringMap(packageType + ".dynamic")
	// TODO: add spliting dynamic and static split if needed
	maps.Copy(toolMap, dynamicToolMap)
	return toolMap, nil
}

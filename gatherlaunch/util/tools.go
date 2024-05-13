package util

import (
	"fmt"
	"maps"

	koanfjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

func InitTools(instrumentsPath string) error {
	if err := k.Load(file.Provider(instrumentsPath), koanfjson.Parser()); err != nil {
		return fmt.Errorf("error while loading config file: %v", err)
	}
	return nil
}

func SelectTools(packageType string) (map[string]string, error) {
	toolMap := k.StringMap(packageType + ".static")
	dynamicToolMap := k.StringMap(packageType + ".dynamic")
	// TODO: add spliting dynamic and static split if needed
	maps.Copy(toolMap, dynamicToolMap)
	return toolMap, nil
}

package config

import (
	"log"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var Koanf = koanf.New(".")

func Load(path string) {
	err := Koanf.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}

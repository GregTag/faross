package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var Koanf = koanf.New(".")

func Load(path string) {
	err := Koanf.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		log.Fatalf("Error loading yaml config: %v", err)
	}
	err = Koanf.Load(env.Provider("FIREWALL_", ".", func(s string) string {
		return strings.TrimPrefix(s, "FIREWALL_")
	}), nil)
	if err != nil {
		log.Fatalf("Error reading envs: %v", err)
	}
}

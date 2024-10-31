package config

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
	"os"
)

func ParseConfig(file string) *Config {
	if _, err := os.Stat(file); err != nil {
		log.Fatal().Msg(err.Error())
	}
	tomlBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	var appConfig Config
	if err := toml.Unmarshal(tomlBytes, &appConfig); err != nil {
		log.Fatal().Msg(err.Error())
	}
	return &appConfig
}

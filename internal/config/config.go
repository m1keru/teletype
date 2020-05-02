package config

import (
	log "github.com/sirupsen/logrus"
	"os"

	"gopkg.in/yaml.v2"
)

//Telegram - holder for Telegram Configuration
type Telegram struct {
	Token    string   `yaml:"Token"`
	Interval int      `yaml:"Interval"`
	Chats    []string `yaml:"Chats"`
}

//Speech - SpeechConfig
type Speech struct {
	SampleRateHertz uint32 `yaml:"SampleRateHertz"`
	LanguageCode    string `yaml:"LanguageCode"`
	SecretFile      string `yaml:"SecretFile"`
}

//Daemon - DaemonConfig
type Daemon struct {
	LogFile string `yaml:"LogFile"`
	Debug   bool   `yaml:"Debug"`
}

//Config - Config
type Config struct {
	Daemon   Daemon   `yaml:"Daemon"`
	Speech   Speech   `yaml:"Speech"`
	Telegram Telegram `yaml:"Telegram"`
}

//Setup - Setup
func (cfg *Config) Setup(filename *string) error {
	configFile, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("Unable to read config file. Error:\n %v\n", err)
	}
	defer configFile.Close()
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Unable to Unmarshal Config, Error:\n %v\n", err)
	}
	return nil
}

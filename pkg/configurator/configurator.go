package configurator

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Constants struct {
	Right     string
	Left      string
	VeryRight string
	VeryLeft  string
	Today     string
	Tomorrow  string
	ThisWeek  string
	NextWeek  string
}

type Config struct {
	Url        string
	Port       string
	DbLogin    string
	DbPassword string
	TgKey      string
	PageSize   int
	Consts     Constants
}

var Cfg *Config

func init() {
	// Load config
	var err error
	Cfg, err = NewConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
}

func NewConfig(cfgPath string) (*Config, error) {
	jsonFile, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(jsonBytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

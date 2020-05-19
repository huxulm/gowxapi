package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config is the expression that yaml is configured in struct
type Config struct {
	Mongo   `yaml:"mongo"`
	AppInfo `yaml:"app_info"`
	Seed    `yaml:"seed"`
}

// Mongo is the expression that yaml is configured in struct
type Mongo struct {
	MongodbURL string `yaml:"mongodb_url"`
	DbName     string `yaml:"db_name"`
}

// AppInfo is the expression that yaml is configured in struct
type AppInfo struct {
	AppID  string `yaml:"appid"`
	Secret string `yaml:"secret"`
}

// Seed is the expression that yaml is configured in struct
type Seed struct {
	SeedFile string `yaml:"seed_file"`
}

// C is a global config initilized when application starts
var C *Config

func init() {
	C = &Config{}
	configFilePath := os.Getenv("CONFIG")
	if len(configFilePath) == 0 {
		configFilePath = "config.yaml"
	}
	confB, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalln("Please make you have a config.yaml")
	}
	err = yaml.Unmarshal(confB, C)
	if err != nil {
		panic(err)
	}
	log.Println("Config loaded successfully.")
}

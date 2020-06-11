package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config is the expression that yaml is configured in struct
type Config struct {
	HttpListen     string `yaml:"http_listen"`
	WsOrigin       string `yaml:"ws_origin"`
	Mongo          `yaml:"mongo"`
	AppInfo        `yaml:"app_info"`
	Seed           `yaml:"seed"`
	SeedTour       `yaml:"seed_tour"`
	StaticResource `yaml:"static_resource"`
}

// StaticResource defines static resource for this app
type StaticResource struct {
	LogoPath   string `yaml:"logo_path"`
	AvatarPath string `yaml:"avatar_path"`
	FontsPath  string `yaml:"fonts_path"`
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
	Seed     bool   `yaml:"seed"`
	SeedFile string `yaml:"seed_file"`
}

// SeedTour is the expression that yaml is configured in struct
type SeedTour struct {
	Seed     bool   `yaml:"seed"`
	TourPath string `yaml:"tour_path"`
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

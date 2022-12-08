package conf

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var _once = sync.Once{}
var homeDir, _ = os.UserHomeDir()
var config = Config{
	Port: 8080,
	Database: database{
		File: filepath.Join(homeDir, "gask", "data.sqlite3"),
	},
}

type database struct {
	File  string `json:"file" yaml:"file"`
	Debug bool   `json:"debug" yaml:"debug"`
}

type Config struct {
	Debug      bool     `json:"debug" yaml:"debug"`
	SwaggerDoc bool     `json:"swaggerDoc" yaml:"swaggerDoc"`
	Port       int      `json:"port" yaml:"port"`
	Host       string   `json:"host" yaml:"host"`
	Database   database `json:"database" yaml:"database"`
}

func load() {
	viper.SetDefault("server", config)
	if viper.InConfig("server") {
		err := viper.UnmarshalKey("server", &config)
		if err != nil {
			panic(err)
		}
	}
	if viper.GetBool("debug") {
		config.Debug = true
	}
}

func GetConfig() Config {
	_once.Do(load)
	return config
}

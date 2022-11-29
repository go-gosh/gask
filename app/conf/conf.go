package conf

import (
	"os"
	"path"
	"sync"

	"github.com/spf13/viper"
)

var _once = sync.Once{}
var homeDir, _ = os.UserHomeDir()
var config = Config{
	Port: 8080,
	Database: database{
		File:  path.Join(homeDir, "gask", "data.sqlite3"),
		Debug: true,
	},
}

type database struct {
	File  string `json:"file" yaml:"file"`
	Debug bool   `json:"debug" yaml:"debug"`
}

type Config struct {
	Port     int      `json:"port" yaml:"port"`
	Database database `json:"database" yaml:"database"`
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
		config.Database.Debug = true
	}
}

func GetConfig() Config {
	_once.Do(load)
	return config
}

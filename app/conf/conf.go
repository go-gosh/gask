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
	Database: Database{
		File:  path.Join(homeDir, "gask", "data.sqlite3"),
		Debug: false,
	},
}

type Database struct {
	File  string `json:"file" yaml:"file"`
	Debug bool   `json:"debug" yaml:"debug"`
}

type Config struct {
	Port     int      `json:"port" yaml:"port"`
	Database Database `json:"database" yaml:"database"`
}

func load() {
	viper.SetDefault("server", config)
	if viper.InConfig("server") {
		err := viper.UnmarshalKey("server", &config)
		if err != nil {
			panic(err)
		}
	}
}

func GetConfig() Config {
	_once.Do(load)
	return config
}

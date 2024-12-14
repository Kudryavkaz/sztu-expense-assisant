package config

import (
	"strings"

	"github.com/spf13/viper"
)

var Config *viper.Viper

func init() {
	Config = viper.New()

	Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	Config.SetDefault("mysql.host", "127.0.0.1")
	Config.SetDefault("mysql.port", "3306")
	Config.SetDefault("mysql.user", "")
	Config.SetDefault("mysql.database", "tree_sitter_parser")
	Config.SetDefault("mysql.password", "")

	Config.SetDefault("log.level", "")

	Config.AutomaticEnv() // read in environment variables that match
}

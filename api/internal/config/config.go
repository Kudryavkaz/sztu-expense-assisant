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
	Config.SetDefault("mysql.database", "sztuea")
	Config.SetDefault("mysql.password", "")

	Config.SetDefault("log.level", "")

	Config.SetDefault("middleware.jwt.secret", "")
	Config.SetDefault("middleware.token.expired", 24)

	Config.SetDefault("redis.host", "127.0.0.1")
	Config.SetDefault("redis.port", "6379")
	Config.SetDefault("redis.password", "")

	Config.SetDefault("grpc.client.address", "127.0.0.1")
	Config.SetDefault("grpc.client.port", "50051")
	Config.SetDefault("crawler.jwt.secret.key", "")

	Config.AutomaticEnv()
}

package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	ServerHost string `json:"host"`
	ServerPort string `json:"port"`
	RedisHost  string `json:"redhost"`
	RedisPort  string `json:"redport"`
	RedisPwd   string
	RedisDb    int
	PgHost     string
	PgPort     string
	PgPwd      string
	PgUser     string
	PgDb       string
	Psql       string
}

var Config Configuration

func GetConfig() *Configuration {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "81")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Config file not found \n", err)
	}
	Config = Configuration{
		ServerPort: viper.GetString("server.port"),
		ServerHost: viper.GetString("server.host"),
		RedisHost:  viper.GetString("redis.host"),
		RedisPort:  viper.GetString("redis.port"),
		RedisPwd:   viper.GetString("redis.password"),
		RedisDb:    viper.GetInt("redis.db"),
		PgHost:     viper.GetString("postgres.host"),
		PgPort:     viper.GetString("postgres.port"),
		PgUser:     viper.GetString("postgres.user"),
		PgPwd:      viper.GetString("postgres.password"),
		PgDb:       viper.GetString("postgres.db"),
	}

	Config.Psql = fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		Config.PgHost,
		Config.PgPort,
		Config.PgUser,
		Config.PgPwd,
		Config.PgDb)

	return &Config
}

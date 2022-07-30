package app

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	NumOfWorkers int
	JobsBuffer   int
	LogFile      string
}

type UserInput struct {
	ChannelId  string `mapstructure:"channel_id"`
	NumOfChats int    `mapstructure:"num_of_chats"`
	APIToken   string `mapstructure:"api_token"`
}

var (
	Config ServerConfig
	Input  *UserInput
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("server.buffered_jobs", 500)
	viper.SetDefault("server.workers", 10)
	viper.SetDefault("server.log_file", "main.log")

	// Load config to viper
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error while reading config file: ", err)
	}

	// Store server config
	Config = ServerConfig{
		NumOfWorkers: viper.GetInt("server.workers"),
		JobsBuffer:   viper.GetInt("server.buffered_jobs"),
		LogFile:      viper.GetString("server.log_file"),
	}

	// Store user input
	err := viper.UnmarshalKey("input", &Input)
	if err != nil {
		log.Fatal("Unable to convert user input into struct:", err)
	}
	if len(Input.APIToken) == 0 {
		log.Fatal("Empty API token")
	}
}

package app

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	NumOfWorkers int    `mapstructure:"workers"`
	JobsBuffer   int    `mapstructure:"buffered_jobs"`
	LogFile      string `mapstructure:"log_file"`
}

type UserInput struct {
	UseJSON    bool   `mapstructure:"use_json"`
	ChannelId  string `mapstructure:"channel_id"`
	NumOfChats int    `mapstructure:"num_of_chats"`
	APIToken   string `mapstructure:"api_token"`
}

var (
	Config *ServerConfig
	Input  *UserInput
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// Load config to viper
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error while reading config file: ", err)
	}

	// Store server config
	if err := viper.UnmarshalKey("server", &Config); err != nil {
		log.Fatal("Unable to convert server config into struct:", err)
	}

	// Store user input
	if err := viper.UnmarshalKey("input", &Input); err != nil {
		log.Fatal("Unable to convert user input into struct:", err)
	}

	if !Input.UseJSON && len(Input.APIToken) == 0 {
		log.Fatal("Empty API token")
	}
}

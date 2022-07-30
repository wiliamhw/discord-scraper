package util

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
	ChannelId  string
	NumOfChats string
}

var (
	Config ServerConfig
	Input  UserInput
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("server.bufferedJobs", 500)
	viper.SetDefault("server.workers", 10)
	viper.SetDefault("server.logFile", "./main.log")

	// Load config to viper
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error while reading config file", err)
	}

	// Store server config
	Config = ServerConfig{
		NumOfWorkers: viper.GetInt("server.workers"),
		JobsBuffer:   viper.GetInt("server.bufferedJobs"),
		LogFile:      viper.GetString("server.logFile"),
	}

	// Store user input
	err := viper.UnmarshalKey("input", Input)
	if err != nil {
		log.Fatal("Unable to convvert user input into struct:", err)
	}
}

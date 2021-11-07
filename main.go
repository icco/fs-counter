package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Paths []*PathConfig
}

type Counter string

const (
	Directory Counter = "dir"
	File      Counter = "file"
)

type PathConfig struct {
	FilePath string
	Count    Counter
	Include  string
	Exclude  string
}

func main() {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.fs-counter")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w", err))
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("Fatal error unmarshaling: %w", err))
	}
}

package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

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
}

func main() {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.fs-counter")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %+v", err)
	}

	var c *Config
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("Fatal error unmarshaling: %+v", err)
	}

	log.Printf("found config: %+v", c)

	counts := map[string]int64{}

	for _, pc := range c.Paths {
		log.Printf("searching path: %+v", pc)
		fileSystem := os.DirFS(pc.FilePath)
		if err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("could not walk %q: %w", path, err)
			}

			// Ignore hidden files
			if strings.HasPrefix(path, ".") {
				return fs.SkipDir
			}

			switch pc.Count {
			case Directory:
				if d.IsDir() {
					counts[pc.FilePath] += 1
				}
			case File:
				if !d.IsDir() {
					counts[pc.FilePath] += 1
				}
			default:
				return fmt.Errorf("not a valid Counter")
			}

			return nil
		}); err != nil {
			log.Fatalf("walk error: %+v", err)
		}
	}

	log.Printf("Found: %+v", counts)
}

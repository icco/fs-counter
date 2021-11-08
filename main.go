package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
	Include  string
}

func (pc *PathConfig) IsValid(path string, d fs.DirEntry) (bool, error) {
	// Ignore hidden files
	if path != "." && d.IsDir() {
		if strings.HasPrefix(path, ".") {
			log.Printf("skipping %q", path)
			return false, fs.SkipDir
		}
	}

	switch pc.Count {
	case Directory:
		if d.IsDir() {
			if pc.Include == "" {
				return true, nil
			}
			return filepath.Match(pc.Include, path)
		}
	case File:
		if !d.IsDir() {
			if pc.Include == "" {
				return true, nil
			}
			return filepath.Match(pc.Include, path)
		}
	default:
		return false, fmt.Errorf("not a valid Counter")
	}

	return false, nil
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

		counts[pc.FilePath] = 0
		fileSystem := os.DirFS(pc.FilePath)
		if err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("could not walk %q: %w", path, err)
			}

			inc, err := pc.IsValid(path, d)
			if err != nil {
				return err
			}

			if inc {
				counts[pc.FilePath] += 1
			}

			return nil
		}); err != nil {
			log.Fatalf("walk error: %+v", err)
		}
	}

	log.Printf("Found: %+v", counts)
}

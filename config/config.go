package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func init() {
	file, err := os.Open("/etc/pram/pram.conf")
	if err != nil {
		// file was not found so use default settings
		Settings = &Config{
			Index: Index{
				Host: "127.0.0.1",
				Port: 5005,
			},
			Directories: Directories{
				AssetsDir: "/data/prim/assets/",
			},
		}
		return
	}

	// if the file is found fill settings with json
	Settings = &Config{}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&Settings)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// Settings holds the current config options
var Settings *Config

// Config represents the possible configurable parameters
// for the local daemon
type Config struct {
	Index       Index
	Directories Directories
	Database    Database
}

// Index sets what the daemon listens on
type Index struct {
	Host                   string
	Port                   uint
	DatabaseMaxIdle        int
	DatabaseMaxConnections int
}

// Database holds the connection settings for MySQL
type Database struct {
	Host     string
	Protocol string
	User     string
	Password string
	Database string
}

// Directories sets where files will be stored locally
type Directories struct {
	AssetsDir string
}

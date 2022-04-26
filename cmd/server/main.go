package main

import (
	"flag"
	"fmt"
	"os"

	"jane_tech/internal/config"
	"jane_tech/internal/database"
	"jane_tech/internal/logger"
	"jane_tech/internal/server"

	log "github.com/sirupsen/logrus"

	"go.uber.org/dig"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "configs/config.yaml", "Used for set path to config file.")
	flag.Parse()

	c := dig.New()
	err := c.Provide(func() (*config.Config, error) {
		return config.Configure(configPath)
	})
	if err != nil {
		fmt.Println("Failed to create config", err)
		os.Exit(1)
	}
	err = c.Provide(func(cfg *config.Config) *log.Logger {
		return logger.CreateLogger(cfg.Log)
	})
	if err != nil {
		fmt.Println("Failed to create logger", err)
		os.Exit(1)
	}
	err = c.Provide(database.NewDatabaseConnector)
	if err != nil {
		fmt.Println("Failed to create new db connection", err)
		os.Exit(1)
	}
	err = c.Provide(server.NewServer)
	if err != nil {
		fmt.Println("Failed to create server", err)
		os.Exit(1)
	}

	err = c.Invoke(func(server *server.Server) {
		server.Run()
	})
	if err != nil {
		fmt.Println("Failed to run server", err)
		os.Exit(1)
	}
}

package main

import (
	"go-fs/internal/conf"
	"go-fs/internal/server"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	data, err := os.ReadFile("./etc/server_config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	cfg := &conf.ServerConfig{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalln(err)
	}
	server := server.NewServer(cfg)
	panic(server.Run())
}

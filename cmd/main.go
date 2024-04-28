package main

import (
	"go-fs/internal/conf"
	"go-fs/internal/repo"
	"go-fs/internal/server"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	data, err := os.ReadFile("./etc/server_config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	cfg := &conf.Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalln(err)
	}
	server := server.NewServer(cfg.ServerCfg)
	repo.SetDataBaseConfig(cfg.ReadDBCfg)

	panic(server.Run())
}

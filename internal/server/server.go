package server

import (
	"go-fs/common/geeweb"
	"go-fs/internal/conf"
	"go-fs/internal/handler"
	"net"
)

type Server struct {
	cfg *conf.ServerConfig
	e   *geeweb.Engine
}

func NewServer(cfg *conf.ServerConfig) *Server {
	return &Server{
		cfg: cfg,
		e:   geeweb.NewEngine(),
	}
}

func (s *Server) Run() error {
	handler.RegistryHandleFunc(s.e.RouterGroup)
	lis, err := net.Listen(s.cfg.Network, s.cfg.Host+":"+s.cfg.Port)
	if err != nil {
		return err
	}
	return s.e.Serve(lis)
}

package apiserver

import (
	"fmt"

	genericclioptions "github.com/onexstack_practice/fast_blog/pkg/options"
)

// Config存储应用配置
type Config struct {
	MysqlOptions *genericclioptions.MysqlOptions
}

// Server是一个服务器结构体类型
type Server struct {
	Config *Config
}

func (cfg *Config) NewServer() *Server {
	return &Server{
		Config: cfg,
	}
}

func (s *Server) Run() error {
	fmt.Printf("Read Mysql addr from Viper: %s\n\n", s.Config.MysqlOptions.Addr)

	select {}
}

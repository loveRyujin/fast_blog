package apiserver

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	mw "github.com/onexstack_practice/fast_blog/internal/pkg/middleware"
	genericclioptions "github.com/onexstack_practice/fast_blog/pkg/options"
)

// Config存储应用配置
type Config struct {
	MysqlOptions *genericclioptions.MysqlOptions
	Addr         string
}

// Server是一个服务器结构体类型
type Server struct {
	Config *Config
	srv    *http.Server
}

func (cfg *Config) NewServer() *Server {
	engine := gin.New()

	// 注册全局中间件
	middlewares := []gin.HandlerFunc{gin.Recovery(), mw.NoCache(), mw.Cors(), mw.RequestID()}
	engine.Use(middlewares...)

	// 注册404Handler
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "404 Not Found."})
	})

	// 注册/healthz路由
	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 创建httpServer实例
	httpServer := &http.Server{
		Addr:    cfg.Addr,
		Handler: engine,
	}

	return &Server{
		Config: cfg,
		srv:    httpServer,
	}
}

func (s *Server) Run() error {
	slog.Info("Read Mysql addr from Viper", "mysql.addr", s.Config.MysqlOptions.Addr)
	slog.Info("Start to listen the incoming request on http address", "addr", s.Config.Addr)

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

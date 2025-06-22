package apiserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/fast_blog/internal/apiserver/biz"
	"github.com/loveRyujin/fast_blog/internal/apiserver/handler"
	"github.com/loveRyujin/fast_blog/internal/apiserver/pkg/validation"
	"github.com/loveRyujin/fast_blog/internal/apiserver/store"
	"github.com/loveRyujin/fast_blog/internal/pkg/core"
	"github.com/loveRyujin/fast_blog/internal/pkg/errorx"
	"github.com/loveRyujin/fast_blog/internal/pkg/known"
	"github.com/loveRyujin/fast_blog/internal/pkg/log"
	mw "github.com/loveRyujin/fast_blog/internal/pkg/middleware"
	genericclioptions "github.com/loveRyujin/fast_blog/pkg/options"
	"github.com/loveRyujin/fast_blog/pkg/token"
)

const (
	GRPCServerMode        = "grpc"
	HTTPServerMode        = "http"
	GRPCGatewayServerMode = "grpc-gateway"
)

// Config存储应用配置
type Config struct {
	ServerMode   string
	MysqlOptions *genericclioptions.MysqlOptions
	HTTPOptions  *genericclioptions.HTTPOptions
	GRPCOptions  *genericclioptions.GRPCOptions
	JWTKey       string
	Expiration   time.Duration
}

// Server是一个服务器结构体类型
type Server struct {
	Config *Config
	srv    *http.Server
}

func (cfg *Config) NewServer() (*Server, error) {
	// 初始化 JWT token
	token.Init(cfg.JWTKey, known.XUserID, cfg.Expiration)
	// 创建gin引擎
	engine := gin.New()

	// 注册全局中间件
	middlewares := []gin.HandlerFunc{gin.Recovery(), mw.NoCache(), mw.Cors(), mw.RequestID()}
	engine.Use(middlewares...)

	// 初始化数据库连接
	db, err := cfg.MysqlOptions.NewDB()
	if err != nil {
		return nil, err
	}
	store := store.NewStore(db)

	cfg.SetupRouter(engine, store)

	// 创建httpServer实例
	httpServer := &http.Server{
		Addr:    cfg.HTTPOptions.Addr,
		Handler: engine,
	}

	return &Server{
		Config: cfg,
		srv:    httpServer,
	}, nil
}

// 注册 API 路由。路由的路径和 HTTP 方法，严格遵循 REST 规范.
func (cfg *Config) SetupRouter(engine *gin.Engine, store store.IStore) {
	// 注册 404 Handler.
	engine.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errorx.ErrNotFound.WithMessage("Page not found"), nil)
	})

	// 注册 /healthz handler.
	engine.GET("/healthz", func(c *gin.Context) {
		core.WriteResponse(c, nil, gin.H{"status": "ok"})
	})

	// 创建核心业务处理器
	handler := handler.NewHandler(biz.NewBiz(store), validation.NewValidator(store))

	engine.POST("/login", handler.Login)
	engine.POST("/refresh-token", mw.Authn(), handler.RefreshToken)

	authMiddlewares := []gin.HandlerFunc{mw.Authn()}

	// 注册 v1 版本 API 路由分组
	v1 := engine.Group("/v1")
	{
		// 用户相关路由
		userv1 := v1.Group("/users")
		{
			// 创建用户。这里要注意：创建用户是不用进行认证和授权的
			userv1.POST("", handler.CreateUser)
			userv1.Use(authMiddlewares...)
			userv1.PUT(":userID/change-password", handler.ChangePassword) // 修改用户密码
			userv1.PUT(":userID", handler.UpdateUser)                     // 更新用户信息
			userv1.DELETE(":userID", handler.DeleteUser)                  // 删除用户
			userv1.GET(":userID", handler.GetUser)                        // 查询用户详情
			userv1.GET("", handler.ListUser)                              // 查询用户列表
		}

		// 博客相关路由
		postv1 := v1.Group("/posts", authMiddlewares...)
		{
			postv1.POST("", handler.CreatePost)       // 创建博客
			postv1.PUT(":postID", handler.UpdatePost) // 更新博客
			postv1.DELETE("", handler.DeletePost)     // 删除博客
			postv1.GET(":postID", handler.GetPost)    // 查询博客详情
			postv1.GET("", handler.ListPost)          // 查询博客列表
		}
	}
}

func (s *Server) Run() error {
	log.Infow("Read Mysql addr from Viper", "mysql.addr", s.Config.MysqlOptions.Addr)
	log.Infow("Start to listen the incoming request on http address", "addr", s.Config.HTTPOptions.Addr)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infow("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 先关闭依赖的服务，再关闭被依赖的服务
	// 10s内关闭所有服务，超过10s就超时退出
	if err := s.srv.Shutdown(ctx); err != nil {
		slog.Error("Insecure Server forced to shutdown:", "err", err)
		return err
	}

	log.Infow("Server exited")

	return nil
}

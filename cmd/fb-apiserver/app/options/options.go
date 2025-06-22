package options

import (
	"fmt"
	"time"

	"github.com/loveRyujin/fast_blog/internal/apiserver"
	genericoptions "github.com/loveRyujin/fast_blog/pkg/options"
	"k8s.io/apimachinery/pkg/util/sets"
)

var availableServerOptions = sets.New(
	apiserver.GRPCServerMode,
	apiserver.HTTPServerMode,
	apiserver.GRPCGatewayServerMode,
)

type ServerOptions struct {
	ServerMode   string                       `json:"server-mode" mapstructure:"server-mode"` // 服务器模式，支持grpc、http、grpc-gateway
	MysqlOptions *genericoptions.MysqlOptions `json:"mysql" mapstructure:"mysql"`
	GRPCOptions  *genericoptions.GRPCOptions  `json:"grpc" mapstructure:"grpc"`
	HTTPOptions  *genericoptions.HTTPOptions  `json:"http" mapstructure:"http"`
	JWTKey       string                       `json:"jwt-key" mapstructure:"jwt-key"`
	Expiration   time.Duration                `json:"expiration" mapstructure:"expiration"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		ServerMode:   apiserver.GRPCGatewayServerMode,
		MysqlOptions: genericoptions.NewMysqlOptions(),
		GRPCOptions:  genericoptions.NewGRPCOptions(),
		HTTPOptions:  genericoptions.NewHTTPOptions(),
		Expiration:   2 * time.Hour,
	}
}

// 校验ServerOptions配置
func (o *ServerOptions) Validate() error {
	if !availableServerOptions.Has(o.ServerMode) {
		fmt.Printf("invalid server mode: %s, available modes: %v\n", o.ServerMode, availableServerOptions.UnsortedList())
	}

	if err := o.MysqlOptions.Validate(); err != nil {
		return err
	}

	if err := o.HTTPOptions.Validate(); err != nil {
		return err
	}

	if err := o.GRPCOptions.Validate(); err != nil {
		return err
	}

	return nil
}

// Config 基于ServerOptions配置生成apiserver.Config
func (o *ServerOptions) Config() *apiserver.Config {
	return &apiserver.Config{
		ServerMode:   o.ServerMode,
		MysqlOptions: o.MysqlOptions,
		HTTPOptions:  o.HTTPOptions,
		GRPCOptions:  o.GRPCOptions,
		JWTKey:       o.JWTKey,
		Expiration:   o.Expiration,
	}
}

package options

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/onexstack_practice/fast_blog/internal/apiserver"
	genericoptions "github.com/onexstack_practice/fast_blog/pkg/options"
)

type ServerOptions struct {
	MysqlOptions *genericoptions.MysqlOptions `json:"mysql" mapstructure:"mysql"`
	Addr         string                       `json:"addr" mapstructure:"addr"`
	JWTKey       string                       `json:"jwt-key" mapstructure:"jwt-key"`
	Expiration   time.Duration                `json:"expiration" mapstructure:"expiration"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		MysqlOptions: genericoptions.NewMysqlOptions(),
		Addr:         "0.0.0.0:6666",
		Expiration:   2 * time.Hour,
	}
}

// 校验ServerOptions配置
func (o *ServerOptions) Validate() error {
	if err := o.MysqlOptions.Validate(); err != nil {
		return err
	}

	// 校验服务器地址
	if o.Addr == "" {
		return errors.New("server addr cannot be empty")
	}
	_, portStr, err := net.SplitHostPort(o.Addr)
	if err != nil {
		return fmt.Errorf("invalid server addr format: '%s': %v", o.Addr, err)
	}
	// 校验服务器端口是否为数字且值是否合理
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		return fmt.Errorf("invalid server port: %s", portStr)
	}

	return nil
}

// Config 基于ServerOptions配置生成apiserver.Config
func (o *ServerOptions) Config() *apiserver.Config {
	return &apiserver.Config{
		MysqlOptions: o.MysqlOptions,
		Addr:         o.Addr,
		JWTKey:       o.JWTKey,
		Expiration:   o.Expiration,
	}
}

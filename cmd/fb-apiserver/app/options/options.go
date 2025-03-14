package options

import (
	"github.com/onexstack_practice/fast_blog/internal/apiserver"
	genericoptions "github.com/onexstack_practice/fast_blog/pkg/options"
)

type ServerOptions struct {
	MysqlOptions *genericoptions.MysqlOptions `json:"mysql" mapstructure:"mysql"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		MysqlOptions: genericoptions.NewMysqlOptions(),
	}
}

// 校验ServerOptions配置
func (o *ServerOptions) Validate() error {
	if err := o.MysqlOptions.Validate(); err != nil {
		return err
	}

	return nil
}

func (o *ServerOptions) Config() *apiserver.Config {
	return &apiserver.Config{
		MysqlOptions: o.MysqlOptions,
	}
}

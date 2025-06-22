package options

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

type HTTPOptions struct {
	Addr    string        `json:"addr" mapstructure:"addr"`       // HTTP服务器地址
	Timeout time.Duration `json:"timeout" mapstructure:"timeout"` // HTTP请求超时时间，默认是30秒
}

func NewHTTPOptions() *HTTPOptions {
	return &HTTPOptions{
		Addr:    "0.0.0.0:6666",
		Timeout: 30 * time.Second,
	}
}

func (o *HTTPOptions) Validate() error {
	// 校验服务器地址
	if o.Addr == "" {
		return errors.New("http server addr cannot be empty")
	}
	_, portStr, err := net.SplitHostPort(o.Addr)
	if err != nil {
		return fmt.Errorf("invalid http server addr format: '%s': %v", o.Addr, err)
	}
	// 校验服务器端口是否为数字且值是否合理
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		return fmt.Errorf("invalid http server port: %s", portStr)
	}

	return nil
}

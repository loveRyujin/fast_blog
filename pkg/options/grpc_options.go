package options

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

type GRPCOptions struct {
	Addr    string        `json:"addr" mapstructure:"addr"`       // gRPC服务器地址，格式为 "host:port"
	Timeout time.Duration `json:"timeout" mapstructure:"timeout"` // gRPC请求超时时间，默认是30秒
}

func NewGRPCOptions() *GRPCOptions {
	return &GRPCOptions{
		Addr:    "0.0.0.0:39090",
		Timeout: 30 * time.Second,
	}
}

func (o *GRPCOptions) Validate() error {
	// 校验服务器地址
	if o.Addr == "" {
		return errors.New("grpc server addr cannot be empty")
	}
	_, portStr, err := net.SplitHostPort(o.Addr)
	if err != nil {
		return fmt.Errorf("invalid grpc server addr format: '%s': %v", o.Addr, err)
	}
	// 校验服务器端口是否为数字且值是否合理
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		return fmt.Errorf("invalid grpc server port: %s", portStr)
	}

	return nil
}

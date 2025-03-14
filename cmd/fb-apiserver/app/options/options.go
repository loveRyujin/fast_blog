package options

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	MINPORTNUM = 1
	MAXPORTNUM = 65535
)

type MysqlOptions struct {
	Addr                  string        `json:"addr,omitempty" mapstructure:"addr"`
	Username              string        `json:"username,omitempty" mapstructure:"username"`
	Password              string        `json:"-" mapstructure:"password"`
	Database              string        `json:"database" mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty" mapstructure:"max-idle-connections,omitempty"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty" mapstructure:"max-open-connections,omitempty"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-lifetime,omitempty" mapstructure:"max-connection-life-time"`
}

func NewMysqlOptions() *MysqlOptions {
	return &MysqlOptions{
		Addr:                  "localhost:3306",
		Username:              "ryujin",
		Password:              "LCSZBD20030208lc!",
		Database:              "fastgo",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
	}
}

type ServerOptions struct {
	MysqlOptions *MysqlOptions `json:"mysql" mapstructure:"mysql"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		MysqlOptions: NewMysqlOptions(),
	}
}

func (o *ServerOptions) Validate() error {
	// 校验mysql地址格式
	if o.MysqlOptions.Addr == "" {
		return fmt.Errorf("mysql.addr is required")
	}
	// 检验mysql地址格式，必须为host:port格式
	host, portStr, err := net.SplitHostPort(o.MysqlOptions.Addr)
	if err != nil {
		return fmt.Errorf("mysql addr wrong format: %s: %v", o.MysqlOptions.Addr, err)
	}
	// 校验主机号
	if host == "" {
		return fmt.Errorf("mysql addr host is required")
	}
	// 校验端口号
	port, err := strconv.Atoi(portStr)
	if err != nil || port < MINPORTNUM || port > MAXPORTNUM {
		return fmt.Errorf("mysql addr port is invalid: %s", portStr)
	}

	// 校验mysql用户名、密码、数据库
	if o.MysqlOptions.Username == "" {
		return fmt.Errorf("mysql.username is required")
	}
	if o.MysqlOptions.Password == "" {
		return fmt.Errorf("mysql.password is required")
	}
	if o.MysqlOptions.Database == "" {
		return fmt.Errorf("mysql.database is required")
	}

	// 校验mysql连接池配置
	if o.MysqlOptions.MaxIdleConnections <= 0 {
		return fmt.Errorf("mysql.max-idle-connections must be greater than 0")
	}
	if o.MysqlOptions.MaxOpenConnections <= 0 {
		return fmt.Errorf("mysql.max-open-connections must be greater than 0")
	}
	if o.MysqlOptions.MaxOpenConnections < o.MysqlOptions.MaxIdleConnections {
		return fmt.Errorf("mysql.max-open-connections must be greater than or equal to mysql.max-idle-connections")
	}
	if o.MysqlOptions.MaxConnectionLifeTime <= 0 {
		return fmt.Errorf("mysql.max-connection-lifetime must be greater than 0")
	}

	return nil
}

package options

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

// 校验mysql配置
func (o *MysqlOptions) Validate() error {
	// 校验mysql地址格式
	if o.Addr == "" {
		return fmt.Errorf("mysql.addr is required")
	}
	// 检验mysql地址格式，必须为host:port格式
	host, portStr, err := net.SplitHostPort(o.Addr)
	if err != nil {
		return fmt.Errorf("mysql addr wrong format: %s: %v", o.Addr, err)
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
	if o.Username == "" {
		return fmt.Errorf("mysql.username is required")
	}
	if o.Password == "" {
		return fmt.Errorf("mysql.password is required")
	}
	if o.Database == "" {
		return fmt.Errorf("mysql.database is required")
	}

	// 校验mysql连接池配置
	if o.MaxIdleConnections <= 0 {
		return fmt.Errorf("mysql.max-idle-connections must be greater than 0")
	}
	if o.MaxOpenConnections <= 0 {
		return fmt.Errorf("mysql.max-open-connections must be greater than 0")
	}
	if o.MaxOpenConnections < o.MaxIdleConnections {
		return fmt.Errorf("mysql.max-open-connections must be greater than or equal to mysql.max-idle-connections")
	}
	if o.MaxConnectionLifeTime <= 0 {
		return fmt.Errorf("mysql.max-connection-lifetime must be greater than 0")
	}

	return nil
}

// DSN return DSN from MySQLOptions.
func (o *MysqlOptions) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Addr,
		o.Database,
		true,
		"Local")
}

// NewDB create mysql store with the given config.
func (o *MysqlOptions) NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(o.DSN()), &gorm.Config{
		// PrepareStmt executes the given query in cached statement.
		// This can improve performance.
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(o.MaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(o.MaxConnectionLifeTime)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(o.MaxIdleConnections)

	return db, nil
}

package app

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/loveRyujin/fast_blog/cmd/fb-apiserver/app/options"
	"github.com/loveRyujin/fast_blog/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configPath string

func NewFastBlogCommand() *cobra.Command {
	// 创建默认的应用命令行选项
	opts := options.NewServerOptions()

	cmd := &cobra.Command{
		// 指定命令名称
		Use:   "fb-apiserver",
		Short: "FastBlog API Server",
		Long:  "FastBlog API Server, a simple blog API server",
		// 命令出错时，不打印帮助信息，保证出错时可以一眼看见错误信息
		SilenceUsage: true,
		// 指定调用cmd.Execute()时的执行的Run函数
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
		// 设置命令的参数检查，无需传入参数
		Args: cobra.NoArgs,
	}

	// 初始化配置函数，每个命令执行时调用
	cobra.OnInitialize(OnInitialize)

	// 将命令行参数解析到变量当中
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to fb-apiserver confuguration file.")

	// 增加--version标志
	version.AddFlags(cmd.PersistentFlags())

	return cmd
}

func run(opts *options.ServerOptions) error {
	// 如果命令行参数中包含--version，则打印版本信息并退出
	version.PrintAndExitIfRequested()

	// 初始化日志
	initLog()

	// 从配置文件中读取配置
	if err := viper.Unmarshal(opts); err != nil {
		fmt.Printf("读取配置文件失败: %v", err)
		return err
	}

	// 验证配置
	if err := opts.Validate(); err != nil {
		fmt.Printf("验证配置失败: %v", err)
		return err
	}

	// 获取应用配置.
	// 将命令行选项和应用配置分开，方便灵活处理两种不同类型的配置
	cfg := opts.Config()

	// 创建服务器实例
	server, err := cfg.NewServer()
	if err != nil {
		return err
	}

	// 启动服务器
	return server.Run()
}

func initLog() {
	slevel := getLogLevel()
	sopts := &slog.HandlerOptions{Level: slevel}
	soutput := getLogOutput()
	handler := getLogHandler(soutput, sopts)
	slog.SetDefault(slog.New(handler))
}

// getLogLevel 从配置文件中获取日志级别
func getLogLevel() slog.Level {
	level := viper.GetString("log.level")
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// getLogOutput 从配置文件中获取日志输出
func getLogOutput() io.Writer {
	output := viper.GetString("log.output")
	switch output {
	case "":
		return os.Stdout
	case "stdout":
		return os.Stdout
	default:
		fd, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		return fd
	}
}

// getLogHandler 根据配置文件配置获取日志处理器
func getLogHandler(output io.Writer, opts *slog.HandlerOptions) slog.Handler {
	format := viper.GetString("log.format")
	switch format {
	case "json":
		return slog.NewJSONHandler(output, opts)
	case "text":
		return slog.NewTextHandler(output, opts)
	default:
		return slog.NewJSONHandler(output, opts)
	}
}

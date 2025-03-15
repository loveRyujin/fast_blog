package app

import (
	"fmt"

	"github.com/onexstack_practice/fast_blog/cmd/fb-apiserver/app/options"
	"github.com/onexstack_practice/fast_blog/pkg/version"
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
	server := cfg.NewServer()

	// 启动服务器
	return server.Run()
}

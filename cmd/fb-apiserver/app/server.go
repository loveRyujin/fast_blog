package app

import (
	"encoding/json"
	"fmt"

	"github.com/onexstack_practice/fast_blog/cmd/fb-apiserver/app/options"
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

			fmt.Printf("Read Mysql addr from Viper: %s\n\n", viper.GetString("mysql.addr"))
			fmt.Printf("Read Mysql username from Viper: %s\n\n", viper.GetString("mysql.username"))

			jsonData, err := json.MarshalIndent(opts, "", "  ")
			if err != nil {
				fmt.Printf("序列化配置失败: %v", err)
				return err
			}
			fmt.Println(string(jsonData))

			return nil
		},
		// 设置命令的参数检查，无需传入参数
		Args: cobra.NoArgs,
	}

	cobra.OnInitialize(OnInitialize)

	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to fb-apiserver confuguration file.")

	return cmd
}

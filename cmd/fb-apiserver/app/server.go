package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewFastBlogCommand() *cobra.Command {
	cmd := &cobra.Command{
		// 指定命令名称
		Use:   "fb-apiserver",
		Short: "FastBlog API Server",
		Long:  "FastBlog API Server, a simple blog API server",
		// 命令出错时，不打印帮助信息，保证出错时可以一眼看见错误信息
		SilenceUsage: true,
		// 指定调用cmd.Execute()时的执行的Run函数
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Hello, FastBlog!")
			return nil
		},
		// 设置命令的参数检查，无需传入参数
		Args: cobra.NoArgs,
	}

	return cmd

}

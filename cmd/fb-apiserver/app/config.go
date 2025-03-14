package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// defaultConfigDir用来存放默认的配置文件目录
	defaultConfigDir = ".fastblog"
	// defaultConfigFile用来存放默认的配置文件
	defaultConfigFile = "fb-apiserver.yaml"
)

func OnInitialize() {
	if configPath != "" {
		// 从命令行指定的参数路径下读取配置文件
		viper.SetConfigFile(configPath)
	} else {
		for _, dir := range searchDirs() {
			// 将dir添加到viper的配置文件搜索路径中
			viper.AddConfigPath(dir)
		}

		// 设置viper读取的配置文件类型
		viper.SetConfigType("yaml")

		// 设置viper读取的文件名称
		viper.SetConfigName(defaultConfigFile)
	}
	// 读取环境变量并设置前缀
	setupEnvironmentVariables()

	// 读取配置文件，如果配置文件不存在则在注册的搜索路径中查找
	if err := viper.ReadInConfig(); err != nil {
		cobra.CheckErr(err)
	}
}

func setupEnvironmentVariables() {
	// 允许viper自动匹配环境变量
	viper.AutomaticEnv()
	// 设置环境变量前缀
	viper.SetEnvPrefix("FASTBLOG")
	// 设置环境变量的key替换规则
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)
}

// searchDirs返回配置文件的搜索目录
func searchDirs() []string {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	// 如果获取用户主目录失败，则打印错误并退出程序
	cobra.CheckErr(err)
	return []string{filepath.Join(homeDir, defaultConfigDir), ".", "./configs"}
}

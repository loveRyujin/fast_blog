package log

import "go.uber.org/zap"

// Options 定义了日志配置的选项结构体.
// 通过该结构体，可以自定义日志的输出格式、级别以及其他相关配置.
type Options struct {
	// Level 指定日志级别.
	// 可选值包括：debug、info、warn、error、dpanic、panic、fatal.
	// 默认值为 info.
	Level string

	// Format 指定日志的输出格式.
	// 可选值包括：console（控制台格式）和 json（JSON 格式）.
	// 默认值为 console.
	Format string

	// Output 指定日志的输出位置.
	// 默认值为标准输出（stdout），也可以指定文件路径或其他输出目标.
	Output []string

	// CallerEnabled 指定是否开启 caller 信息.
	// 如果设置为 true（默认值），日志中会显示调用日志所在的文件名和行号，例如："caller":"main.go:42".
	CallerEnabled bool

	// StacktraceEnabled 指定是否开启堆栈追踪.
	// 如果设置为 true（默认值），在日志级别为 panic 或更高时，会打印堆栈跟踪信息.
	StacktraceEnabled bool
}

func NewOptions() *Options {
	return &Options{
		Level:             zap.InfoLevel.String(),
		Format:            "console",
		Output:            []string{"stdout"},
		CallerEnabled:     true,
		StacktraceEnabled: true,
	}
}

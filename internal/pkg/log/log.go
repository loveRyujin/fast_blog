package log

import (
	"context"
	"sync"
	"time"

	"github.com/loveRyujin/fast_blog/internal/pkg/contextx"
	"github.com/loveRyujin/fast_blog/internal/pkg/known"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 定义了 fast_blog 项目的日志接口。
// 该接口包含了项目中支持的日志记录方法，提供对不同日志级别的支持。
type Logger interface {
	// Debugw 用于记录调试级别的日志，通常用于开发阶段，包含详细的调试信息。
	Debugw(msg string, kvs ...any)

	// Infow 用于记录信息级别的日志，表示系统的正常运行状态。
	Infow(msg string, kvs ...any)

	// Warnw 用于记录警告级别的日志，表示可能存在问题但不影响系统正常运行。
	Warnw(msg string, kvs ...any)

	// Errorw 用于记录错误级别的日志，表示系统运行中出现的错误，需要开发人员介入处理。
	Errorw(msg string, kvs ...any)

	// Panicw 用于记录严重错误级别的日志，表示系统无法继续运行，记录日志后会触发 panic。
	Panicw(msg string, kvs ...any)

	// Fatalw 用于记录致命错误级别的日志，表示系统无法继续运行，记录日志后会直接退出程序。
	Fatalw(msg string, kvs ...any)

	// Sync 用于刷新日志缓冲区，确保日志被完整写入目标存储。
	Sync()
}

type zapLogger struct {
	z *zap.Logger
}

var _ Logger = (*zapLogger)(nil)

var (
	once sync.Once

	def = New(NewOptions())
)

func Init(opts *Options) {
	once.Do(func() {
		def = New(opts)
	})
}

func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	// 将Options里的level转换为zapcore.Level
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	// 创建 encoder 配置，用于控制日志的输出格式
	encoderCfg := zap.NewProductionEncoderConfig()
	// 自定义TimeKey，明确语义
	encoderCfg.TimeKey = "timestamp"
	// 自定义MessageKey，明确语义
	encoderCfg.MessageKey = "message"
	// 指定时间序列化函数，将时间序列化为 `2006-01-02 15:04:05.000` 格式，更易读
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// 指定 time.Duration 序列化函数，将 time.Duration 序列化为经过的毫秒数的浮点数
	// 毫秒数比默认的秒数更精确
	encoderCfg.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}

	// 构建zap logger配置
	cfg := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Encoding:          opts.Format,
		EncoderConfig:     encoderCfg,
		OutputPaths:       opts.Output,
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     !opts.CallerEnabled,
		DisableStacktrace: !opts.StacktraceEnabled,
	}

	// 使用zap config创建logger
	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}

	zap.RedirectStdLog(z)

	return &zapLogger{z: z}
}

func Debugw(msg string, kvs ...any) {
	def.Debugw(msg, kvs...)
}

func (l *zapLogger) Debugw(msg string, kvs ...any) {
	l.z.Sugar().Debugw(msg, kvs...)
}

func Infow(msg string, kvs ...any) {
	def.Infow(msg, kvs...)
}

func (l *zapLogger) Infow(msg string, kvs ...any) {
	l.z.Sugar().Infow(msg, kvs...)
}

func Warnw(msg string, kvs ...any) {
	def.Warnw(msg, kvs...)
}

func (l *zapLogger) Warnw(msg string, kvs ...any) {
	l.z.Sugar().Warnw(msg, kvs...)
}

func Errorw(msg string, kvs ...any) {
	def.Errorw(msg, kvs...)
}

func (l *zapLogger) Errorw(msg string, kvs ...any) {
	l.z.Sugar().Errorw(msg, kvs...)
}

func Panicw(msg string, kvs ...any) {
	def.Panicw(msg, kvs...)
}

func (l *zapLogger) Panicw(msg string, kvs ...any) {
	l.z.Sugar().Panicw(msg, kvs...)
}

func Fatalw(msg string, kvs ...any) {
	def.Fatalw(msg, kvs...)
}

func (l *zapLogger) Fatalw(msg string, kvs ...any) {
	l.z.Sugar().Fatalw(msg, kvs...)
}

func Sync() {
	def.Sync()
}

func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}

func With(ctx context.Context) Logger {
	return def.With(ctx)
}

// With 返回一个新的 Logger 实例，且日志输出包含从上下文中提取的请求 ID 和用户 ID。
func (l *zapLogger) With(ctx context.Context) Logger {
	cl := l.clone()

	extractorMap := map[string]func(ctx context.Context) string{
		known.XRequestID: contextx.RequestID,
		known.XUserID:    contextx.UserID,
	}

	for key, extractor := range extractorMap {
		if value := extractor(ctx); value != "" {
			cl.z = cl.z.With(zap.String(key, value))
		}
	}

	return cl
}

func (l *zapLogger) clone() *zapLogger {
	newLogger := *l
	return &newLogger
}

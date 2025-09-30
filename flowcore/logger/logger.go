package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	flowcorepkg "github.com/hanyougame/glib/flowcore/pkg"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ========================= 接口定义 =========================

// Logger Temporal 日志接口
type Logger interface {
	log.Logger
	WithFields(fields map[string]interface{}) Logger
}

// ========================= 基础 Logger 实现 =========================

// ZapLogger 使用 zap 实现 Temporal 的日志接口
type ZapLogger struct {
	logger *zap.Logger
	level  LogLevel
	fields map[string]interface{} // 附加字段
}

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// ========================= 颜色定义 =========================

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"

	// 粗体颜色组合
	ColorRedBold     = "\033[1;31m"
	ColorGreenBold   = "\033[1;32m"
	ColorYellowBold  = "\033[1;33m"
	ColorBlueBold    = "\033[1;34m"
	ColorMagentaBold = "\033[1;35m"
	ColorCyanBold    = "\033[1;36m"
)

var (
	projectRootOnce sync.Once
	projectRootPath string
	enableColor     = isColorTerminal()
)

// ========================= 配置 =========================

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level       string                 // debug, info, warn, error
	WithIcon    bool                   // 是否显示图标
	WithColor   bool                   // 是否启用颜色
	Format      string                 // console, json
	ShowCaller  bool                   // 是否显示调用者信息
	ProjectRoot string                 // 项目根目录
	Fields      map[string]interface{} // 默认附加字段
}

// 默认配置
func defaultConfig() LoggerConfig {
	return LoggerConfig{
		Level:      "info",
		WithIcon:   true,
		WithColor:  true,
		Format:     "console",
		ShowCaller: true,
		Fields:     make(map[string]interface{}),
	}
}

// ========================= 工厂方法 =========================

// NewLogger 创建默认日志记录器
func NewLogger(level string) Logger {
	config := defaultConfig()
	config.Level = level
	return NewLoggerWithConfig(config)
}

// NewLoggerWithConfig 使用配置创建 logger
func NewLoggerWithConfig(config LoggerConfig) Logger {
	logLevel := parseLogLevel(config.Level)

	if config.ProjectRoot != "" {
		projectRootPath = config.ProjectRoot
	}

	if config.WithColor {
		enableColor = true
	}

	zapLogger := createZapLogger(logLevel, config)

	return &ZapLogger{
		logger: zapLogger,
		level:  logLevel,
		fields: config.Fields,
	}
}

// NewWorkflowLogger 获取带 trace 的 Workflow Logger
func NewWorkflowLogger(ctx workflow.Context) Logger {
	logger := workflow.GetLogger(ctx)
	info := workflow.GetInfo(ctx)

	// 提取 trace 信息
	traceID := ""
	spanID := ""

	// 从 context 中获取 trace 信息
	if val := ctx.Value("trace"); val != nil {
		traceID, _ = val.(string)
	}
	if val := ctx.Value("span"); val != nil {
		spanID, _ = val.(string)
	}

	// 如果 context 中没有，尝试从 Memo 获取
	if traceID == "" && info.Memo != nil {
		contextData := flowcorepkg.ExtractContextFromMemo(info.Memo)
		traceID = contextData.TraceID
		spanID = contextData.SpanID
	}

	// 创建带附加字段的 logger
	fields := map[string]interface{}{
		"workflow_id":   info.WorkflowExecution.ID,
		"run_id":        info.WorkflowExecution.RunID,
		"workflow_type": info.WorkflowType.Name,
	}

	if traceID != "" {
		fields["trace"] = traceID
		fields["span"] = spanID
	}

	// 包装原生 logger
	return wrapTemporalLogger(logger, fields)
}

// NewActivityLogger 获取带 trace 的 Activity Logger
func NewActivityLogger(ctx context.Context) Logger {
	logger := activity.GetLogger(ctx)
	info := activity.GetInfo(ctx)

	// 创建带附加字段的 logger
	fields := map[string]interface{}{
		"activity_type": info.ActivityType.Name,
		"workflow_id":   info.WorkflowExecution.ID,
		"activity_id":   info.ActivityID,
	}

	// 包装原生 logger
	return wrapTemporalLogger(logger, fields)
}

// ========================= 便捷方法 =========================

// NewDevelopmentLogger 开发环境日志记录器
func NewDevelopmentLogger(level string) Logger {
	return NewLoggerWithConfig(LoggerConfig{
		Level:      level,
		WithIcon:   true,
		WithColor:  true,
		Format:     "console",
		ShowCaller: true,
	})
}

// NewProductionLogger 生产环境日志记录器
func NewProductionLogger() Logger {
	return NewLoggerWithConfig(LoggerConfig{
		Level:      "info",
		WithIcon:   false,
		WithColor:  false,
		Format:     "json",
		ShowCaller: true,
	})
}

// ========================= Logger 接口实现 =========================

// Debug 记录调试信息
func (l *ZapLogger) Debug(msg string, keyvals ...interface{}) {
	if l.level <= LevelDebug {
		fields := l.parseKeyvals(keyvals...)
		fields = append(fields, l.getDefaultFields()...)
		fields = append(fields, l.getCallerField())
		l.logger.Debug(msg, fields...)
	}
}

// Info 记录信息消息
func (l *ZapLogger) Info(msg string, keyvals ...interface{}) {
	if l.level <= LevelInfo {
		fields := l.parseKeyvals(keyvals...)
		fields = append(fields, l.getDefaultFields()...)
		fields = append(fields, l.getCallerField())
		l.logger.Info(msg, fields...)
	}
}

// Warn 记录警告消息
func (l *ZapLogger) Warn(msg string, keyvals ...interface{}) {
	if l.level <= LevelWarn {
		fields := l.parseKeyvals(keyvals...)
		fields = append(fields, l.getDefaultFields()...)
		fields = append(fields, l.getCallerField())
		l.logger.Warn(msg, fields...)
	}
}

// Error 记录错误消息
func (l *ZapLogger) Error(msg string, keyvals ...interface{}) {
	if l.level <= LevelError {
		fields := l.parseKeyvals(keyvals...)
		fields = append(fields, l.getDefaultFields()...)
		fields = append(fields, l.getCallerField())
		l.logger.Error(msg, fields...)
	}
}

// WithFields 创建带附加字段的新日志记录器
func (l *ZapLogger) WithFields(fields map[string]interface{}) Logger {
	// 合并字段
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}

	// 创建新的 logger
	zapFields := make([]zap.Field, 0, len(newFields))
	for k, v := range newFields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	newLogger := l.logger.With(zapFields...)
	return &ZapLogger{
		logger: newLogger,
		level:  l.level,
		fields: newFields,
	}
}

// ========================= 内部辅助方法 =========================

// wrapTemporalLogger 包装 Temporal 原生 logger
func wrapTemporalLogger(logger log.Logger, fields map[string]interface{}) Logger {
	// 如果已经是我们的 Logger，直接添加字段
	if zl, ok := logger.(*ZapLogger); ok {
		return zl.WithFields(fields)
	}

	// 否则创建一个包装器
	return &temporalLoggerWrapper{
		Logger: logger,
		fields: fields,
	}
}

// temporalLoggerWrapper 包装 Temporal 原生 logger
type temporalLoggerWrapper struct {
	log.Logger
	fields map[string]interface{}
}

func (w *temporalLoggerWrapper) Debug(msg string, keyvals ...interface{}) {
	keyvals = w.appendFields(keyvals)
	w.Logger.Debug(msg, keyvals...)
}

func (w *temporalLoggerWrapper) Info(msg string, keyvals ...interface{}) {
	keyvals = w.appendFields(keyvals)
	w.Logger.Info(msg, keyvals...)
}

func (w *temporalLoggerWrapper) Warn(msg string, keyvals ...interface{}) {
	keyvals = w.appendFields(keyvals)
	w.Logger.Warn(msg, keyvals...)
}

func (w *temporalLoggerWrapper) Error(msg string, keyvals ...interface{}) {
	keyvals = w.appendFields(keyvals)
	w.Logger.Error(msg, keyvals...)
}

func (w *temporalLoggerWrapper) WithFields(fields map[string]interface{}) Logger {
	newFields := make(map[string]interface{})
	for k, v := range w.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	return &temporalLoggerWrapper{
		Logger: w.Logger,
		fields: newFields,
	}
}

func (w *temporalLoggerWrapper) appendFields(keyvals []interface{}) []interface{} {
	for k, v := range w.fields {
		keyvals = append(keyvals, k, v)
	}
	return keyvals
}

// getDefaultFields 获取默认字段
func (l *ZapLogger) getDefaultFields() []zap.Field {
	fields := make([]zap.Field, 0, len(l.fields))
	for k, v := range l.fields {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

// createZapLogger 创建 zap logger
func createZapLogger(level LogLevel, config LoggerConfig) *zap.Logger {
	var encoder zapcore.Encoder

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder,
		EncodeLevel:    customLevelEncoder(config.WithIcon),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	if config.ShowCaller {
		encoderConfig.EncodeCaller = customCallerEncoder(config.WithIcon)
	}

	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.ConsoleSeparator = " ║ "
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		mapToZapLevel(level),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(3))
	return logger.Named("temporal")
}

// ========================= 编码器 =========================

// customTimeEncoder 时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// customLevelEncoder 级别编码器
func customLevelEncoder(withIcon bool) zapcore.LevelEncoder {
	return func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		var levelStr string
		switch level {
		case zapcore.DebugLevel:
			if withIcon {
				levelStr = getColor(ColorCyanBold) + "🔍 [DEBUG]" + getColor(ColorReset)
			} else {
				levelStr = "[DEBUG]"
			}
		case zapcore.InfoLevel:
			if withIcon {
				levelStr = getColor(ColorGreenBold) + "ℹ️ [INFO]" + getColor(ColorReset)
			} else {
				levelStr = "[INFO] "
			}
		case zapcore.WarnLevel:
			if withIcon {
				levelStr = getColor(ColorYellowBold) + "⚠️ [WARN]" + getColor(ColorReset)
			} else {
				levelStr = "[WARN] "
			}
		case zapcore.ErrorLevel:
			if withIcon {
				levelStr = getColor(ColorRedBold) + "🚨 [ERROR]" + getColor(ColorReset)
			} else {
				levelStr = "[ERROR]"
			}
		default:
			levelStr = "[UNKNOWN]"
		}
		enc.AppendString(levelStr)
	}
}

// customCallerEncoder 调用者编码器
func customCallerEncoder(withIcon bool) zapcore.CallerEncoder {
	return func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		file := getRelativePath(caller.File)

		pc := caller.PC
		fn := runtime.FuncForPC(pc)
		funcName := "unknown"
		if fn != nil {
			funcName = fn.Name()
			if idx := strings.LastIndex(funcName, "."); idx != -1 {
				funcName = funcName[idx+1:]
			}
		}

		if withIcon {
			fileStr := getColor(ColorBlueBold) + fmt.Sprintf("📂 %s:%d", file, caller.Line) + getColor(ColorReset)
			funcStr := getColor(ColorMagentaBold) + funcName + getColor(ColorReset)
			enc.AppendString(fmt.Sprintf("%s:%s", fileStr, funcStr))
		} else {
			enc.AppendString(fmt.Sprintf("@ %s:%d:%s", file, caller.Line, funcName))
		}
	}
}

// ========================= 辅助函数 =========================

// isColorTerminal 检测终端是否支持颜色
func isColorTerminal() bool {
	term := os.Getenv("TERM")
	colorTerm := os.Getenv("COLORTERM")

	if colorTerm != "" {
		return true
	}

	colorTerms := []string{"xterm", "screen", "tmux", "rxvt", "ansi"}
	for _, ct := range colorTerms {
		if strings.Contains(term, ct) {
			return true
		}
	}
	return false
}

// getColor 根据配置返回颜色代码
func getColor(color string) string {
	if enableColor {
		return color
	}
	return ""
}

// findProjectRoot 查找项目根目录
func findProjectRoot() string {
	projectRootOnce.Do(func() {
		_, currentFile, _, ok := runtime.Caller(0)
		if !ok {
			if wd, err := os.Getwd(); err == nil {
				projectRootPath = wd
			} else {
				projectRootPath = "."
			}
			return
		}

		dir := filepath.Dir(currentFile)
		for {
			for _, marker := range []string{"go.mod", ".git", "Makefile", "README.md"} {
				if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
					projectRootPath = dir
					return
				}
			}

			parent := filepath.Dir(dir)
			if parent == dir {
				if wd, err := os.Getwd(); err == nil {
					projectRootPath = wd
				} else {
					projectRootPath = "."
				}
				return
			}
			dir = parent
		}
	})
	return projectRootPath
}

// getRelativePath 获取相对路径
func getRelativePath(fullPath string) string {
	projectRoot := findProjectRoot()

	if relPath, err := filepath.Rel(projectRoot, fullPath); err == nil {
		if !strings.HasPrefix(relPath, "..") {
			cleanPath := filepath.ToSlash(relPath)

			// 智能压缩长路径
			parts := strings.Split(cleanPath, "/")
			if len(parts) > 4 {
				cleanPath = parts[0] + "/.../" + strings.Join(parts[len(parts)-2:], "/")
			}
			return cleanPath
		}
	}

	// 兜底方案
	parts := strings.Split(filepath.ToSlash(fullPath), "/")
	if len(parts) >= 3 {
		return strings.Join(parts[len(parts)-3:], "/")
	}
	return filepath.Base(fullPath)
}

// isFrameworkCode 判断是否是框架代码
func isFrameworkCode(file string) bool {
	frameworkPaths := []string{
		"/go/pkg/mod/", "/usr/local/go/src/", "go.temporal.io",
		"go.uber.org/zap", "/runtime/", "/reflect/", "src/",
	}

	for _, framework := range frameworkPaths {
		if strings.Contains(file, framework) {
			return true
		}
	}

	return strings.HasSuffix(file, "logger.go")
}

// findUserCode 查找用户代码位置
func findUserCode(skip int) (string, int, string, bool) {
	for i := skip; i < skip+10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		if isFrameworkCode(file) {
			continue
		}

		fn := runtime.FuncForPC(pc)
		funcName := "unknown"
		if fn != nil {
			funcName = fn.Name()
			if idx := strings.LastIndex(funcName, "."); idx != -1 {
				funcName = funcName[idx+1:]
			}
		}

		return file, line, funcName, true
	}
	return "unknown", 0, "unknown", false
}

// getCallerField 获取调用者信息
func (l *ZapLogger) getCallerField() zap.Field {
	file, line, funcName, ok := findUserCode(3)
	if !ok {
		return zap.String("source", "unknown")
	}

	relativePath := getRelativePath(file)
	return zap.String("source", fmt.Sprintf("%s:%d:%s", relativePath, line, funcName))
}

// parseKeyvals 转换 keyvals 为 zap 字段
func (l *ZapLogger) parseKeyvals(keyvals ...interface{}) []zap.Field {
	if len(keyvals) == 0 {
		return nil
	}

	fields := make([]zap.Field, 0, len(keyvals)/2)

	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key := toString(keyvals[i])
			value := keyvals[i+1]
			fields = append(fields, zap.Any(key, value))
		} else {
			key := toString(keyvals[i])
			fields = append(fields, zap.String(key, "<missing>"))
		}
	}

	return fields
}

// toString 安全转换为字符串
func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return "unknown"
}

// mapToZapLevel 转换日志级别
func mapToZapLevel(level LogLevel) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

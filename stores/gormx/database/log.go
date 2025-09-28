package database

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

func NewLog() logger2.Interface {
	return &ormLog{}
}

type ormLog struct {
	LogLevel logger2.LogLevel
}

func (l *ormLog) LogMode(level logger2.LogLevel) logger2.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *ormLog) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger2.Info {
		logx.WithContext(ctx).Info("sql:", fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

func (l *ormLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger2.Warn {
		logx.WithContext(ctx).Error("sql:", fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

func (l *ormLog) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger2.Error {
		logx.WithContext(ctx).Error("sql:", fmt.Sprintf(msg, append([]interface{}{utils.FileWithLineNum()}, data...)...))
	}
}

func (l *ormLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger2.Silent {
		return
	}
	elapsed := time.Since(begin).Milliseconds()
	sql, rows := fc()
	logx.WithContext(ctx).WithCallerSkip(l.fileIndex()).Info("sql:", sql, " , time_ms:", elapsed, " , row:", rows)
}

func (l *ormLog) fileIndex() int {
	pcs := [13]uintptr{}
	// the third caller usually from gorm internal
	length := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:length])
	for i := 0; i < length; i++ {
		// second return value is "more", not "ok"
		frame, _ := frames.Next()
		if (!strings.Contains(frame.File, "/gorm.io/") ||
			strings.HasSuffix(frame.File, "_test.go")) && !strings.HasSuffix(frame.File, ".gen.go") && !strings.Contains(frame.File, "/dao/") {
			return i + 1
		}
	}

	return 1
}

package pkg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

type GormLogger struct {
	internal Logger
	level    logger.LogLevel
}

func NewGormLogger(l Logger, level logger.LogLevel) *GormLogger {
	return &GormLogger{internal: l, level: level}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		l.internal.Info(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Warn {
		l.internal.Warn(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Error {
		l.internal.Error(fmt.Sprintf(msg, data...))
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.level >= logger.Error && !errors.Is(err, logger.ErrRecordNotFound):
		sql, rows := fc()
		l.internal.Error(fmt.Sprintf(
			"[%.3fms] [rows:%d] %s | ERROR: %v",
			float64(elapsed.Nanoseconds())/1e6, rows, sql, err,
		))
	case elapsed > 200*time.Millisecond && l.level >= logger.Warn:
		sql, rows := fc()
		l.internal.Warn(fmt.Sprintf(
			"[%.3fms] [rows:%d] %s | SLOW QUERY",
			float64(elapsed.Nanoseconds())/1e6, rows, sql,
		))
	case l.level >= logger.Info:
		sql, rows := fc()
		l.internal.Info(fmt.Sprintf(
			"[%.3fms] [rows:%d] %s",
			float64(elapsed.Nanoseconds())/1e6, rows, sql,
		))
	}
}

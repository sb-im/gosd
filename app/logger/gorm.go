package logger

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	slowThreshold = 200 * time.Millisecond

	traceStr = "%s %s\t[%.3fms] [rows:%v] %s"
)

type Logger struct {
	log logger.Interface
}

func NewGorm() *Logger {
	return &Logger{
		log: logger.Default,
	}
}

func (l *Logger) LogMode(logger.LogLevel) logger.Interface { return l }

func (l *Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	WithContext(ctx).Info(append([]interface{}{msg}, args...))
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	WithContext(ctx).Warn(append([]interface{}{msg}, args...))
}

func (l *Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	WithContext(ctx).Error(append([]interface{}{msg}, args...))
}

func (t *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	l := WithContext(ctx)

	elapsed := time.Since(begin)
	switch {
	case err != nil:
		sql, rows := fc()
		if rows == -1 {
			l.Errorf(traceStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Errorf(traceStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > slowThreshold && slowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", slowThreshold)
		if rows == -1 {
			l.Warnf(traceStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Warnf(traceStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			l.Tracef(traceStr, utils.FileWithLineNum(), "", float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Tracef(traceStr, utils.FileWithLineNum(), "", float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

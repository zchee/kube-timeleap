// Copyright 2020 The kube-timeleap Authors.
// SPDX-License-Identifier: BSD-3-Clause

package logging

import (
	"time"

	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	zapf "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Logger represents a logr.Logger embedded logger.
type Logger struct {
	logr logr.Logger
}

// compile time check whether the Logger implements logr.Logger interface.
var _ logr.Logger = (*Logger)(nil)

// NewLogger returns the new zapr implemented logr.Logger.
func NewLogger(debug bool, zapfOpts ...zapf.Opts) *Logger {
	opts := make([]zapf.Opts, 0, len(zapfOpts)+4) // +3 is JSONEncoder, Level, StacktraceLevel and RawZapOpts

	lvl := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	stacktracelvl := zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	rawOpts := []zap.Option{zap.WithCaller(debug), zap.AddCallerSkip(1)}

	if debug {
		opts = make([]zapf.Opts, 0, len(zapfOpts)+5) // lazy optimize
		opts = append(opts, zapf.UseDevMode(true))

		lvl = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		stacktracelvl = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	}

	opts = append(opts,
		zapf.JSONEncoder(EncoderConfigs()...),
		zapf.Level(lvl),
		zapf.StacktraceLevel(stacktracelvl),
		zapf.RawZapOpts(rawOpts...),
	)

	// append zapfOpts variadic args to end of opts
	opts = append(opts, zapfOpts...)

	return &Logger{
		logr: zapf.New(opts...),
	}
}

// Enabled implements logr.Logger.Enabled.
func (l *Logger) Enabled() bool {
	return l.logr.Enabled()
}

// Info implements logr.Logger.Info.
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.logr.Info(msg, keysAndValues...)
}

// Error implements logr.Logger.Error.
func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logr.Error(err, msg, keysAndValues...)
}

// V implements logr.Logger.V.
func (l *Logger) V(level int) logr.Logger {
	return l.logr.V(level)
}

// WithValues implements logr.Logger.WithValues.
func (l *Logger) WithValues(keysAndValues ...interface{}) logr.Logger {
	return l.logr.WithValues(keysAndValues...)
}

// WithName implements logr.Logger.WithName.
func (l *Logger) WithName(name string) logr.Logger {
	return l.logr.WithName(name)
}

// EncoderConfigs returns the zapf.EncoderConfigOption.
func EncoderConfigs() []zapf.EncoderConfigOption {
	return []zapf.EncoderConfigOption{
		func(enc *zapcore.EncoderConfig) {
			enc.TimeKey = "ts"
			enc.LevelKey = "level"
			enc.NameKey = "logger"
			enc.CallerKey = "caller"
			enc.MessageKey = "msg"
			enc.StacktraceKey = "stacktrace"
			enc.LineEnding = zapcore.DefaultLineEnding
			enc.EncodeLevel = zapcore.LowercaseLevelEncoder
			enc.EncodeTime = zapcore.TimeEncoder(ISO8601TimeEncoder)
			enc.EncodeDuration = zapcore.StringDurationEncoder
			enc.EncodeCaller = zapcore.ShortCallerEncoder
		},
	}
}

// ISO8601TimeEncoder serializes a time.Time to an ISO8601-formatted string
// with millisecond precision.
// This is same as t.Format("2006-01-02T15:04:05.000Z").
//
// It optimized at byte slice level, faster than time.Format.
func ISO8601TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.UTC()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	msec := t.Nanosecond() / 1e6

	buf := make([]byte, 24)

	buf[0] = byte((year/1000)%10) + '0'
	buf[1] = byte((year/100)%10) + '0'
	buf[2] = byte((year/10)%10) + '0'
	buf[3] = byte(year%10) + '0'
	buf[4] = '-'
	buf[5] = byte((month)/10) + '0'
	buf[6] = byte((month)%10) + '0'
	buf[7] = '-'
	buf[8] = byte((day)/10) + '0'
	buf[9] = byte((day)%10) + '0'
	buf[10] = 'T'
	buf[11] = byte((hour)/10) + '0'
	buf[12] = byte((hour)%10) + '0'
	buf[13] = ':'
	buf[14] = byte((minute)/10) + '0'
	buf[15] = byte((minute)%10) + '0'
	buf[16] = ':'
	buf[17] = byte((second)/10) + '0'
	buf[18] = byte((second)%10) + '0'
	buf[19] = '.'
	buf[20] = byte((msec/100)%10) + '0'
	buf[21] = byte((msec/10)%10) + '0'
	buf[22] = byte((msec)%10) + '0'
	buf[23] = 'Z'

	enc.AppendString(string(buf))
}

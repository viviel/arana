//
// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package log

import (
	"bytes"
	"errors"
	"fmt"
)

import (
	"github.com/natefinch/lumberjack"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level represents the level of logging.
type LogLevel int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = LogLevel(zapcore.DebugLevel)
	// InfoLevel is the default logging priority.
	InfoLevel = LogLevel(zapcore.InfoLevel)
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = LogLevel(zapcore.WarnLevel)
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = LogLevel(zapcore.ErrorLevel)
	// PanicLevel logs a message, then panics.
	PanicLevel = LogLevel(zapcore.PanicLevel)
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = LogLevel(zapcore.FatalLevel)
)

func (l *LogLevel) UnmarshalText(text []byte) error {
	if l == nil {
		return errors.New("can't unmarshal a nil *Level")
	}
	if !l.unmarshalText(text) && !l.unmarshalText(bytes.ToLower(text)) {
		return fmt.Errorf("unrecognized level: %q", text)
	}
	return nil
}

func (l *LogLevel) unmarshalText(text []byte) bool {
	switch string(text) {
	case "debug", "DEBUG":
		*l = DebugLevel
	case "info", "INFO", "": // make the zero value useful
		*l = InfoLevel
	case "warn", "WARN":
		*l = WarnLevel
	case "error", "ERROR":
		*l = ErrorLevel
	case "panic", "PANIC":
		*l = PanicLevel
	case "fatal", "FATAL":
		*l = FatalLevel
	default:
		return false
	}
	return true
}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

var (
	log       Logger
	zapLogger *zap.Logger

	zapLoggerConfig        = zap.NewDevelopmentConfig()
	zapLoggerEncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
)

func init() {
	zapLoggerConfig.EncoderConfig = zapLoggerEncoderConfig
	zapLogger, _ = zapLoggerConfig.Build(zap.AddCallerSkip(1))
	log = zapLogger.Sugar()
}

func Init(logPath string, level LogLevel) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	syncer := zapcore.AddSync(lumberJackLogger)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, syncer, zap.NewAtomicLevelAt(zapcore.Level(level)))
	zapLogger = zap.New(core, zap.AddCaller())

	log = zapLogger.Sugar()
}

// SetLogger: customize yourself logger.
func SetLogger(logger Logger) {
	log = logger
}

// GetLogger get logger
func GetLogger() Logger {
	return log
}

// Debug ...
func Debug(v ...interface{}) {
	log.Debug(v...)
}

// Debugf ...
func Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

// Info ...
func Info(v ...interface{}) {
	log.Info(v...)
}

// Infof ...
func Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

// Warn ...
func Warn(v ...interface{}) {
	log.Warn(v...)
}

// Warnf ...
func Warnf(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

// Error ...
func Error(v ...interface{}) {
	log.Error(v...)
}

// Errorf ...
func Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

// Panic ...
func Panic(v ...interface{}) {
	log.Panic(v...)
}

// Panicf ...
func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

// Fatal ...
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

// Fatalf ...
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

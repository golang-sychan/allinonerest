package server

import "go.uber.org/zap"

type Logger interface {
	LoggerSetter
	LoggerGetter
}

type LoggerSetter interface {
	SetLogger(logger *zap.Logger)
}

type LoggerGetter interface {
	L() *zap.Logger
	Log() *zap.Logger
}

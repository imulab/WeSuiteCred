package internal

import (
	"github.com/eclipse/paho.golang/paho"
	"github.com/rs/zerolog"
)

func NewPahoZeroLogger(logger *zerolog.Logger) paho.Logger {
	return &pahoZeroLogger{logger: logger}
}

type pahoZeroLogger struct {
	logger *zerolog.Logger
}

func (p *pahoZeroLogger) Println(v ...interface{}) {
	p.logger.Print(v...)
}

func (p *pahoZeroLogger) Printf(format string, v ...interface{}) {
	p.logger.Printf(format, v...)
}

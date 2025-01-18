package logger

import (
	"fmt"
	"log"
)

type prefixLogger struct {
	prefix string
}

func (l *prefixLogger) Println(args ...any) {
	log.Println([]any{l.prefix, args}...)
}

func NewLoggerWithPrefix(prefix string) Logger {
	return &prefixLogger{prefix: fmt.Sprint("[", prefix, "]")}
}

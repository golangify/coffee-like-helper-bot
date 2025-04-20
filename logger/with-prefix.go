package logger

import (
	"fmt"
	"log"
)

type prefixLogger struct {
	prefix string
}

func NewLoggerWithPrefix(prefix string) Logger {
	return &prefixLogger{prefix: fmt.Sprint("[", prefix, "]")}
}

func (l *prefixLogger) Println(v ...any) {
	log.Println(append([]any{l.prefix}, v...)...)
}

func (l *prefixLogger) Printf(format string, v ...any) {
	log.Printf("%s "+format, append([]any{l.prefix}, v...)...)
}

func (l *prefixLogger) Fatal(err error) {
	log.Fatalf("%s %+v", l.prefix, err)
}

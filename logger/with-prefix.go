package logger

import (
	"fmt"
	"log"
	"os"
)

type prefixLogger struct {
	prefix string
}

func NewLoggerWithPrefix(prefix string) Logger {
	return &prefixLogger{prefix: fmt.Sprint("[", prefix, "]")}
}

func (l *prefixLogger) Println(args ...any) {
	log.Println([]any{l.prefix, args}...)
}

func (l *prefixLogger) Fatal(err error) {
	log.Fatalf("%+v", err)
	os.Exit(1)
}

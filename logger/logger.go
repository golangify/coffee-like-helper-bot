package logger

type Logger interface {
	Println(v ...any)
	Printf(format string, v ...any)
	Fatal(err error)
}

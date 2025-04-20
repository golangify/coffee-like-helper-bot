package logger

type Logger interface {
	Println(...any)
	Fatal(err error)
}

package db

import "os"

type Logger struct{}

func (l *Logger) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

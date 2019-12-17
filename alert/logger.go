package alert

import "os"

type DebugLogger struct{}

func (l *DebugLogger) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

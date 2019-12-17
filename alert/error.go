package alert

type FlushTimeoutError struct{}

func (e *FlushTimeoutError) Error() string {
	return "cannot flush within timeout"
}

type AlertError interface {
	error
	ClientKey() ClientKey
}

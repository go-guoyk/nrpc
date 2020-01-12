package nrpc

import "fmt"

type Error struct {
	Status  string
	Message string
	Tried   int
}

func (e *Error) Error() string {
	if e.Tried > 0 {
		return fmt.Sprintf("%s (tried %d times): %s", e.Status, e.Tried, e.Message)
	} else {
		return fmt.Sprintf("%s: %s", e.Status, e.Message)
	}
}

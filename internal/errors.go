package internal

import "fmt"

type CustomError struct {
	OriginalError error
	Message       string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.OriginalError)
}

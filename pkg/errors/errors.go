package errors

import "fmt"

func Wrap(err error, message string) error {
	return fmt.Errorf("error: %s, err: %v", message, err)
}

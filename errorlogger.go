/*
Logger interface for sendcloud
*/
package sendcloud

import (
    "fmt"
)

type ErrorLogger interface {
    ErrorLog(source string, code int, msg string) error
}

type FmtErrorLogger struct {
}

func (l FmtErrorLogger) ErrorLog(source string, code int, msg string) error {
    return fmt.Errorf("%s error: code=%d, msg=%s", source, code, msg)
}

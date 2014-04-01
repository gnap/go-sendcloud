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
    if code != 200 {
        return fmt.Errorf("%s: code=%d, msg=%s", source, code, msg)
    }
    fmt.Printf("%s: code=%d, msg=%s", source, code, msg)
    return nil
}

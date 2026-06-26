package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

var (
	w    io.Writer
	once sync.Once
)

func SetLogFile(path string) error {
	var e error
	once.Do(func() {
		if path != "" {
			f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
			if err != nil {
				e = err
				return
			}

			w = f
		}
	})

	return e
}

func callerName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	return fn.Name()
}

func Printf(format string, a ...any) {
	str := fmt.Sprintf(format, a...)
	if w != nil {
		fmt.Fprint(w, str)
	}

	fmt.Print(str)
}

func logf(ltype, host, format string, a ...any) {
	if host != "main.main" {
		format = fmt.Sprintf("arr-go: %s: [%s] %s\n", host, ltype, format)
	} else {
		format = fmt.Sprintf("arr-go: [%s] %s\n", ltype, format)
	}

	Printf(format, a...)
}
func sLogf(ltype, host, format string, a ...any) string {
	var fm string
	if host != "main.main" {
		fm = host + ":"
	}

	if ltype != "" {
		fm += fmt.Sprintf(" [%s]", ltype)
	}

	return fm + fmt.Sprintf(format, a...)
}

func Infof(format string, a ...any) {
	caller := callerName()
	logf("info", caller, format, a...)
}

func AsError(format string, a ...any) error {
	caller := callerName()
	return errors.New(sLogf("error", caller, format, a...))
}

func Errorf(format string, a ...any) {
	caller := callerName()
	logf("error", caller, format, a...)
}

func Warnf(format string, a ...any) {
	caller := callerName()
	logf("warn", caller, format, a...)
}

func Fatalf(format string, a ...any) {
	caller := callerName()
	logf("fatal", caller, format, a...)
	os.Exit(1)
}

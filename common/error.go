package common

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type err struct {
	original      error
	wrapped       error
	stacktrace    string
	stack         []uintptr
	logCtx        string
	keyerr        error
	isNotify      bool
	isSuccessResp bool
}

func (e *err) Error() string {
	var original string
	var wrapped string

	if e.wrapped == nil {
		return e.original.Error() + e.stacktrace
	}

	wrapped = e.wrapped.Error() + e.stacktrace

	if e.original != nil {
		if _, ok := e.original.(*err); !ok {
			original = "root cause : " + e.original.Error()
		} else {
			original = e.original.Error()
		}
	}

	return wrapped + ": " + original
}

// WrapWithErr wrap new error into an existing error
func WrapWithErr(original error, wrapped error) *err {
	pc, file, no, _ := runtime.Caller(1)

	// Current working directory
	dir, _ := os.Getwd()
	splitDir := strings.Split(dir, "/")
	rootDir := splitDir[len(splitDir)-1]

	logCtx := generateLogCtx(pc, file, rootDir)

	return &err{
		original:   original,
		wrapped:    wrapped,
		keyerr:     GetErrKey(wrapped),
		stacktrace: " -- At : " + fmt.Sprintf("%s:%d", file, no),
		stack:      []uintptr{pc},
		logCtx:     logCtx,
	}
}

// get error as key to compare what the output response will be
func GetErrKey(err_ error) error {
	if val, ok := err_.(*err); ok {
		return val.keyerr
	}

	return err_
}

// generateLogCtx will generate the log context
func generateLogCtx(pc uintptr, file, rootDir string) string {
	var filePath, funcName string
	filePath = strings.Split(file, fmt.Sprintf("/%s/", rootDir))[1]
	filePath = strings.Split(filePath, ".")[0]
	filePath = strings.ReplaceAll(filePath, "/", ".")
	funcName = runtime.FuncForPC(pc).Name()
	i := strings.LastIndex(funcName, ".")
	funcName = funcName[i+1:]

	return fmt.Sprintf("%s.%s", filePath, funcName)
}

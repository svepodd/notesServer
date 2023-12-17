package pkg

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type EWrapper struct {
	functionName string
	comment      string
	err          error
	logFile      *os.File
}

func NewEWrapper(f string) *EWrapper {
	return &EWrapper{f, "", nil, nil}
}

func NewEWrapperWithFile(f string) (*EWrapper, error) {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // May be any file
	if err != nil {
		return nil, err
	}
	return &EWrapper{f, "", nil, file}, nil
}

func (e *EWrapper) Wrap(err error, comment string) *EWrapper {
	if err != nil {
		e.err = err
		e.comment = comment
	}
	return e
}

func (e *EWrapper) Error() error {
	if e.err == nil {
		return nil
	}
	return errors.New(e.comment + "    __IN__    " + e.functionName + ":\n" + e.err.Error() + "\n")
}

func (e *EWrapper) WrapError(err error, comment string) error {
	if err != nil {
		return e.Wrap(err, comment).Error()
	}
	return nil
}

func (e *EWrapper) LogError(err error, comment string) {
	if err != nil {
		e.comment = comment
		log.Println("\033[31m", e.Wrap(err, comment).Error(), "\033[0m")
		fmt.Fprintln(e.logFile, e.Wrap(err, comment).Error())

	}
}

func (e *EWrapper) Close() error {
	return e.logFile.Close()
}

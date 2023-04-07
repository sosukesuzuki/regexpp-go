package parser

import "fmt"

type ParserError struct {
	msg string
	err error
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("Error from parser: %s (%s)", e.msg, e.err.Error())
}

func (e *ParserError) Unwrap() error {
	return e.err
}

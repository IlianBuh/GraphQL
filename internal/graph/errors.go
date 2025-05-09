package graph

import "errors"

const (
	InvalidInput = "invalid input"
	Internal     = "internal"
)

func sendErr(mess string, err error) error {
	return errors.New(mess + ": " + err.Error())
}

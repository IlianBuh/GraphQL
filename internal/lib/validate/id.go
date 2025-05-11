package validate

import "errors"

type number interface {
	int | uint |
		int8 | uint8 |
		int16 | uint16 |
		int32 | uint32 |
		int64 | uint64
}

func Id[T number](id T) error {
	if id < 0 {
		return errors.New("id must be positive integer")
	}

	return nil
}

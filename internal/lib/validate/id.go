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
func Ids[T number](ids []T) error {
	var err error
	for _, id := range ids {

		if err = Id(id); err != nil {
			return err
		}

	}

	return nil
}

package validate

import (
	"errors"
	"strings"
)

func Login(login string) error {
	if len(strings.Trim(login, " ")) == 0 {
		return errors.New("login required")
	}

	return nil
}

package validate

import (
	"errors"
)

func SignUp(login, email, password string) error {
	if isEmpty(login) {
		return errors.New("login is required")
	}

	if isEmpty(email) {
		return errors.New("email is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 symbols")
	}

	return nil
}

package validate

import "errors"

func LogIn(login, password string) error {
	if isEmpty(login) {
		return errors.New("login is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 symbols")
	}

	return nil
}

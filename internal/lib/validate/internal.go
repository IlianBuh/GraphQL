package validate

import (
	"strings"
)

func isEmpty(str string) bool {
	return len(strings.Trim(str, " ")) == 0
}

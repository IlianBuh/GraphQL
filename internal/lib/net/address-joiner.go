package net

import (
	"fmt"
)

func Join(host string, port int) string {
	return host + ":" + fmt.Sprintf("%d", port)
}

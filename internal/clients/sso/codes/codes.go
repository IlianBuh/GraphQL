package codes

const (
	Ok              = 0
	InvalidArgument = 1
	Internal        = 2
	Unknown         = -1
)

var codeMessages map[int]string = map[int]string{
	Ok:              "ok",
	InvalidArgument: "invalid argument",
	Internal:        "internal",
	Unknown:         "unknown",
}

func Text(code int) string {
	message, ok := codeMessages[code]
	if !ok {
		return "unknown message code"
	}

	return message
}

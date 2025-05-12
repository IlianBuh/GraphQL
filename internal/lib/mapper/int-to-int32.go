package mapper

type number interface {
	int | uint |
		int8 | uint8 |
		int16 | uint16 |
		int32 | uint32 |
		int64 | uint64
}

func NumsTToNumsE[T, E number](vals []T) []E {
	res := make([]E, len(vals))

	for i, val := range vals {
		res[i] = E(val)
	}

	return res
}

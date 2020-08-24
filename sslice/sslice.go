package sslice

func AppendStr(slice []string, elems ...string) []string {
	slice = append(slice, elems...)
	return slice
}

func RemoveStr(slice []string, elems ...string) []string {
	if len(slice) == 0 || len(elems) == 0 {
		return slice
	}

	for i := 0; i < len(slice); i++ {
		v := slice[i]
		for j := 0; j < len(elems); j++ {
			if v == elems[j] {
				slice = append(slice[:i], slice[i+1:]...)
				i--
				// elems = append(elems[:j], elems[j+1:]...)
				// j--
			}
		}
	}

	return slice
}

// AppendStrToNew create a new slice, then append string
func AppendStrToNew(slice []string, elems ...string) []string {
	r := make([]string, len(slice))
	copy(r, slice)

	return AppendStr(r, elems...)
}

// RemoveStrToNew create a new slice, then remove string
func RemoveStrToNew(slice []string, elems ...string) []string {
	r := make([]string, len(slice))
	copy(r, slice)
	return RemoveStr(r, elems...)
}

package sslice

func AppendString(slice []string, elems ...string) []string {
	slice = append(slice, elems...)
	return slice
}

func RemoveString(slice []string, elems ...string) []string {
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

func AppendStringToNew(slice []string, elems ...string) []string {
	// 制作拷贝，防止复用
	r := make([]string, len(slice))
	copy(r, r)

	return AppendString(r, elems...)
}

func RemoveStringToNew(slice []string, elems ...string) []string {
	// 制作拷贝，防止复用
	r := make([]string, len(slice))
	copy(r, r)
	return RemoveString(r, elems...)
}

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
		for _, e := range elems {
			if slice[i] == e {
				slice = append(slice[:i], slice[i+1:]...)
				i-- // maintain the correct index
			}
		}
	}

	return slice
}

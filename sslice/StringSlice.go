package sslice

type StringSlice []string

func (x *StringSlice) Add(elem string) *StringSlice {
	*x = append(*x, elem)
	return x
}

func (x *StringSlice) Remove(elem string) *StringSlice {
	if len(*x) == 0 {
		return x
	}
	for i, v := range *x {
		if v == elem {
			*x = append((*x)[:i], (*x)[i+1:]...)
			return x.Remove(elem)
		}
	}
	return x
}

func (x *StringSlice) Has(value string) bool {
	for _, v := range *x {
		if v == value {
			return true
		}
	}
	return false
}

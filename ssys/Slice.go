package ssys

type StringSlice []string

func (x *StringSlice) Has(value string) bool {
	for _, v := range *x {
		if v == value {
			return true
		}
	}
	return false
}

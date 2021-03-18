package sconv

import "strconv"

func ToString(obj interface{}) string {
	switch v := obj.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 2, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 10)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}

func ToInt32(obj interface{}) int32 {
	switch v := obj.(type) {
	case int32:
		return v
	case int64:
		return int32(v)
	case int:
		return int32(v)
	case string:
		r, _ := strconv.ParseInt(v, 10, 32)
		return int32(r)
	default:
		return 0
	}
}

func ToInt64(obj interface{}) int64 {
	switch v := obj.(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case int:
		return int64(v)
	case string:
		r, _ := strconv.ParseInt(v, 10, 64)
		return r
	default:
		return 0
	}
}

func ToFloat32(obj interface{}) float32 {
	switch v := obj.(type) {
	case float32:
		return v
	case float64:
		return float32(v)
	case string:
		r, _ := strconv.ParseFloat(v, 10)
		return float32(r)
	default:
		return 0
	}
}

func ToFloat64(obj interface{}) float64 {
	switch v := obj.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case string:
		r, _ := strconv.ParseFloat(v, 10)
		return r
	default:
		return 0
	}
}

package sconv

import (
	"strconv"
	"strings"
)

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
		return strconv.FormatFloat(float64(v), 'f', 2, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', 2, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ""
	}
}

func ToInt(obj interface{}) int {
	switch v := obj.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case int32:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		v = strings.ReplaceAll(v, ",", "") // 移除逗号
		r, _ := strconv.ParseInt(v, 10, 32)
		return int(r)
	default:
		return 0
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
	case float32:
		return int32(v)
	case float64:
		return int32(v)
	case string:
		v = strings.ReplaceAll(v, ",", "") // 移除逗号
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
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		v = strings.ReplaceAll(v, ",", "") // 移除逗号
		r, _ := strconv.ParseInt(v, 10, 64)
		return r
	default:
		return 0
	}
}

func ToFloat32(obj interface{}) float32 {
	switch v := obj.(type) {
	case int:
		return float32(v)
	case int32:
		return float32(v)
	case int64:
		return float32(v)
	case float32:
		return v
	case float64:
		return float32(v)
	case string:
		v = strings.ReplaceAll(v, ",", "") // 移除逗号
		r, _ := strconv.ParseFloat(v, 32)
		return float32(r)
	default:
		return 0
	}
}

func ToFloat64(obj interface{}) float64 {
	switch v := obj.(type) {
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case float64:
		return v
	case float32:
		return float64(v)
	case string:
		v = strings.ReplaceAll(v, ",", "") // 移除逗号
		r, _ := strconv.ParseFloat(v, 64)
		return r
	default:
		return 0
	}
}

func ToBool(obj interface{}) bool {
	switch v := obj.(type) {
	case bool:
		return v
	case string:
		r, _ := strconv.ParseBool(v)
		return r
	default:
		return false
	}
}

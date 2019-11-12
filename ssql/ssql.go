package ssql

import (
	"database/sql"
	"errors"
	"reflect"
)

// MapRowsToSlice map db data to slice
func MapRowsToSlice(rows *sql.Rows, itemType *reflect.Type, cap int) (interface{}, error) {
	// 根据Item类型创建Slice
	sliceValue := reflect.MakeSlice(reflect.SliceOf(*itemType), 0, cap)
	// 获取数据库列名
	columns, _ := rows.Columns()

	var err error
	//轮询数据行
	for rows.Next() {
		// 构造Item实例，获得实例指针
		itemPtr := reflect.New(*itemType)
		// 获取Item各个字段的指针slice
		itemFieldPtrs, err := StructForScan(&columns, itemPtr.Interface())
		if err != nil {
			return nil, err
		}

		// 数据行扫描，将匹配的列数据存入Item实例各字段指针所指向的地址
		err = rows.Scan(itemFieldPtrs...)
		if err != nil {
			return nil, err
		}
		// 添加到结果slice里
		itemValue := reflect.Indirect(itemPtr)
		sliceValue = reflect.Append(sliceValue, itemValue)
	}

	return sliceValue.Interface(), err
}

// StructForScan _
func StructForScan(columns *[]string, itemPtr interface{}) (slicePtrs []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = errors.New("unknown panic")
			}

			slicePtrs = nil
		}
	}()

	// 获取itemPtr的类型
	inputParaType := reflect.TypeOf(itemPtr)
	// 确保itemPtr是指针
	if inputParaType.Kind() != reflect.Ptr {
		panic("sql: itemPtr must be a struct pointer")
		// return nil, errors.New("sql: itemPtr must be a struct pointer")
	}

	// 获取input指针指向的类型
	structType := inputParaType.Elem()
	// 确保input指向的值是struct
	if structType.Kind() != reflect.Struct {
		return nil, errors.New("sql: itemPtr must be a struct pointer")
	}

	// 获取struct的指针
	structPtr := reflect.ValueOf(itemPtr)
	// 根据指针获取struct对象
	structValue := reflect.Indirect(structPtr)

	// 创建返回slice的指针
	slicePtrs = make([]interface{}, 0, len(*columns))

	for _, column := range *columns {
		// 根据数据库列名查找struct是否包含对应字段
		_, found := structType.FieldByName(column)
		// 如果struct有此字段
		if found {
			// 获取字段
			valueField := structValue.FieldByName(column)
			// 判断字段是否可以取指针
			if valueField.CanAddr() {
				// 取指针
				fieldPtr := valueField.Addr().Interface()
				// 存入返回slice
				slicePtrs = append(slicePtrs, fieldPtr)
			} else {
				// 如果不能取指针，传入一个随机地址指针
				slicePtrs = append(slicePtrs, new(interface{}))
				// slicePtrs = append(slicePtrs, valueField)
			}
		} else {
			// 数据库的字段跟DTO的不匹配，传入一个随机地址指针
			slicePtrs = append(slicePtrs, new(interface{}))
		}
	}

	return slicePtrs, nil
}

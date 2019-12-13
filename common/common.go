package common

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

func DataToStruct(data map[string]string, obj interface{}) {
	//通过反射拿到obj的Value,通过Elem判断是否为指针类型，如果不是，则panic
	objVal := reflect.ValueOf(obj).Elem()
	objType := objVal.Type()
	var key string
	//循环这个objVal的参数
	for i := 0; i < objVal.NumField(); i++ {
		//获取到第一个值的键
		key = objType.Field(i).Tag.Get("sql")
		value := data[key]
		//获取这个值的反射对象
		val := reflect.ValueOf(value)
		fieldType := objVal.Field(i).Type()
		//判断当前这个键的类型，如果是string，则直接取，如果是其他的，则需要转换一层
		if fieldType != val.Type() {
			val, _ = TypeConversion(value, fieldType.Name())
		}
		//将这个键值放入到obj的结构体中
		objVal.Field(i).Set(val)
	}

}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}

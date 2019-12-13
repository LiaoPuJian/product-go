package main

import (
	"fmt"
	"product-go/common"
	"product-go/models"
	"reflect"
)

func main() {
	p := make(map[string]string)
	p["id"] = "1"
	p["product_name"] = "lpj"
	p["product_num"] = "1"
	p["product_image"] = "https://www.baidu.com/image"
	p["product_url"] = "https://www.baidu.com"

	model := &models.Product{}

	DataToStruct(p, model)

	fmt.Printf("type:%T, value:%v", model, model)
}

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
			val, _ = common.TypeConversion(value, fieldType.Name())
		}
		//将这个键值放入到obj的结构体中
		objVal.Field(i).Set(val)
	}

}

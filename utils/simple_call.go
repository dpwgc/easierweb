package utils

import (
	"github.com/dpwgc/easierweb"
	"reflect"
)

func SimpleCallJSON(ctx *easierweb.Context, simpleHandle any) (any, error) {

	// 反射获取函数类型
	funcType := reflect.TypeOf(simpleHandle)

	// 创建参数值的切片
	paramValues := make([]reflect.Value, 2)
	paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
	paramValues[1] = reflect.New(funcType.In(1)).Elem()

	if len(ctx.Form) > 0 {
		_ = ctx.BindForm(paramValues[1].Addr().Interface())
	}
	if len(ctx.Body) > 0 {
		_ = ctx.BindJSON(paramValues[1].Addr().Interface())
	}
	if len(ctx.Query) > 0 {
		_ = ctx.BindQuery(paramValues[1].Addr().Interface())
	}

	// 调用函数
	returnValues := reflect.ValueOf(simpleHandle).Call(paramValues)

	// 处理返回值
	var resultValue any
	if returnValues[0].IsValid() && returnValues[0].Kind() == reflect.Ptr && returnValues[0].Elem().IsValid() {
		resultValue = returnValues[0].Elem().Interface()
	} else if returnValues[0].IsValid() && returnValues[0].Kind() == reflect.Slice {
		resultValue = returnValues[0].Interface()
	}
	errValue, _ := returnValues[1].Interface().(error)
	return resultValue, errValue
}

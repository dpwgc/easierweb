package plugins

import (
	"github.com/dpwgc/easierweb"
	"reflect"
)

func JSONRequestHandle(ctx *easierweb.Context, easyHandle any) (any, error) {

	// 反射获取函数类型
	funcType := reflect.TypeOf(easyHandle)

	// 创建参数值的切片
	var paramValues []reflect.Value

	// 如果没有第二个参数，就不传了
	if funcType.NumIn() == 1 {
		paramValues = make([]reflect.Value, 1)
		paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
	} else {
		paramValues = make([]reflect.Value, 2)
		paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
		paramValues[1] = reflect.New(funcType.In(1)).Elem()
		if len(ctx.Form) > 0 {
			err := ctx.BindForm(paramValues[1].Addr().Interface())
			if err != nil {
				return nil, err
			}
		}
		if len(ctx.Body) > 0 {
			err := ctx.BindJSON(paramValues[1].Addr().Interface())
			if err != nil {
				return nil, err
			}
		}
		if len(ctx.Query) > 0 {
			err := ctx.BindQuery(paramValues[1].Addr().Interface())
			if err != nil {
				return nil, err
			}
		}
	}

	// 调用函数
	returnValues := reflect.ValueOf(easyHandle).Call(paramValues)

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

func YAMLRequestHandle(ctx *easierweb.Context, easyHandle any) (any, error) {

	// 反射获取函数类型
	funcType := reflect.TypeOf(easyHandle)

	// 创建参数值的切片
	var paramValues []reflect.Value

	// 如果没有第二个参数，就不传了
	if funcType.NumIn() == 1 {
		paramValues = make([]reflect.Value, 1)
		paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
	} else {
		paramValues = make([]reflect.Value, 2)
		paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
		paramValues[1] = reflect.New(funcType.In(1)).Elem()
		if len(ctx.Form) > 0 {
			err := ctx.BindForm(paramValues[1].Addr().Interface())
			if err != nil {
				return nil, err
			}
		}
		if len(ctx.Body) > 0 {
			err := ctx.BindYAML(paramValues[1].Addr().Interface())
			if err != nil {
				return nil, err
			}
		}
		if len(ctx.Query) > 0 {
			err := ctx.BindQuery(paramValues[1].Addr().Interface())
			if err != nil {
				return nil, err
			}
		}
	}

	// 调用函数
	returnValues := reflect.ValueOf(easyHandle).Call(paramValues)

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

func XMLRequestHandle(ctx *easierweb.Context, easyHandle any) (any, error) {

	// 反射获取函数类型
	funcType := reflect.TypeOf(easyHandle)

	// 创建参数值的切片
	var paramValues []reflect.Value

	// 如果没有第二个参数，就不传了
	if funcType.NumIn() == 1 {
		paramValues = make([]reflect.Value, 1)
		paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
	} else {
		paramValues = make([]reflect.Value, 2)
		paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
		paramValues[1] = reflect.New(funcType.In(1)).Elem()
		if len(ctx.Form) > 0 {
			err := ctx.BindForm(paramValues[1].Addr().Interface())
			if err != nil {
				return nil, err
			}
		}
		if len(ctx.Body) > 0 {
			err := ctx.BindXML(paramValues[1].Addr().Interface())
			if err != nil {
				return nil, err
			}
		}
		if len(ctx.Query) > 0 {
			err := ctx.BindQuery(paramValues[1].Addr().Interface())
			if err != nil {
				return nil, err
			}
		}
	}

	// 调用函数
	returnValues := reflect.ValueOf(easyHandle).Call(paramValues)

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

func NoActionRequestHandle(ctx *easierweb.Context, easyHandle any) (any, error) {

	// 反射获取函数类型
	funcType := reflect.TypeOf(easyHandle)

	// 创建参数值的切片
	var paramValues []reflect.Value

	// 如果没有第二个参数，就不传了
	if funcType.NumIn() == 1 {
		paramValues = make([]reflect.Value, 1)
		paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
	} else {
		paramValues = make([]reflect.Value, 2)
		paramValues[0] = reflect.ValueOf(ctx).Elem().Addr()
		paramValues[1] = reflect.New(funcType.In(1)).Elem()
	}

	// 调用函数
	returnValues := reflect.ValueOf(easyHandle).Call(paramValues)

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

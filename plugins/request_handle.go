package plugins

import (
	"github.com/dpwgc/easierweb"
	"reflect"
)

func JSONRequestHandle(ctx *easierweb.Context, paramValues []reflect.Value) error {
	if len(ctx.Body) > 0 {
		err := ctx.BindJSON(paramValues[1].Addr().Interface())
		if err != nil {
			return err
		}
	}
	if len(ctx.Query) > 0 {
		err := ctx.BindQuery(paramValues[1].Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func YAMLRequestHandle(ctx *easierweb.Context, paramValues []reflect.Value) error {
	if len(ctx.Body) > 0 {
		err := ctx.BindYAML(paramValues[1].Addr().Interface())
		if err != nil {
			return err
		}
	}
	if len(ctx.Query) > 0 {
		err := ctx.BindQuery(paramValues[1].Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func XMLRequestHandle(ctx *easierweb.Context, paramValues []reflect.Value) error {
	if len(ctx.Body) > 0 {
		err := ctx.BindXML(paramValues[1].Addr().Interface())
		if err != nil {
			return err
		}
	}
	if len(ctx.Query) > 0 {
		err := ctx.BindQuery(paramValues[1].Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

func FormRequestHandle(ctx *easierweb.Context, paramValues []reflect.Value) error {
	if len(ctx.Form) > 0 {
		err := ctx.BindForm(paramValues[1].Addr().Interface())
		if err != nil {
			return err
		}
	}
	if len(ctx.Query) > 0 {
		err := ctx.BindQuery(paramValues[1].Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

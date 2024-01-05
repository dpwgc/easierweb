package plugins

import (
	"github.com/dpwgc/easierweb"
)

func JSONRequestHandle() easierweb.RequestHandle {
	return func(ctx *easierweb.Context, reqObj any) error {
		if len(ctx.Form) > 0 {
			err := ctx.BindForm(reqObj)
			if err != nil {
				return err
			}
		} else if len(ctx.Body) > 0 {
			err := ctx.BindJSON(reqObj)
			if err != nil {
				return err
			}
		}
		if len(ctx.Query) > 0 {
			err := ctx.BindQuery(reqObj)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func YAMLRequestHandle() easierweb.RequestHandle {
	return func(ctx *easierweb.Context, reqObj any) error {
		if len(ctx.Form) > 0 {
			err := ctx.BindForm(reqObj)
			if err != nil {
				return err
			}
		} else if len(ctx.Body) > 0 {
			err := ctx.BindYAML(reqObj)
			if err != nil {
				return err
			}
		}
		if len(ctx.Query) > 0 {
			err := ctx.BindQuery(reqObj)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func XMLRequestHandle() easierweb.RequestHandle {
	return func(ctx *easierweb.Context, reqObj any) error {
		if len(ctx.Form) > 0 {
			err := ctx.BindForm(reqObj)
			if err != nil {
				return err
			}
		} else if len(ctx.Body) > 0 {
			err := ctx.BindXML(reqObj)
			if err != nil {
				return err
			}
		}
		if len(ctx.Query) > 0 {
			err := ctx.BindQuery(reqObj)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

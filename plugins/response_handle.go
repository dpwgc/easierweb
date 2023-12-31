package plugins

import (
	"github.com/dpwgc/easierweb"
	"net/http"
)

func JSONResponseHandle() easierweb.ResponseHandle {
	return func(ctx *easierweb.Context, result any, err error) {
		if err != nil {
			if result != nil {
				ctx.WriteJSON(http.StatusBadRequest, result)
				return
			}
			panic(err)
		}
		if result == nil {
			ctx.NoContent(http.StatusNoContent)
			return
		}
		ctx.WriteJSON(http.StatusOK, result)
	}
}

func YAMLResponseHandle() easierweb.ResponseHandle {
	return func(ctx *easierweb.Context, result any, err error) {
		if err != nil {
			if result != nil {
				ctx.WriteYAML(http.StatusBadRequest, result)
				return
			}
			panic(err)
		}
		if result == nil {
			ctx.NoContent(http.StatusNoContent)
			return
		}
		ctx.WriteYAML(http.StatusOK, result)
	}
}

func XMLResponseHandle() easierweb.ResponseHandle {
	return func(ctx *easierweb.Context, result any, err error) {
		if err != nil {
			if result != nil {
				ctx.WriteXML(http.StatusBadRequest, result)
				return
			}
			panic(err)
		}
		if result == nil {
			ctx.NoContent(http.StatusNoContent)
			return
		}
		ctx.WriteXML(http.StatusOK, result)
	}
}

func BytesResponseHandle() easierweb.ResponseHandle {
	return func(ctx *easierweb.Context, result any, err error) {
		if err != nil {
			if result != nil {
				ctx.Write(http.StatusBadRequest, result.([]byte))
				return
			}
			panic(err)
		}
		if result == nil {
			ctx.NoContent(http.StatusNoContent)
			return
		}
		ctx.Write(http.StatusOK, result.([]byte))
	}
}

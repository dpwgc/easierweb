package plugins

import (
	"github.com/dpwgc/easierweb"
	"net/http"
)

func JSONResponseHandle(ctx *easierweb.Context, result any, err error) {
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

func YAMLResponseHandle(ctx *easierweb.Context, result any, err error) {
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

func XMLResponseHandle(ctx *easierweb.Context, result any, err error) {
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

func BytesResponseHandle(ctx *easierweb.Context, result any, err error) {
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

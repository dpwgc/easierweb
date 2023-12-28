package plugins

import (
	"github.com/dpwgc/easierweb"
	"net/http"
)

func JSONResponseHandle(ctx *easierweb.Context, result any, err error) {
	if err != nil {
		errRes := make(map[string]string, 1)
		errRes["error"] = err.Error()
		ctx.WriteJSON(http.StatusBadRequest, errRes)
		return
	}
	if result == nil {
		ctx.Write(http.StatusNoContent, nil)
		return
	}
	ctx.WriteJSON(http.StatusOK, result)
}

func YAMLResponseHandle(ctx *easierweb.Context, result any, err error) {
	if err != nil {
		errRes := make(map[string]string, 1)
		errRes["error"] = err.Error()
		ctx.WriteYAML(http.StatusBadRequest, errRes)
		return
	}
	if result == nil {
		ctx.Write(http.StatusNoContent, nil)
		return
	}
	ctx.WriteYAML(http.StatusOK, result)
}

func XMLResponseHandle(ctx *easierweb.Context, result any, err error) {
	if err != nil {
		errRes := make(map[string]string, 1)
		errRes["error"] = err.Error()
		ctx.WriteXML(http.StatusBadRequest, errRes)
		return
	}
	if result == nil {
		ctx.Write(http.StatusNoContent, nil)
		return
	}
	ctx.WriteXML(http.StatusOK, result)
}

func HTMLResponseHandle(ctx *easierweb.Context, result any, err error) {
	if err != nil {
		ctx.WriteString(http.StatusBadRequest, err.Error())
		return
	}
	if result == nil {
		ctx.Write(http.StatusNoContent, nil)
		return
	}
	ctx.WriteHTML(http.StatusOK, result.(string))
}

func StringResponseHandle(ctx *easierweb.Context, result any, err error) {
	if err != nil {
		ctx.WriteString(http.StatusBadRequest, err.Error())
		return
	}
	if result == nil {
		ctx.Write(http.StatusNoContent, nil)
		return
	}
	ctx.WriteString(http.StatusOK, result.(string))
}

func BytesResponseHandle(ctx *easierweb.Context, result any, err error) {
	if err != nil {
		ctx.Write(http.StatusBadRequest, []byte(err.Error()))
		return
	}
	if result == nil {
		ctx.Write(http.StatusNoContent, nil)
		return
	}
	ctx.Write(http.StatusOK, result.([]byte))
}

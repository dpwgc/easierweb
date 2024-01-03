package easierweb

import (
	"net/http"
)

type Group struct {
	router      *Router
	path        string
	middlewares []Handle
}

func (r *Router) Group(path string, middlewares ...Handle) *Group {
	return &Group{
		router:      r,
		path:        path,
		middlewares: middlewares,
	}
}

// easier usage function

func (g *Group) EasyGET(path string, easyHandle any, middlewares ...Handle) *Group {
	g.router.EasyGET(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyHEAD(path string, easyHandle any, middlewares ...Handle) *Group {
	g.router.EasyHEAD(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyOPTIONS(path string, easyHandle any, middlewares ...Handle) *Group {
	g.router.EasyOPTIONS(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyPOST(path string, easyHandle any, middlewares ...Handle) *Group {
	g.router.EasyPOST(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyPUT(path string, easyHandle any, middlewares ...Handle) *Group {
	g.router.EasyPUT(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyPATCH(path string, easyHandle any, middlewares ...Handle) *Group {
	g.router.EasyPATCH(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyDELETE(path string, easyHandle any, middlewares ...Handle) *Group {
	g.router.EasyDELETE(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyAny(path string, easyHandle any, middlewares ...Handle) *Group {
	g.router.EasyAny(g.path+path, easyHandle, middlewares...)
	return g
}

// basic usage function

func (g *Group) GET(path string, handle Handle, middlewares ...Handle) *Group {
	g.router.GET(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) HEAD(path string, handle Handle, middlewares ...Handle) *Group {
	g.router.HEAD(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) OPTIONS(path string, handle Handle, middlewares ...Handle) *Group {
	g.router.OPTIONS(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) POST(path string, handle Handle, middlewares ...Handle) *Group {
	g.router.POST(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) PUT(path string, handle Handle, middlewares ...Handle) *Group {
	g.router.PUT(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) PATCH(path string, handle Handle, middlewares ...Handle) *Group {
	g.router.PATCH(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) DELETE(path string, handle Handle, middlewares ...Handle) *Group {
	g.router.DELETE(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) Any(path string, handle Handle, middlewares ...Handle) *Group {
	g.router.Any(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) WS(path string, handle Handle, middlewares ...Handle) *Group {
	g.router.WS(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) Static(path, dir string) *Group {
	g.router.Static(g.path+path, dir)
	return g
}

func (g *Group) StaticFS(path string, fs http.FileSystem) *Group {
	g.router.StaticFS(g.path+path, fs)
	return g
}

func (g *Group) Use(middlewares ...Handle) *Group {
	g.middlewares = append(g.middlewares, middlewares...)
	return g
}

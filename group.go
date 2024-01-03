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
	middlewares = append(g.middlewares, middlewares...)
	g.router.EasyGET(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyHEAD(path string, easyHandle any, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.EasyHEAD(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyOPTIONS(path string, easyHandle any, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.EasyOPTIONS(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyPOST(path string, easyHandle any, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.EasyPOST(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyPUT(path string, easyHandle any, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.EasyPUT(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyPATCH(path string, easyHandle any, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.EasyPATCH(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyDELETE(path string, easyHandle any, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.EasyDELETE(g.path+path, easyHandle, middlewares...)
	return g
}

func (g *Group) EasyAny(path string, easyHandle any, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.EasyAny(g.path+path, easyHandle, middlewares...)
	return g
}

// basic usage function

func (g *Group) GET(path string, handle Handle, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.GET(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) HEAD(path string, handle Handle, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.HEAD(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) OPTIONS(path string, handle Handle, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.OPTIONS(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) POST(path string, handle Handle, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.POST(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) PUT(path string, handle Handle, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.PUT(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) PATCH(path string, handle Handle, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.PATCH(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) DELETE(path string, handle Handle, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.DELETE(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) Any(path string, handle Handle, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
	g.router.Any(g.path+path, handle, middlewares...)
	return g
}

func (g *Group) WS(path string, handle Handle, middlewares ...Handle) *Group {
	middlewares = append(g.middlewares, middlewares...)
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

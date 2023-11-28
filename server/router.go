package server

import (
	"fmt"
	"net/http"
)

type (
	UsesRoutes interface {
		Routes()
	}
	Router struct {
		UsesRoutes
		Server             *server
		prefix             string
		Routes             []*Route
		domainMiddleware   []Middleware
		sequenceMiddleware []MiddlewareSequence[*Route]
	}
	UsesSubRoutes interface {
		SubRoutes()
	}
	Route struct {
		UsesSubRoutes
		Server             *server
		Router             *Router
		prefix             string
		subroutes          []*Subroute
		domainMiddleware   []Middleware
		sequenceMiddleware []MiddlewareSequence[*Subroute]
	}
	Subroute struct {
		Server      *server
		Router      *Router
		Route       *Route
		middleware  []Middleware
		prefix      string
		handler     *http.Handler
		handlerFunc http.HandlerFunc
	}
	MiddlewareSequence[T interface{}] struct {
		middleware []Middleware
		routes     []T
	}
	Middleware func(next http.Handler) http.Handler
)

func newRouter(prefix string, s *server) *Router {
	return &Router{
		Server:             s,
		prefix:             prefix,
		Routes:             []*Route{},
		domainMiddleware:   []Middleware{},
		sequenceMiddleware: []MiddlewareSequence[*Route]{},
	}
}

func (rtr *Router) Handle(prefix string, r Route) {
	rtr.Routes = append(rtr.Routes, &r)
}

func (rtr *Router) Use(middleware ...Middleware) {
	rtr.domainMiddleware = append(rtr.domainMiddleware, middleware...)
}

func (rtr *Router) UseSequence(middleware ...Middleware) {
	rtr.sequenceMiddleware = append(
		rtr.sequenceMiddleware,
		MiddlewareSequence[*Route]{
			middleware: middleware,
			routes:     rtr.Routes,
		},
	)
}

func (rtr *Router) addMiddleware(mw []Middleware, srs ...*Subroute) {
	if len(mw) == 0 || srs == nil {
		return
	}

	for _, sr := range srs {
		// Intializatize handler and append middleware
		if len(mw) > 1 {
			*sr.handler = mw[len(mw)-1](sr.handlerFunc)
			// Wrap handler in middleware in reverse order of intialization
			for i := len(mw) - 2; i >= len(mw)-2; i-- {
				*sr.handler = mw[i](*sr.handler)
			}
		} else {
			*sr.handler = mw[0](sr.handlerFunc)
		}

		fmt.Println("ok")
		//TODO: What is this?
		// Append middleware to existing handler
		// for i := len(mw) - 1; i >= len(mw)-1; i-- {
		// 	*sr.handler = mw[i](*sr.handler)
		// }
	}
}

func (rtr *Router) ApplyMiddleware() {
	// Apply subroute middleware
	for _, r := range rtr.Routes {
		for _, sr := range r.subroutes {
			sr.applyMiddleware()
		}
		// Apply route middleware
		r.applyMiddleware()
	}

	// Apply router sequential middleware to subroutes
	for _, smw := range rtr.sequenceMiddleware {
		for _, r := range smw.routes {
			rtr.addMiddleware(smw.middleware, r.subroutes...)
		}
	}
	// Apply router domain middleware to subroutes
	for _, r := range rtr.Routes {
		rtr.addMiddleware(rtr.domainMiddleware, r.subroutes...)
	}
}

func NewRoute(prefix string, subroutes ...Subroute) *Route {
	r := Route{
		Server:             Server,
		prefix:             Server.Router.prefix + prefix,
		subroutes:          []*Subroute{},
		domainMiddleware:   []Middleware{},
		sequenceMiddleware: []MiddlewareSequence[*Subroute]{},
		Router:             Server.Router,
	}
	return &r
}

func (r *Route) Handle(prefix string, h http.HandlerFunc) {
	r.subroutes = append(r.subroutes, newSubRoute(r.prefix+prefix, h, r))
}

func (r *Route) Use(middleware ...Middleware) {
	r.domainMiddleware = append(r.domainMiddleware, middleware...)
}

func (r *Route) UseSequence(middleware ...Middleware) {
	r.sequenceMiddleware = append(
		r.sequenceMiddleware,
		MiddlewareSequence[*Subroute]{
			middleware: middleware,
			routes:     r.subroutes,
		},
	)
}

func (r *Route) applyMiddleware() {
	// Apply sequential middleware to subroutes
	for _, s := range r.sequenceMiddleware {
		r.Router.addMiddleware(s.middleware, s.routes...)
	}
	// Apply domain middleware to subroutes
	for _, sr := range r.subroutes {
		r.Router.addMiddleware(r.domainMiddleware, sr)
	}
}

func newSubRoute(prefix string, h http.HandlerFunc, r *Route) *Subroute {
	return &Subroute{
		Server:      Server,
		Router:      Server.Router,
		Route:       r,
		middleware:  []Middleware{},
		prefix:      prefix,
		handler:     nil,
		handlerFunc: h,
	}
}

func (sr *Subroute) applyMiddleware() {
	if sr.handler == nil {
		h := http.Handler(sr.handlerFunc)
		sr.handler = &h
	}
	sr.Router.addMiddleware(sr.middleware, sr)
}

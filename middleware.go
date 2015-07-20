package main

import (
	"log"
	"net/http"
	"os"
)

// Middleware allows middleware handlers to be chained before and after the
// request handler.
type Middleware struct {
	mux  *http.ServeMux
	pre  []http.Handler
	post []http.Handler
}

// NewMiddleware constructs a new middleware type around a mux. The mux ServeHTTP is called
// after all the Pre middlewares have been called.
func NewMiddleware(mux *http.ServeMux) *Middleware {
	return &Middleware{
		mux:  mux,
		pre:  make([]http.Handler, 0),
		post: make([]http.Handler, 0),
	}
}

// AddPre adds a middleware handler that will be executed before the request handler.
func (m *Middleware) AddPre(h http.Handler) *Middleware {
	m.pre = append(m.pre, h)
	return m
}

// AddPost addsd a middleware handler that will be executed after the request handler.
func (m *Middleware) AddPost(h http.Handler) *Middleware {
	m.post = append(m.post, h)
	return m
}

// ServeHTTP implements the http.Handler interface
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, h := range m.pre {
		h.ServeHTTP(w, r)
	}
	m.mux.ServeHTTP(w, r)
	for _, h := range m.post {
		h.ServeHTTP(w, r)
	}
}

// logRequests is a simple POST middle ware that logs the request
func logRequests(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
}

// TemplateReloader is a PRE middleware that reloads templates on each request if the
// DEVELOPER environment variable is set to anything that isn't an empty string.
type TemplateReloader struct {
	reloadables []TemplateReloadable
}

// Interface for types that can have their templates reloaded
type TemplateReloadable interface {
	reloadTemplates()
}

// NewTemplatedReloader returns an initialized reloader.
func NewTemplateReloader() *TemplateReloader {
	return &TemplateReloader{
		reloadables: make([]TemplateReloadable, 0),
	}
}

// Add appends a type to have its templates reloaded.
func (tr *TemplateReloader) Add(t TemplateReloadable) *TemplateReloader {
	tr.reloadables = append(tr.reloadables, t)
	return tr
}

// ServeHTTP checks the `DEVELOPER` environment variable and if set, instructs all
// registered types to reload for each request.
func (tr *TemplateReloader) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// TODO: move template lifecycle management outside of the controlers

	if os.Getenv("DEVELOPER") != "" {
		log.Println("Reloading templates.")
		for _, r := range tr.reloadables {
			r.reloadTemplates()
		}
	}
}

// methodRouter extends the default net/http routing so that you can specify separate handlers
// for each request method type (GET, POST, PUT, DELETE ...). This is not a middleware but is
// cosely related so in the same file.
type methodRouter struct {
	handlers map[string]http.Handler
}

// newMethodRouter takes the server mux that will dispatch to this router based on the matched path.
func newMethodRouter(mux *http.ServeMux, path string) *methodRouter {
	mr := &methodRouter{
		handlers: map[string]http.Handler{},
	}
	mux.Handle(path, mr)
	return mr
}

// Add saves a handler along with the method is should be called for.
func (mr *methodRouter) Add(method string, handler http.Handler) {
	mr.handlers[method] = handler
}

// ServeHTTP checks for a registered handler for the request method and calls it. If there
// isn't a registerd handler for this method then a 404 response is generated.
func (mr *methodRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := mr.handlers[r.Method]
	if h != nil {
		h.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Invalid method for route: "+r.Method, http.StatusNotFound)
}

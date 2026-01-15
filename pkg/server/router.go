package server

import (
	"io/fs"
	"net/http"

	"github.com/ParkerGits/go-backend-starter/pkg/config"
	"github.com/ParkerGits/go-backend-starter/pkg/httplib"
	"github.com/ParkerGits/go-backend-starter/pkg/ogimage"
	"github.com/ParkerGits/go-backend-starter/pkg/templates"
	"github.com/ParkerGits/go-backend-starter/web/static"
)

type router struct {
	router         *http.ServeMux
	config         *config.Server
	ogImageHandler *ogimage.Handler
}

func newRouter(config *config.Server, generator *ogimage.Generator) *router {
	router := &router{
		router:         http.NewServeMux(),
		config:         config,
		ogImageHandler: ogimage.NewHandler(generator),
	}
	router.registerDefaultRoutes()
	return router
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.withDefaultMiddleware(r.router).ServeHTTP(w, req)
}

func (r *router) registerDefaultRoutes() {
	// Static files (embedded)
	staticFS, _ := fs.Sub(static.Assets, ".")
	r.router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Preview route for rod to screenshot
	r.registerRoute("GET /_preview", r.previewHandler())

	// Main OG image endpoint
	r.registerRoute("GET /", r.ogImageHandler)
}

func (r *router) previewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		title := req.URL.Query().Get("title")
		subtitle := req.URL.Query().Get("subtitle")

		params := templates.OGImageParams{
			Title:     title,
			Subtitle:  subtitle,
			Handle:    "@lordprkr",
			AvatarURL: "/static/avatar.png",
		}

		w.Header().Set("Content-Type", "text/html")
		templates.OGImage(params).Render(req.Context(), w)
	})
}

func (r *router) registerRoute(pattern string, baseHandler http.Handler, middlewares ...httplib.Middleware) {
	withMiddleware := httplib.CreateStack(middlewares...)
	r.router.Handle(pattern, withMiddleware(baseHandler))
}

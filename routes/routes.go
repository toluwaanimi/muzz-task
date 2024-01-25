package routes

import (
	"api/controllers"
	"github.com/go-chi/chi"
	swaggerMiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/rs/cors"
	"net/http"
)

func configureCORS() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
	}).Handler
}

func configureSwaggerDocs(router *chi.Mux) {
	redocOpts := swaggerMiddleware.RedocOpts{SpecURL: "/docs/swagger.yaml", Path: "/docs"}
	shareDocs := swaggerMiddleware.Redoc(redocOpts, nil)
	router.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))
	router.Handle("/docs", shareDocs)
}

func ConfigureRoutes(router *chi.Mux, controller *controllers.Controller) {
	router.Use(configureCORS())
	configureSwaggerDocs(router)
	router.Route("/", func(r chi.Router) {
		r.Get("/", controller.Home)
		r.Post("/user/create", controller.RegisterUser)
		r.Post("/login", controller.LoginUser)
		r.Get("/user", controller.GetUser)
		r.Get("/discover", controller.DiscoverUsers)
		r.Post("/swipe", controller.SwipeUser)
	})
}

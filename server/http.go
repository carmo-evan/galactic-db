package server

import (
	"alchemy/galacticdb/db"
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

func GetRouter(s db.Store) http.Handler {
	r := chi.NewRouter()
	routesFunc := routes(s)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/spaceship", routesFunc)
	return r
}

func routes(s db.Store) func(router chi.Router) {
	return func (router chi.Router) {
		router.Get("/", getSpaceshipsRoute(s))
		router.Post("/", createSpaceshipRoute(s))
		router.Patch("/", updateSpaceshipRoute(s))
		router.Route("/{id}", func(router chi.Router) {
			router.Use(SpaceshipContext)
			router.Get("/", getSpaceshipRoute(s))
			router.Delete("/", deleteSpaceshipRoute(s))
		})
	}
}


func SpaceshipContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		if idParam == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("id is required")))
			return
		}
		id, err := strconv.Atoi(idParam)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("id must be a number")))
			return
		}
		ctx := context.WithValue(r.Context(), "spaceshipId", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

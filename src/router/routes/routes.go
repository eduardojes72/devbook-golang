package routes

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI                    string
	Method                 string
	Function               func(http.ResponseWriter, *http.Request)
	RequiresAuthentication bool
}

func Config(r *mux.Router) *mux.Router {
	routes := routesUsers
	routes = append(routes, routeLogin)

	for _, route := range routes {
		if route.RequiresAuthentication {
			r.HandleFunc(
				route.URI,
				middlewares.Logger(middlewares.IsAuthenticated(route.Function)),
			).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, middlewares.Logger(route.Function)).Methods(route.Method)
		}

	}

	return r
}

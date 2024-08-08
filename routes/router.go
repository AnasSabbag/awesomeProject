package routes

import (
	"awesomeProject/handler"
	"awesomeProject/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateRoutes() http.Handler {
	router := mux.NewRouter()
	router.Handle("/", handler.BasicRoutes(router))
	router.Handle("/", handler.UserRoutes(router))

	//apiPrefix := router.PathPrefix("/api").Subrouter()
	router1 := router.PathPrefix("/api").Subrouter()

	authRoutes(router1)
	return middlewares.EnableCORS(router)
}

func authRoutes(router *mux.Router) {

	router.Use(middlewares.AuthHandler)
	router.Use(middlewares.PermissionMiddleware)
	router.Handle("/", handler.AdminRoutes(router))
	router.Handle("/", handler.UserProfileRoutes(router))

}

//func hotelRoutes(router *mux.Router) {}

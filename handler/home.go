package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func BasicRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/api", Basic)
	return router
}

func HomeRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/api/home", home)
	return router
}

func Basic(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "hello test")
	if err != nil {
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode("welcome!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("home:encoding error", err)
		return
	}

}

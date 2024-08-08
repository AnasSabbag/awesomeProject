package main

import (
	"awesomeProject/global"
	"awesomeProject/routes"
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("server starting.. ")
	global.InitDb()
	port := os.Getenv("SERVER_PORT")
	srv := routes.CreateRoutes()
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handlers.RecoveryHandler()(srv)); err != nil {
		log.Println("Listen and serve error : ", err)
	}
	fmt.Println("server successfully connected")
}

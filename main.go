package main

import (
	"fmt"
	"log"
	"net/http"
	"templator/handlers"
	"templator/handlers/admin"
	"templator/services"
)

func main() {
	// Registering route handlers
	http.HandleFunc("/admin/", admin.PanelHandler)
	http.HandleFunc("/admin/login", admin.LoginHandler)
	http.HandleFunc("/templates", handlers.TemplatesHandler)
	http.HandleFunc("/params", handlers.ParamsHandler)
	http.HandleFunc("/generate", handlers.GeneratorHandler)
	http.HandleFunc("/welcome", admin.Welcome)
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Running HTTP server
	var httpPort = services.GetHttpPort()
	log.Printf("Running HTTP-server on port %d", httpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil))
}

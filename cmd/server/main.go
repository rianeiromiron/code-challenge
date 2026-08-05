package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rianeiromiron/code-challenge/internal/handlers"
)

func main() {
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/clima", handlers.EnableCORS(handlers.ClimaHandler))

	addr := ":8080"
	fmt.Printf("🌤️  Servidor iniciado en http://localhost%s\n", addr)
	fmt.Printf("   Abre tu navegador en http://localhost%s\n", addr)
	fmt.Printf("   API: curl \"http://localhost%s/clima?pais=Mexico&limit=2\"\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

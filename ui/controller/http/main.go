package http

import (
	"github.com/DJustProgrameR/internshipPVZ/ui/app"
	"log"
	"net/http"
)

func main() {
	container := app.NewContainer()
	server := app.NewHTTPServer(container)
	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", server.Router()); err != nil {
		log.Fatal("HTTP server failed: ", err)
	}
}

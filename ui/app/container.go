package app

import (
	"go.uber.org/fx"
	"log"
	"net/http"
)

func main() {
	fx.New(
		fx.Provide(
			app.NewContainer,  // DI container builder
			app.NewHTTPServer, // provides *HTTPServer
		),
		fx.Invoke(func(server *app.HTTPServer) {
			log.Println("Starting HTTP server on :8080")
			if err := http.ListenAndServe(":8080", server.Router()); err != nil {
				log.Fatal("HTTP server failed: ", err)
			}
		}),
	).Run()
}

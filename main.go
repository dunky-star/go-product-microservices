package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dunky-star/protobv1/handlers"
	"github.com/dunky-star/protobv1/product-api/data"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

//var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {
    l := log.New(os.Stdout, "product-api", log.LstdFlags)
	v := data.NewValidation()
	// create the handlers
	hProducts := handlers.NewProducts(l, v)

	// Using Gorilla Mux and subrouters for routing.
	sm := mux.NewRouter()
    getRouter := sm.Methods(http.MethodGet).Subrouter()
	// Create a new serve mux and register the handlers
	getRouter.HandleFunc("/products", hProducts.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", hProducts.ListSingle)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", hProducts.Create)
	postRouter.Use(hProducts.MiddlewareValidateProduct)
	

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{:id[0-9]+}", hProducts.Update)
	putRouter.Use(hProducts.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", hProducts.Delete)
	
	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// create a new server
	s := &http.Server{
		Addr:         ":9090",      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func(){
		err := s.ListenAndServe()
		if err != nil{
			l.Fatal(err)
		}

	}()
	
	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    signal.Notify(c, os.Kill)

    // Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

    // gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dunky-star/protobv1/handlers"
	"github.com/gorilla/mux"
)

func main() {
    l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hProducts := handlers.NewProducts(l)

	// Using Gorilla Mux and subrouters for routing.
	sm := mux.NewRouter()
    getRouter := sm.Methods(http.MethodGet).Subrouter()
	// Create a new serve mux and register the handlers
	getRouter.HandleFunc("/products", hProducts.GetProducts)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", hProducts.AddProduct)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{:id[0-9]+}", hProducts.UpdateProducts)
	

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func(){
		err := s.ListenAndServe()
		if err != nil{
			l.Fatal(err)
		}

	}()
	
	sigChan := make(chan os.Signal)
    signal.Notify(sigChan, os.Interrupt)
    signal.Notify(sigChan, os.Kill)
    sig := <- sigChan
    log.Println("Received terminate, graceful shutdown", sig)

    tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
	s.Shutdown(tc)

}
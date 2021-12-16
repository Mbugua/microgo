package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mbugua/micro/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	p := handlers.NewProducts(l)
	sm := http.NewServeMux()
	sm.Handle("/", p)
	sm.Handle("/products", p)

	// http.ListenAndServe(":8000", sm)
	s := &http.Server{
		Addr:         ":8000",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sch := make(chan os.Signal)
	signal.Notify(sch, os.Interrupt)
	signal.Notify(sch, os.Kill)
	sig := <-sch

	l.Println("Received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)

}

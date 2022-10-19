package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Başarılı"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/deneme", TestEndpoint).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Sunucu başarı ile başlatıldı")

	<-done
	log.Print("Sunucu kapatıldı")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {

		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Sunucu kapanma hatası:%+v", err)
	}
	log.Print("Sunucudan Çıkıldı")
}
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"spe-trx-gateway/config"
	"spe-trx-gateway/controllers"
	"spe-trx-gateway/routers"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	db, err := config.NewDB(ctx)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer db.Close()

	srvCtl := controllers.NewServer(db.Pool)

	r, port := routers.Route(srvCtl)

	httpSrv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Println("listening on :" + port)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctxShutdown); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("server stopped gracefully")
}

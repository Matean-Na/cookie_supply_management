package main

import (
	"context"
	"cookie_supply_management/core/server"
	"fmt"
	"log"
	"os/signal"

	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := server.Server{}
	r, e := s.Init()
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	s.Run(r)

	<-ctx.Done()

	stop()
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

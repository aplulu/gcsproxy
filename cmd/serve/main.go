package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aplulu/gcsproxy/internal/config"
	"github.com/aplulu/gcsproxy/internal/infrastructure/http"
)

func main() {
	if err := config.LoadConf(); err != nil {
		panic(err)
	}

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-quitCh
		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := http.StopServer(shutdownCtx); err != nil {
			log.Println(fmt.Sprintf("command.ServeCommand: failed to stop server: %+v", err))
			os.Exit(1)
			return
		}
	}()

	log.Println("Starting server...")
	if err := http.RunServer(); err != nil {
		log.Println(fmt.Sprintf("command.ServeCommand: failed to start server: %+v", err))
		os.Exit(1)
	}
}

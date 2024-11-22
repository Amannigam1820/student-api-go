package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Amannigam1820/student-api-go/internal/config"
	"github.com/Amannigam1820/student-api-go/internal/http/handler/student"
)

func main() {
	//load config
	cfg := config.MustLoad()

	// database setup
	// setup router
	router := http.NewServeMux() // router initialized

	router.HandleFunc("POST /api/students", student.New())

	// setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1) // create a channel type signal

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // signal notify channel about the signal

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}

	}()

	<-done

	slog.Info("Shutting Down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server ShutDown SuccessFully..")

}

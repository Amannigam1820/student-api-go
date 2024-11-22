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
	"github.com/Amannigam1820/student-api-go/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux() // router initialized

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetAllStudent(storage))

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

	errs := server.Shutdown(ctx)
	if errs != nil {
		slog.Error("failed to shutdown server", slog.String("error", errs.Error()))
	}
	slog.Info("Server ShutDown SuccessFully..")

}

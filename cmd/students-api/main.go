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
	"github.com/Amannigam1820/student-api-go/internal/http/handler/user"
	"github.com/Amannigam1820/student-api-go/internal/middleware"
	"github.com/Amannigam1820/student-api-go/internal/storage/sqlite"
	"github.com/rs/cors"
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

	// User Registration Routes

	router.HandleFunc("POST /api/users/register", user.RegisterUser(storage))
	router.HandleFunc("POST /api/users/login", user.Login(storage))
	router.HandleFunc("POST /api/users/logout", user.Logout)

	// Students Routes

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetAllStudent(storage))

	//router.HandleFunc("GET /api/students/filter", student.GetStudentByFilter(storage))

	// router.Handle("/api/students", middleware.AuthMiddleware(http.HandlerFunc(student.GetAllStudent(storage))))
	router.Handle("/api/user/me", middleware.AuthMiddleware(http.HandlerFunc(user.GetLoggedInUser(storage))))

	router.HandleFunc("DELETE /api/students/{id}", student.DeleteStudent(storage))
	router.HandleFunc("PUT /api/student/{id}", student.UpdateStudent(storage))

	// setup server

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: corsHandler,
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

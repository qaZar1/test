package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/qaZar1/test/autogen"
	"github.com/qaZar1/test/internal/config"
	"github.com/qaZar1/test/internal/repository"
	"github.com/qaZar1/test/internal/service"
)

func main() {
	cfg := config.New()

	repo := repository.NewRepository(
		repository.Config{
			Hostname: cfg.Hostname,
			Database: cfg.Database,
			User:     cfg.User,
			Password: cfg.Password,
			Port:     cfg.Port,
		})

	transport := service.NewTransport(repo)

	router := chi.NewRouter()
	router.Handle("/*", autogen.Handler(transport))

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	srv := service.New(server)
	if err := srv.Run(); err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch

	srv.Stop()
}

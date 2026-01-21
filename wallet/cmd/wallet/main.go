package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/qaZar1/test/wallet/autogen"
	"github.com/qaZar1/test/wallet/internal/config"
	"github.com/qaZar1/test/wallet/internal/postgres"
	"github.com/qaZar1/test/wallet/internal/service"
)

func main() {
	cfg := config.New()

	transport := service.NewTransport(postgres.Config{
		Hostname: cfg.Hostname,
		Database: cfg.Database,
		User:     cfg.User,
		Password: cfg.Password,
		Port:     cfg.Port,
	})

	router := http.NewServeMux()
	router.Handle("/", autogen.Handler(transport))

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
		fmt.Println("Service stopped")
	}()
	fmt.Println("Service started")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		syscall.SIGABRT,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	<-ch

	if err := server.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}

package service

import (
	"fmt"
	"net/http"

	"github.com/qaZar1/test/internal/repository"
)

type Service struct {
	server *http.Server
	repo   repository.IRepository
}

func New(server *http.Server) *Service {
	return &Service{
		server: server,
	}
}

func (s *Service) Run() error {
	fmt.Println("Service started")
	return s.server.ListenAndServe()
}

func (s *Service) Stop() {
	fmt.Println("Service stopped")
	if err := s.repo.Close(); err != nil {
		fmt.Printf("Can not close database: %s", err)
	}
}

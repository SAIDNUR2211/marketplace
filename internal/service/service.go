package service

import (
	"marketplace/internal/contracts"
	"os"

	"github.com/rs/zerolog"
)

type Service struct {
	repository contracts.RepositoryI
	logger     zerolog.Logger
}

func NewService(repository contracts.RepositoryI) *Service {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("entity", "service").Logger()
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

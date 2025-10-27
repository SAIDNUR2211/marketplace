package service

import (
	"marketplace/internal/contracts"
	"os"

	"github.com/rs/zerolog"
)

type Service struct {
	repository contracts.RepositoryI
	cache      contracts.CacheI // НОВАЯ СТРОКА
	logger     zerolog.Logger
}

// ИЗМЕНЕНИЕ: NewService теперь также принимает CacheI
func NewService(repository contracts.RepositoryI, cache contracts.CacheI) *Service {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("entity", "service").Logger()
	return &Service{
		repository: repository,
		cache:      cache, // НОВАЯ СТРОКА
		logger:     logger,
	}
}

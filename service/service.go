package service

import (
	"context"

	"github.com/go-kit/kit/log"
)

// Service describes a service that generates QR Code images.
type Service interface {
	Generate(ctx context.Context, data string) ([]byte, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(logger log.Logger) Service {
	var svc Service
	{
		svc = NewService()
		svc = LoggingMiddleware(logger)(svc)
	}
	return svc
}

// NewService constructs a service
func NewService() Service {
	return service{}
}

type service struct{}

// Generate encodes the input data and returns the
func (s service) Generate(ctx context.Context, data string) ([]byte, error) {
	return nil, nil
}

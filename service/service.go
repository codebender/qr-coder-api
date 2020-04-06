package service

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	qrcode "github.com/skip2/go-qrcode"
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
	err := validateData(data)

	if err != nil {
		return nil, err
	}

	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("An error occured while endcoding the data, err=%w", err)
	}

	return png, nil
}

func validateData(data string) error {
	if data == "" {
		return fmt.Errorf("Data is required")
	}

	return nil
}

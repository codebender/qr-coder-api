package transport

import (
	"context"
	"time"

	qrcodev1 "github.com/codebender/qrcode-api/proto/codebender/qrcode/v1"
	"github.com/codebender/qrcode-api/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// Set collects all of the endpoints that compose an QR Code service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	GenerateEndpoint endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc service.Service, logger log.Logger) Set {
	var generateEndpoint endpoint.Endpoint
	{
		generateEndpoint = MakeGenerateEndpoint(svc)
		generateEndpoint = LoggingMiddleware(log.With(logger, "method", "Generate"))(generateEndpoint)
	}

	return Set{
		GenerateEndpoint: generateEndpoint,
	}
}

// MakeGenerateEndpoint constructs a Generate endpoint wrapping the service.
func MakeGenerateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*qrcodev1.GenerateRequest)

		qrCode, err := s.Generate(ctx, req.Data)
		if err != nil {
			return nil, err
		}

		resp := &qrcodev1.GenerateResponse{QrCode: qrCode}

		return resp, nil
	}
}

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				logger.Log("transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)

		}
	}
}

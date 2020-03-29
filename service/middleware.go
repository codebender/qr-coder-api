package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{next, logger}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) log(ctx context.Context, err error, params ...interface{}) {
	logLevel := level.Info(mw.logger)
	if err != nil {
		logLevel = level.Error(mw.logger)
	}
	logLevel.Log(params...)
}

func (mw loggingMiddleware) Generate(ctx context.Context, data string) ([]byte, error) {
	r, e := mw.next.Generate(ctx, data)
	mw.log(ctx, e, "method", "Generate", "data", data)

	return r, e
}

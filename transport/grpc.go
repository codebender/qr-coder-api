package transport

import (
	"context"

	qrcodev1 "github.com/codebender/qrcode-api/proto/codebender/qrcode/v1"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	generate grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC AddServer.
func NewGRPCServer(endpoints Set, logger log.Logger) qrcodev1.QRCodeAPIServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		generate: grpctransport.NewServer(
			endpoints.GenerateEndpoint,
			decodeGRPCRequest,
			encodeGRPCResponse,
			options...,
		),
	}
}

func (s *grpcServer) Generate(ctx context.Context, req *qrcodev1.GenerateRequest) (*qrcodev1.GenerateResponse, error) {
	_, rep, err := s.generate.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*qrcodev1.GenerateResponse), nil
}

// decodeGRPCRequest returns the gRPC request.
func decodeGRPCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return grpcReq, nil
}

// encodeGRPCResponse returns the gRPC response.
func encodeGRPCResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	return grpcResp, nil
}

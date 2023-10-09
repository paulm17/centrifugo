package proxy

import (
	"context"
	"fmt"
	"time"

	"github.com/centrifugal/centrifugo/v5/internal/proxyproto"

	"google.golang.org/grpc"
)

// GRPCRPCProxy ...
type GRPCRPCProxy struct {
	config Config
	client proxyproto.CentrifugoProxyClient
}

var _ RPCProxy = (*GRPCRPCProxy)(nil)

// NewGRPCRPCProxy ...
func NewGRPCRPCProxy(p Config) (*GRPCRPCProxy, error) {
	host, err := getGrpcHost(p.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("error getting grpc host: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.Timeout))
	defer cancel()
	dialOpts, err := getDialOpts(p)
	if err != nil {
		return nil, fmt.Errorf("error creating GRPC dial options: %v", err)
	}
	conn, err := grpc.DialContext(ctx, host, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("error connecting to GRPC proxy server: %v", err)
	}
	return &GRPCRPCProxy{
		config: p,
		client: proxyproto.NewCentrifugoProxyClient(conn),
	}, nil
}

// ProxyRPC ...
func (p *GRPCRPCProxy) ProxyRPC(ctx context.Context, req *proxyproto.RPCRequest) (*proxyproto.RPCResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(p.config.Timeout))
	defer cancel()
	return p.client.RPC(grpcRequestContext(ctx, p.config), req, grpc.ForceCodec(grpcCodec))
}

// Protocol ...
func (p *GRPCRPCProxy) Protocol() string {
	return "grpc"
}

// UseBase64 ...
func (p *GRPCRPCProxy) UseBase64() bool {
	return p.config.BinaryEncoding
}

// IncludeMeta ...
func (p *GRPCRPCProxy) IncludeMeta() bool {
	return p.config.IncludeConnectionMeta
}

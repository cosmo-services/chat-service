package pkg

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type GrpcClient struct {
	opts        []grpc.DialOption
	connections []*grpc.ClientConn
	logger      Logger
}

func NewGrpcClient(logger Logger) *GrpcClient {
	return &GrpcClient{
		opts: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(4*1024*1024),
				grpc.MaxCallSendMsgSize(4*1024*1024),
			),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             3 * time.Second,
				PermitWithoutStream: true,
			}),
		},
		connections: make([]*grpc.ClientConn, 0),
		logger:      logger,
	}
}

func (c *GrpcClient) Connect(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, c.opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client for %s: %w", address, err)
	}

	c.connections = append(c.connections, conn)

	return conn, nil
}

func (c *GrpcClient) CloseAllConnections() {
	for _, conn := range c.connections {
		if err := conn.Close(); err != nil {
			c.logger.Error(err.Error())
		}
	}

	c.logger.Info("all gRPC connections closed")
}

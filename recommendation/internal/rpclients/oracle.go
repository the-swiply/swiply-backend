package rpclients

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/the-swiply/swiply-backend/recommendation/internal/pb/oracle"
)

type OracleClient struct {
	conn   *grpc.ClientConn
	client oracle.OracleClient
}

func NewOracleClient(addr string) (*OracleClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("can't dial: %w", err)
	}

	return &OracleClient{
		conn:   conn,
		client: oracle.NewOracleClient(conn),
	}, nil
}

func (o *OracleClient) CloseConn() error {
	return o.conn.Close()
}

func (o *OracleClient) RetrainLFMv1(ctx context.Context) error {
	_, err := o.client.RetrainLFMv1(ctx, &oracle.RetrainLFMv1Request{})
	return err
}

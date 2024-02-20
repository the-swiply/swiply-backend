package rpclients

import (
	"context"
	"fmt"
	"github.com/the-swiply/swiply-backend/recommendation/internal/pb/oracle"
	"google.golang.org/grpc"
)

type OracleClient struct {
	conn   *grpc.ClientConn
	client oracle.OracleClient
}

func NewOracleClient(addr string) (*OracleClient, error) {
	var opts []grpc.DialOption

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

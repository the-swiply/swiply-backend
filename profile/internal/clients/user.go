package clients

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/the-swiply/swiply-backend/profile/internal/pb/user"
)

type User struct {
	conn   *grpc.ClientConn
	client user.UserClient
}

func NewUserClient(addr string) (*User, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("can't dial: %w", err)
	}

	return &User{
		conn:   conn,
		client: user.NewUserClient(conn),
	}, nil
}

func (c *User) CloseConn() error {
	return c.conn.Close()
}

func (c *User) ValidateAuthorizationCode(ctx context.Context, email, code string) (bool, error) {
	resp, err := c.client.ValidateAuthorizationCode(ctx, &user.ValidateAuthorizationCodeRequest{
		Email: email,
		Code:  code,
	})

	if err != nil {
		return false, err
	}

	return resp.IsCorrect, nil
}

func (c *User) SendAuthorizationCode(ctx context.Context, email string) error {
	_, err := c.client.SendAuthorizationCode(ctx, &user.SendAuthorizationCodeRequest{
		Email: email,
	})

	return err
}

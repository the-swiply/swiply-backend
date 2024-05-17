package clients

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/the-swiply/swiply-backend/profile/internal/pb/notification"
)

type Notification struct {
	conn   *grpc.ClientConn
	client notification.NotificationClient

	s2sToken string
}

func NewNotificationClient(addr, s2sToken string) (*Notification, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("can't dial: %w", err)
	}

	return &Notification{
		conn:     conn,
		client:   notification.NewNotificationClient(conn),
		s2sToken: "Bearer " + s2sToken,
	}, nil
}

func (n *Notification) CloseConn() error {
	return n.conn.Close()
}

func (n *Notification) Send(ctx context.Context, userID uuid.UUID, content string) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "s2s-authorization", n.s2sToken)
	_, err := n.client.Send(ctx, &notification.SendRequest{
		Id:      userID.String(),
		Content: content,
	})

	return err
}

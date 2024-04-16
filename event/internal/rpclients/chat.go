package rpclients

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/event/internal/converter"
	"github.com/the-swiply/swiply-backend/event/internal/pb/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type ChatClient struct {
	conn   *grpc.ClientConn
	client chat.ChatClient

	s2sToken string
}

func NewChatClient(addr string, s2sToken string) (*ChatClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("can't dial: %w", err)
	}

	return &ChatClient{
		conn:     conn,
		client:   chat.NewChatClient(conn),
		s2sToken: "Bearer " + s2sToken,
	}, nil
}

func (c *ChatClient) CloseConn() error {
	return c.conn.Close()
}

func (c *ChatClient) CreateChat(ctx context.Context, members []uuid.UUID) (int64, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "s2s-authorization", c.s2sToken)
	resp, err := c.client.CreateChat(ctx, &chat.CreateChatRequest{Members: converter.UUIDsToStrings(members)})

	if err != nil {
		return 0, err
	}

	return resp.GetChatId(), nil
}

func (c *ChatClient) AddChatMembers(ctx context.Context, chatID int64, members []uuid.UUID) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "s2s-authorization", c.s2sToken)
	_, err := c.client.AddChatMembers(ctx, &chat.AddChatMembersRequest{
		ChatId:  chatID,
		Members: converter.UUIDsToStrings(members),
	})

	return err
}

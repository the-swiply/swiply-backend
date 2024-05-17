package rpclients

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/the-swiply/swiply-backend/recommendation/internal/converter"
	"github.com/the-swiply/swiply-backend/recommendation/internal/domain"
	"github.com/the-swiply/swiply-backend/recommendation/internal/pb/profile"
)

type ProfileClient struct {
	conn   *grpc.ClientConn
	client profile.ProfileClient

	s2sToken string
}

func NewProfileClient(addr string, s2sToken string) (*ProfileClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("can't dial: %w", err)
	}

	return &ProfileClient{
		conn:     conn,
		client:   profile.NewProfileClient(conn),
		s2sToken: "Bearer " + s2sToken,
	}, nil
}

func (o *ProfileClient) CloseConn() error {
	return o.conn.Close()
}

func (o *ProfileClient) GetProfiles(ctx context.Context, after time.Time) ([]domain.Profile, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "s2s-authorization", o.s2sToken)
	profilesResp, err := o.client.ListProfiles(ctx, &profile.ListProfilesRequest{UpdatedAfter: timestamppb.New(after)})
	if err != nil {
		return nil, err
	}

	return converter.ProfilesFromProtoToDomain(profilesResp.GetProfiles()), nil
}

func (o *ProfileClient) GetInteractions(ctx context.Context, after time.Time) ([]domain.Interaction, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "s2s-authorization", o.s2sToken)
	interactionsResp, err := o.client.ListInteractions(ctx, &profile.ListInteractionsRequest{After: timestamppb.New(after)})
	if err != nil {
		return nil, err
	}

	return converter.InteractionsFromProtoToDomain(interactionsResp.GetInteractions()), nil
}

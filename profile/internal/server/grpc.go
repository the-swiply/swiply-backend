package server

import (
	"context"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc/reflection"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"
	"github.com/the-swiply/swiply-backend/pkg/houston/grut"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/the-swiply/swiply-backend/profile/internal/converter"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
	"github.com/the-swiply/swiply-backend/profile/internal/service"
	"github.com/the-swiply/swiply-backend/profile/pkg/api/profile"
)

type GRPCServer struct {
	*grpc.Server

	profileServer *profileServer
	photoServer   *photoServer
}

type profileServer struct {
	profile.UnimplementedProfileServer
	service *service.ProfileService
}

func (p *profileServer) Create(ctx context.Context, req *profile.CreateProfileRequest) (*profile.CreateProfileResponse, error) {
	prof := converter.ProfileFromProtoToDomain(&profile.UserProfile{
		Id:               auf.ExtractUserIDFromContext[uuid.UUID](ctx).String(),
		Email:            req.Email,
		Name:             req.Name,
		Interests:        req.Interests,
		BirthDay:         req.BirthDay,
		Gender:           req.Gender,
		Info:             req.Info,
		SubscriptionType: req.SubscriptionType,
		Location:         req.Location,
	})

	err := p.service.Create(ctx, prof)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &profile.CreateProfileResponse{Id: prof.ID.String()}, nil
}

func (p *profileServer) Update(ctx context.Context, req *profile.UpdateProfileRequest) (*profile.UpdateProfileResponse, error) {
	prof := converter.ProfileFromProtoToDomain(&profile.UserProfile{
		Id:               auf.ExtractUserIDFromContext[uuid.UUID](ctx).String(),
		Name:             req.Name,
		Interests:        req.Interests,
		BirthDay:         req.BirthDay,
		Gender:           req.Gender,
		Info:             req.Info,
		SubscriptionType: req.SubscriptionType,
		Location:         req.Location,
	})

	err := p.service.Update(ctx, prof)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &profile.UpdateProfileResponse{Id: prof.ID.String()}, nil
}

func (p *profileServer) Get(ctx context.Context, req *profile.GetProfileRequest) (*profile.GetProfileResponse, error) {
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "id must have uuid format")
	}
	prof, err := p.service.Get(ctx, userID)
	if err == domain.ErrEntityIsNotExists {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &profile.GetProfileResponse{
		UserProfile: converter.ProfileFromDomainToProto(prof),
	}, nil
}

func (p *profileServer) WhoAmI(ctx context.Context, _ *profile.WhoAmIRequest) (*profile.WhoAmIResponse, error) {
	return &profile.WhoAmIResponse{
		Id: auf.ExtractUserIDFromContext[uuid.UUID](ctx).String(),
	}, nil
}

func (p *profileServer) GetRecommendations(ctx context.Context, req *profile.GetRecommendationsRequest) (*profile.GetRecommendationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "need implementation")
}

func (p *profileServer) Interaction(ctx context.Context, req *profile.InteractionRequest) (*profile.InteractionResponse, error) {
	if err := uuid.Validate(req.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, "id must have uuid format")
	}

	interaction := converter.InteractionFromProtoToDomain(&profile.Interaction{
		From: auf.ExtractUserIDFromContext[uuid.UUID](ctx).String(),
		To:   req.Id,
		Type: req.Type,
	})

	err := p.service.CreateInteraction(ctx, interaction)
	if err != nil {
		return nil, err
	}

	return &profile.InteractionResponse{}, nil
}

func (p *profileServer) Liked(ctx context.Context, _ *profile.LikedRequest) (*profile.LikedResponse, error) {
	userIDs, err := p.service.GetLikedProfiles(ctx, auf.ExtractUserIDFromContext[uuid.UUID](ctx))
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	resp := &profile.LikedResponse{}
	for _, userID := range userIDs {
		resp.Ids = append(resp.Ids, userID.String())
	}

	return resp, nil
}

func (p *profileServer) LikedMe(ctx context.Context, _ *profile.LikedMeRequest) (*profile.LikedMeResponse, error) {
	prof, err := p.service.Get(ctx, auf.ExtractUserIDFromContext[uuid.UUID](ctx))
	if err == domain.ErrEntityIsNotExists {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	if prof.Subscription != domain.SubscriptionTypePrimary {
		return nil, status.Error(codes.PermissionDenied, "user hasn't primary subscription")
	}

	userIDs, err := p.service.GetLikedProfiles(ctx, auf.ExtractUserIDFromContext[uuid.UUID](ctx))
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	resp := &profile.LikedMeResponse{}
	for _, userID := range userIDs {
		resp.Ids = append(resp.Ids, userID.String())
	}

	return resp, nil
}

func (p *profileServer) ListInterests(ctx context.Context, _ *profile.ListInterestsRequest) (*profile.ListInterestsResponse, error) {
	interests, err := p.service.ListInterests(ctx)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	resp := &profile.ListInterestsResponse{}
	for _, interest := range interests {
		resp.Interests = append(resp.Interests, converter.InterestFromDomainToProto(interest))
	}

	return resp, nil
}

func (p *profileServer) ListInteractions(ctx context.Context, req *profile.ListInteractionsRequest) (*profile.ListInteractionsResponse, error) {
	interactions, err := p.service.ListInteractions(ctx, req.After.AsTime())
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	resp := &profile.ListInteractionsResponse{}
	for _, interaction := range interactions {
		resp.Interactions = append(resp.Interactions, converter.InteractionFromDomainToProto(interaction))
	}

	return resp, nil
}

func (p *profileServer) ListProfiles(ctx context.Context, req *profile.ListProfilesRequest) (*profile.ListProfilesResponse, error) {
	profs, err := p.service.ListProfiles(ctx, req.UpdatedAfter.AsTime())
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	resp := &profile.ListProfilesResponse{}
	for _, prof := range profs {
		resp.Profiles = append(resp.Profiles, converter.ProfileFromDomainToProto(prof))
	}

	return resp, nil
}

type photoServer struct {
	profile.UnimplementedPhotoServer
	service *service.PhotoService
}

func (p *photoServer) Create(ctx context.Context, req *profile.CreatePhotoRequest) (*profile.CreatePhotoResponse, error) {
	photo := converter.PhotoFromProtoToDomain(&profile.ProfilePhoto{
		Id:      uuid.New().String(),
		Content: req.Content,
	})

	id, err := p.service.Create(ctx, auf.ExtractUserIDFromContext[uuid.UUID](ctx), photo)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &profile.CreatePhotoResponse{
		Id: id.String(),
	}, nil
}

func (p *photoServer) Get(ctx context.Context, req *profile.GetPhotoRequest) (*profile.GetPhotoResponse, error) {
	userID, err := uuid.Parse(req.ProfileId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "profile id must have uuid format")
	}

	photoID, err := uuid.Parse(req.PhotoId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "photo id must have uuid format")
	}

	photo, err := p.service.Get(ctx, userID, photoID)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &profile.GetPhotoResponse{
		Photo: converter.PhotoFromDomainToProto(photo),
	}, nil
}

func (p *photoServer) List(ctx context.Context, req *profile.ListPhotoRequest) (*profile.ListPhotoResponse, error) {
	userID, err := uuid.Parse(req.ProfileId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "profile id must have uuid format")
	}

	photos, err := p.service.List(ctx, userID)
	if err == domain.ErrEntityIsNotExists {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	resp := &profile.ListPhotoResponse{}
	for _, photo := range photos {
		resp.Photos = append(resp.Photos, converter.PhotoFromDomainToProto(photo))
	}

	return resp, nil
}

func (p *photoServer) Delete(ctx context.Context, req *profile.DeletePhotoRequest) (*profile.DeletePhotoResponse, error) {
	photoID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "id must have uuid format")
	}

	if _, err := p.service.Delete(ctx, auf.ExtractUserIDFromContext[uuid.UUID](ctx), photoID); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &profile.DeletePhotoResponse{Id: photoID.String()}, nil
}

func (p *photoServer) Reorder(ctx context.Context, req *profile.ReorderPhotoRequest) (*profile.ReorderPhotoResponse, error) {
	var photoIDs []uuid.UUID
	for _, photoID := range req.Ids {
		photo, err := uuid.Parse(photoID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "photo id must have uuid format")
		}

		photoIDs = append(photoIDs, photo)
	}

	if err := p.service.Reorder(ctx, auf.ExtractUserIDFromContext[uuid.UUID](ctx), photoIDs); err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	return &profile.ReorderPhotoResponse{}, nil
}

func NewGRPCServer(profileService *service.ProfileService, photoService *service.PhotoService) *GRPCServer {
	srv := &GRPCServer{
		profileServer: &profileServer{service: profileService},
		photoServer:   &photoServer{service: photoService},
	}

	opts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_auth.UnaryServerInterceptor(nil),
			grpc_recovery.UnaryServerInterceptor(grut.WithLogAndRecover()),
		)),
	}
	srv.Server = grpc.NewServer(opts...)
	reflection.RegisterV1(srv.Server)
	profile.RegisterProfileServer(srv.Server, srv.profileServer)
	//profile.RegisterPhotoServer(srv.Server, srv.photoServer)

	return srv
}

func (g *GRPCServer) Shutdown(ctx context.Context) error {
	stopCh := make(chan struct{})

	go func() {
		g.Server.GracefulStop()
		close(stopCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stopCh:
		return nil
	}
}

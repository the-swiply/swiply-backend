// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: api/recommendation.proto

package recommendation

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RecommendationClient is the client API for Recommendation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RecommendationClient interface {
	GetRecommendations(ctx context.Context, in *GetRecommendationsRequest, opts ...grpc.CallOption) (*GetRecommendationsResponse, error)
}

type recommendationClient struct {
	cc grpc.ClientConnInterface
}

func NewRecommendationClient(cc grpc.ClientConnInterface) RecommendationClient {
	return &recommendationClient{cc}
}

func (c *recommendationClient) GetRecommendations(ctx context.Context, in *GetRecommendationsRequest, opts ...grpc.CallOption) (*GetRecommendationsResponse, error) {
	out := new(GetRecommendationsResponse)
	err := c.cc.Invoke(ctx, "/swiply.recommendation.Recommendation/GetRecommendations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecommendationServer is the server API for Recommendation service.
// All implementations must embed UnimplementedRecommendationServer
// for forward compatibility
type RecommendationServer interface {
	GetRecommendations(context.Context, *GetRecommendationsRequest) (*GetRecommendationsResponse, error)
	mustEmbedUnimplementedRecommendationServer()
}

// UnimplementedRecommendationServer must be embedded to have forward compatible implementations.
type UnimplementedRecommendationServer struct {
}

func (UnimplementedRecommendationServer) GetRecommendations(context.Context, *GetRecommendationsRequest) (*GetRecommendationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRecommendations not implemented")
}
func (UnimplementedRecommendationServer) mustEmbedUnimplementedRecommendationServer() {}

// UnsafeRecommendationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecommendationServer will
// result in compilation errors.
type UnsafeRecommendationServer interface {
	mustEmbedUnimplementedRecommendationServer()
}

func RegisterRecommendationServer(s grpc.ServiceRegistrar, srv RecommendationServer) {
	s.RegisterService(&Recommendation_ServiceDesc, srv)
}

func _Recommendation_GetRecommendations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRecommendationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecommendationServer).GetRecommendations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swiply.recommendation.Recommendation/GetRecommendations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecommendationServer).GetRecommendations(ctx, req.(*GetRecommendationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Recommendation_ServiceDesc is the grpc.ServiceDesc for Recommendation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Recommendation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "swiply.recommendation.Recommendation",
	HandlerType: (*RecommendationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRecommendations",
			Handler:    _Recommendation_GetRecommendations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/recommendation.proto",
}

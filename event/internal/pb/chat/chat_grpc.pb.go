// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: api/chat.proto

package chat

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

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatClient interface {
	SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error)
	GetNextMessages(ctx context.Context, in *GetNextMessagesRequest, opts ...grpc.CallOption) (*GetNextMessagesResponse, error)
	GetPreviousMessages(ctx context.Context, in *GetPreviousMessagesRequest, opts ...grpc.CallOption) (*GetPreviousMessagesResponse, error)
	GetChats(ctx context.Context, in *GetChatsRequest, opts ...grpc.CallOption) (*GetChatsResponse, error)
	LeaveChat(ctx context.Context, in *LeaveChatRequest, opts ...grpc.CallOption) (*LeaveChatResponse, error)
	CreateChat(ctx context.Context, in *CreateChatRequest, opts ...grpc.CallOption) (*CreateChatResponse, error)
	AddChatMembers(ctx context.Context, in *AddChatMembersRequest, opts ...grpc.CallOption) (*AddChatMembersResponse, error)
}

type chatClient struct {
	cc grpc.ClientConnInterface
}

func NewChatClient(cc grpc.ClientConnInterface) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) SendMessage(ctx context.Context, in *SendMessageRequest, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	out := new(SendMessageResponse)
	err := c.cc.Invoke(ctx, "/swiply.chat.Chat/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) GetNextMessages(ctx context.Context, in *GetNextMessagesRequest, opts ...grpc.CallOption) (*GetNextMessagesResponse, error) {
	out := new(GetNextMessagesResponse)
	err := c.cc.Invoke(ctx, "/swiply.chat.Chat/GetNextMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) GetPreviousMessages(ctx context.Context, in *GetPreviousMessagesRequest, opts ...grpc.CallOption) (*GetPreviousMessagesResponse, error) {
	out := new(GetPreviousMessagesResponse)
	err := c.cc.Invoke(ctx, "/swiply.chat.Chat/GetPreviousMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) GetChats(ctx context.Context, in *GetChatsRequest, opts ...grpc.CallOption) (*GetChatsResponse, error) {
	out := new(GetChatsResponse)
	err := c.cc.Invoke(ctx, "/swiply.chat.Chat/GetChats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) LeaveChat(ctx context.Context, in *LeaveChatRequest, opts ...grpc.CallOption) (*LeaveChatResponse, error) {
	out := new(LeaveChatResponse)
	err := c.cc.Invoke(ctx, "/swiply.chat.Chat/LeaveChat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) CreateChat(ctx context.Context, in *CreateChatRequest, opts ...grpc.CallOption) (*CreateChatResponse, error) {
	out := new(CreateChatResponse)
	err := c.cc.Invoke(ctx, "/swiply.chat.Chat/CreateChat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatClient) AddChatMembers(ctx context.Context, in *AddChatMembersRequest, opts ...grpc.CallOption) (*AddChatMembersResponse, error) {
	out := new(AddChatMembersResponse)
	err := c.cc.Invoke(ctx, "/swiply.chat.Chat/AddChatMembers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServer is the server API for Chat service.
// All implementations must embed UnimplementedChatServer
// for forward compatibility
type ChatServer interface {
	SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error)
	GetNextMessages(context.Context, *GetNextMessagesRequest) (*GetNextMessagesResponse, error)
	GetPreviousMessages(context.Context, *GetPreviousMessagesRequest) (*GetPreviousMessagesResponse, error)
	GetChats(context.Context, *GetChatsRequest) (*GetChatsResponse, error)
	LeaveChat(context.Context, *LeaveChatRequest) (*LeaveChatResponse, error)
	CreateChat(context.Context, *CreateChatRequest) (*CreateChatResponse, error)
	AddChatMembers(context.Context, *AddChatMembersRequest) (*AddChatMembersResponse, error)
	mustEmbedUnimplementedChatServer()
}

// UnimplementedChatServer must be embedded to have forward compatible implementations.
type UnimplementedChatServer struct {
}

func (UnimplementedChatServer) SendMessage(context.Context, *SendMessageRequest) (*SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedChatServer) GetNextMessages(context.Context, *GetNextMessagesRequest) (*GetNextMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNextMessages not implemented")
}
func (UnimplementedChatServer) GetPreviousMessages(context.Context, *GetPreviousMessagesRequest) (*GetPreviousMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPreviousMessages not implemented")
}
func (UnimplementedChatServer) GetChats(context.Context, *GetChatsRequest) (*GetChatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChats not implemented")
}
func (UnimplementedChatServer) LeaveChat(context.Context, *LeaveChatRequest) (*LeaveChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveChat not implemented")
}
func (UnimplementedChatServer) CreateChat(context.Context, *CreateChatRequest) (*CreateChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateChat not implemented")
}
func (UnimplementedChatServer) AddChatMembers(context.Context, *AddChatMembersRequest) (*AddChatMembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddChatMembers not implemented")
}
func (UnimplementedChatServer) mustEmbedUnimplementedChatServer() {}

// UnsafeChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServer will
// result in compilation errors.
type UnsafeChatServer interface {
	mustEmbedUnimplementedChatServer()
}

func RegisterChatServer(s grpc.ServiceRegistrar, srv ChatServer) {
	s.RegisterService(&Chat_ServiceDesc, srv)
}

func _Chat_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swiply.chat.Chat/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).SendMessage(ctx, req.(*SendMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_GetNextMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNextMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).GetNextMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swiply.chat.Chat/GetNextMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).GetNextMessages(ctx, req.(*GetNextMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_GetPreviousMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPreviousMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).GetPreviousMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swiply.chat.Chat/GetPreviousMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).GetPreviousMessages(ctx, req.(*GetPreviousMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_GetChats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).GetChats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swiply.chat.Chat/GetChats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).GetChats(ctx, req.(*GetChatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_LeaveChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).LeaveChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swiply.chat.Chat/LeaveChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).LeaveChat(ctx, req.(*LeaveChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_CreateChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).CreateChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swiply.chat.Chat/CreateChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).CreateChat(ctx, req.(*CreateChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Chat_AddChatMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddChatMembersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).AddChatMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swiply.chat.Chat/AddChatMembers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).AddChatMembers(ctx, req.(*AddChatMembersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Chat_ServiceDesc is the grpc.ServiceDesc for Chat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "swiply.chat.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _Chat_SendMessage_Handler,
		},
		{
			MethodName: "GetNextMessages",
			Handler:    _Chat_GetNextMessages_Handler,
		},
		{
			MethodName: "GetPreviousMessages",
			Handler:    _Chat_GetPreviousMessages_Handler,
		},
		{
			MethodName: "GetChats",
			Handler:    _Chat_GetChats_Handler,
		},
		{
			MethodName: "LeaveChat",
			Handler:    _Chat_LeaveChat_Handler,
		},
		{
			MethodName: "CreateChat",
			Handler:    _Chat_CreateChat_Handler,
		},
		{
			MethodName: "AddChatMembers",
			Handler:    _Chat_AddChatMembers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/chat.proto",
}

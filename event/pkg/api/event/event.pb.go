// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: api/event.proto

package event

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EventModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId     int64                  `protobuf:"varint,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	Title       string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Photos      [][]byte               `protobuf:"bytes,4,rep,name=photos,proto3" json:"photos,omitempty"`
	Date        *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=date,proto3" json:"date,omitempty"`
}

func (x *EventModel) Reset() {
	*x = EventModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventModel) ProtoMessage() {}

func (x *EventModel) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventModel.ProtoReflect.Descriptor instead.
func (*EventModel) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{0}
}

func (x *EventModel) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

func (x *EventModel) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *EventModel) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *EventModel) GetPhotos() [][]byte {
	if x != nil {
		return x.Photos
	}
	return nil
}

func (x *EventModel) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

type CreateEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title       string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Photos      [][]byte               `protobuf:"bytes,3,rep,name=photos,proto3" json:"photos,omitempty"`
	Date        *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=date,proto3" json:"date,omitempty"`
}

func (x *CreateEventRequest) Reset() {
	*x = CreateEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateEventRequest) ProtoMessage() {}

func (x *CreateEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateEventRequest.ProtoReflect.Descriptor instead.
func (*CreateEventRequest) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{1}
}

func (x *CreateEventRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *CreateEventRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateEventRequest) GetPhotos() [][]byte {
	if x != nil {
		return x.Photos
	}
	return nil
}

func (x *CreateEventRequest) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

type CreateEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId int64 `protobuf:"varint,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
}

func (x *CreateEventResponse) Reset() {
	*x = CreateEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateEventResponse) ProtoMessage() {}

func (x *CreateEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateEventResponse.ProtoReflect.Descriptor instead.
func (*CreateEventResponse) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{2}
}

func (x *CreateEventResponse) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

type UpdateEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event *EventModel `protobuf:"bytes,1,opt,name=event,proto3" json:"event,omitempty"`
}

func (x *UpdateEventRequest) Reset() {
	*x = UpdateEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateEventRequest) ProtoMessage() {}

func (x *UpdateEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateEventRequest.ProtoReflect.Descriptor instead.
func (*UpdateEventRequest) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateEventRequest) GetEvent() *EventModel {
	if x != nil {
		return x.Event
	}
	return nil
}

type UpdateEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateEventResponse) Reset() {
	*x = UpdateEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateEventResponse) ProtoMessage() {}

func (x *UpdateEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateEventResponse.ProtoReflect.Descriptor instead.
func (*UpdateEventResponse) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{4}
}

type GetEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Limit  int64 `protobuf:"varint,1,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset int64 `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *GetEventsRequest) Reset() {
	*x = GetEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEventsRequest) ProtoMessage() {}

func (x *GetEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEventsRequest.ProtoReflect.Descriptor instead.
func (*GetEventsRequest) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{5}
}

func (x *GetEventsRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *GetEventsRequest) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type GetEventsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event []*EventModel `protobuf:"bytes,1,rep,name=event,proto3" json:"event,omitempty"`
}

func (x *GetEventsResponse) Reset() {
	*x = GetEventsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEventsResponse) ProtoMessage() {}

func (x *GetEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEventsResponse.ProtoReflect.Descriptor instead.
func (*GetEventsResponse) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{6}
}

func (x *GetEventsResponse) GetEvent() []*EventModel {
	if x != nil {
		return x.Event
	}
	return nil
}

type GetUserEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetUserEventsRequest) Reset() {
	*x = GetUserEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserEventsRequest) ProtoMessage() {}

func (x *GetUserEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserEventsRequest.ProtoReflect.Descriptor instead.
func (*GetUserEventsRequest) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{7}
}

type GetUserEventsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event []*EventModel `protobuf:"bytes,1,rep,name=event,proto3" json:"event,omitempty"`
}

func (x *GetUserEventsResponse) Reset() {
	*x = GetUserEventsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserEventsResponse) ProtoMessage() {}

func (x *GetUserEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserEventsResponse.ProtoReflect.Descriptor instead.
func (*GetUserEventsResponse) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{8}
}

func (x *GetUserEventsResponse) GetEvent() []*EventModel {
	if x != nil {
		return x.Event
	}
	return nil
}

type JoinEventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId int64 `protobuf:"varint,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
}

func (x *JoinEventRequest) Reset() {
	*x = JoinEventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinEventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinEventRequest) ProtoMessage() {}

func (x *JoinEventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinEventRequest.ProtoReflect.Descriptor instead.
func (*JoinEventRequest) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{9}
}

func (x *JoinEventRequest) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

type JoinEventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *JoinEventResponse) Reset() {
	*x = JoinEventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinEventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinEventResponse) ProtoMessage() {}

func (x *JoinEventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinEventResponse.ProtoReflect.Descriptor instead.
func (*JoinEventResponse) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{10}
}

type AcceptEventJoinRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventId int64  `protobuf:"varint,1,opt,name=event_id,json=eventId,proto3" json:"event_id,omitempty"`
	UserId  string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *AcceptEventJoinRequest) Reset() {
	*x = AcceptEventJoinRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptEventJoinRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptEventJoinRequest) ProtoMessage() {}

func (x *AcceptEventJoinRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptEventJoinRequest.ProtoReflect.Descriptor instead.
func (*AcceptEventJoinRequest) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{11}
}

func (x *AcceptEventJoinRequest) GetEventId() int64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

func (x *AcceptEventJoinRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type AcceptEventJoinResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AcceptEventJoinResponse) Reset() {
	*x = AcceptEventJoinResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_event_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptEventJoinResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptEventJoinResponse) ProtoMessage() {}

func (x *AcceptEventJoinResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_event_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptEventJoinResponse.ProtoReflect.Descriptor instead.
func (*AcceptEventJoinResponse) Descriptor() ([]byte, []int) {
	return file_api_event_proto_rawDescGZIP(), []int{12}
}

var File_api_event_proto protoreflect.FileDescriptor

var file_api_event_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0c, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa7,
	0x01, 0x0a, 0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x19, 0x0a,
	0x08, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x0a, 0x06, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0c,
	0x52, 0x06, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x12, 0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x22, 0x94, 0x01, 0x0a, 0x12, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x68, 0x6f, 0x74, 0x6f,
	0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x73, 0x12,
	0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x22,
	0x30, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x22, 0x44, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x6c,
	0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x15, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x40,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74,
	0x22, 0x43, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52, 0x05,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x16, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x47, 0x0a,
	0x15, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52,
	0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x2d, 0x0a, 0x10, 0x4a, 0x6f, 0x69, 0x6e, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x13, 0x0a, 0x11, 0x4a, 0x6f, 0x69, 0x6e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x4c, 0x0a, 0x16, 0x41, 0x63,
	0x63, 0x65, 0x70, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x19, 0x0a, 0x17, 0x41, 0x63, 0x63, 0x65,
	0x70, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0xb8, 0x05, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x6f, 0x0a,
	0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x2e, 0x73,
	0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21,
	0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x3a, 0x01, 0x2a, 0x22, 0x10, 0x2f, 0x76,
	0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x6f,
	0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x2e,
	0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x21, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x3a, 0x01, 0x2a, 0x22, 0x10, 0x2f,
	0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12,
	0x60, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x1e, 0x2e, 0x73,
	0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x73,
	0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x12, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x12, 0x6f, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x12, 0x22, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x0f, 0x12, 0x0d, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f,
	0x6d, 0x79, 0x12, 0x6f, 0x0a, 0x09, 0x4a, 0x6f, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12,
	0x1e, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x4a,
	0x6f, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1f, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x4a,
	0x6f, 0x69, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x21, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x22, 0x19, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x2f, 0x6a, 0x6f, 0x69, 0x6e, 0x2f, 0x7b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f,
	0x69, 0x64, 0x7d, 0x12, 0x88, 0x01, 0x0a, 0x0f, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x4a, 0x6f, 0x69, 0x6e, 0x12, 0x24, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79,
	0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e,
	0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63,
	0x65, 0x70, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x28, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x22, 0x22, 0x20, 0x2f, 0x76,
	0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x2d, 0x6a,
	0x6f, 0x69, 0x6e, 0x2f, 0x7b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x7d, 0x42, 0x2b,
	0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x65,
	0x2d, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x6b,
	0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_api_event_proto_rawDescOnce sync.Once
	file_api_event_proto_rawDescData = file_api_event_proto_rawDesc
)

func file_api_event_proto_rawDescGZIP() []byte {
	file_api_event_proto_rawDescOnce.Do(func() {
		file_api_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_event_proto_rawDescData)
	})
	return file_api_event_proto_rawDescData
}

var file_api_event_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_api_event_proto_goTypes = []interface{}{
	(*EventModel)(nil),              // 0: swiply.event.EventModel
	(*CreateEventRequest)(nil),      // 1: swiply.event.CreateEventRequest
	(*CreateEventResponse)(nil),     // 2: swiply.event.CreateEventResponse
	(*UpdateEventRequest)(nil),      // 3: swiply.event.UpdateEventRequest
	(*UpdateEventResponse)(nil),     // 4: swiply.event.UpdateEventResponse
	(*GetEventsRequest)(nil),        // 5: swiply.event.GetEventsRequest
	(*GetEventsResponse)(nil),       // 6: swiply.event.GetEventsResponse
	(*GetUserEventsRequest)(nil),    // 7: swiply.event.GetUserEventsRequest
	(*GetUserEventsResponse)(nil),   // 8: swiply.event.GetUserEventsResponse
	(*JoinEventRequest)(nil),        // 9: swiply.event.JoinEventRequest
	(*JoinEventResponse)(nil),       // 10: swiply.event.JoinEventResponse
	(*AcceptEventJoinRequest)(nil),  // 11: swiply.event.AcceptEventJoinRequest
	(*AcceptEventJoinResponse)(nil), // 12: swiply.event.AcceptEventJoinResponse
	(*timestamppb.Timestamp)(nil),   // 13: google.protobuf.Timestamp
}
var file_api_event_proto_depIdxs = []int32{
	13, // 0: swiply.event.EventModel.date:type_name -> google.protobuf.Timestamp
	13, // 1: swiply.event.CreateEventRequest.date:type_name -> google.protobuf.Timestamp
	0,  // 2: swiply.event.UpdateEventRequest.event:type_name -> swiply.event.EventModel
	0,  // 3: swiply.event.GetEventsResponse.event:type_name -> swiply.event.EventModel
	0,  // 4: swiply.event.GetUserEventsResponse.event:type_name -> swiply.event.EventModel
	1,  // 5: swiply.event.Event.CreateEvent:input_type -> swiply.event.CreateEventRequest
	3,  // 6: swiply.event.Event.UpdateEvent:input_type -> swiply.event.UpdateEventRequest
	5,  // 7: swiply.event.Event.GetEvents:input_type -> swiply.event.GetEventsRequest
	7,  // 8: swiply.event.Event.GetUserEvents:input_type -> swiply.event.GetUserEventsRequest
	9,  // 9: swiply.event.Event.JoinEvent:input_type -> swiply.event.JoinEventRequest
	11, // 10: swiply.event.Event.AcceptEventJoin:input_type -> swiply.event.AcceptEventJoinRequest
	2,  // 11: swiply.event.Event.CreateEvent:output_type -> swiply.event.CreateEventResponse
	4,  // 12: swiply.event.Event.UpdateEvent:output_type -> swiply.event.UpdateEventResponse
	6,  // 13: swiply.event.Event.GetEvents:output_type -> swiply.event.GetEventsResponse
	8,  // 14: swiply.event.Event.GetUserEvents:output_type -> swiply.event.GetUserEventsResponse
	10, // 15: swiply.event.Event.JoinEvent:output_type -> swiply.event.JoinEventResponse
	12, // 16: swiply.event.Event.AcceptEventJoin:output_type -> swiply.event.AcceptEventJoinResponse
	11, // [11:17] is the sub-list for method output_type
	5,  // [5:11] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_api_event_proto_init() }
func file_api_event_proto_init() {
	if File_api_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventModel); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateEventRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateEventResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateEventRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateEventResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEventsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetEventsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserEventsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserEventsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinEventRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinEventResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptEventJoinRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_event_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptEventJoinResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_event_proto_goTypes,
		DependencyIndexes: file_api_event_proto_depIdxs,
		MessageInfos:      file_api_event_proto_msgTypes,
	}.Build()
	File_api_event_proto = out.File
	file_api_event_proto_rawDesc = nil
	file_api_event_proto_goTypes = nil
	file_api_event_proto_depIdxs = nil
}
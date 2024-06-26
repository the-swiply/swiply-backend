// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: api/randomcoffee.proto

package randomcoffee

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

type MeetingStatus int32

const (
	MeetingStatus_MEETING_STATUS_UNSPECIFIED MeetingStatus = 0
	MeetingStatus_AWAITING_SCHEDULE          MeetingStatus = 1
	MeetingStatus_SCHEDULING                 MeetingStatus = 2
	MeetingStatus_SCHEDULED                  MeetingStatus = 3
)

// Enum value maps for MeetingStatus.
var (
	MeetingStatus_name = map[int32]string{
		0: "MEETING_STATUS_UNSPECIFIED",
		1: "AWAITING_SCHEDULE",
		2: "SCHEDULING",
		3: "SCHEDULED",
	}
	MeetingStatus_value = map[string]int32{
		"MEETING_STATUS_UNSPECIFIED": 0,
		"AWAITING_SCHEDULE":          1,
		"SCHEDULING":                 2,
		"SCHEDULED":                  3,
	}
)

func (x MeetingStatus) Enum() *MeetingStatus {
	p := new(MeetingStatus)
	*p = x
	return p
}

func (x MeetingStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MeetingStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_randomcoffee_proto_enumTypes[0].Descriptor()
}

func (MeetingStatus) Type() protoreflect.EnumType {
	return &file_api_randomcoffee_proto_enumTypes[0]
}

func (x MeetingStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MeetingStatus.Descriptor instead.
func (MeetingStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{0}
}

type Meeting struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	OwnerId        string                 `protobuf:"bytes,2,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
	MemberId       string                 `protobuf:"bytes,3,opt,name=member_id,json=memberId,proto3" json:"member_id,omitempty"`
	Start          *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=start,proto3" json:"start,omitempty"`
	End            *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=end,proto3" json:"end,omitempty"`
	OrganizationId int64                  `protobuf:"varint,6,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
	Status         MeetingStatus          `protobuf:"varint,7,opt,name=status,proto3,enum=swiply.randomcoffee.MeetingStatus" json:"status,omitempty"`
	CreatedAt      *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
}

func (x *Meeting) Reset() {
	*x = Meeting{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_randomcoffee_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Meeting) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Meeting) ProtoMessage() {}

func (x *Meeting) ProtoReflect() protoreflect.Message {
	mi := &file_api_randomcoffee_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Meeting.ProtoReflect.Descriptor instead.
func (*Meeting) Descriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{0}
}

func (x *Meeting) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Meeting) GetOwnerId() string {
	if x != nil {
		return x.OwnerId
	}
	return ""
}

func (x *Meeting) GetMemberId() string {
	if x != nil {
		return x.MemberId
	}
	return ""
}

func (x *Meeting) GetStart() *timestamppb.Timestamp {
	if x != nil {
		return x.Start
	}
	return nil
}

func (x *Meeting) GetEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.End
	}
	return nil
}

func (x *Meeting) GetOrganizationId() int64 {
	if x != nil {
		return x.OrganizationId
	}
	return 0
}

func (x *Meeting) GetStatus() MeetingStatus {
	if x != nil {
		return x.Status
	}
	return MeetingStatus_MEETING_STATUS_UNSPECIFIED
}

func (x *Meeting) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type CreateMeetingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start          *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	End            *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=end,proto3" json:"end,omitempty"`
	OrganizationId int64                  `protobuf:"varint,3,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
}

func (x *CreateMeetingRequest) Reset() {
	*x = CreateMeetingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_randomcoffee_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateMeetingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMeetingRequest) ProtoMessage() {}

func (x *CreateMeetingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_randomcoffee_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMeetingRequest.ProtoReflect.Descriptor instead.
func (*CreateMeetingRequest) Descriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{1}
}

func (x *CreateMeetingRequest) GetStart() *timestamppb.Timestamp {
	if x != nil {
		return x.Start
	}
	return nil
}

func (x *CreateMeetingRequest) GetEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.End
	}
	return nil
}

func (x *CreateMeetingRequest) GetOrganizationId() int64 {
	if x != nil {
		return x.OrganizationId
	}
	return 0
}

type CreateMeetingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Meeting *Meeting `protobuf:"bytes,1,opt,name=meeting,proto3" json:"meeting,omitempty"`
}

func (x *CreateMeetingResponse) Reset() {
	*x = CreateMeetingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_randomcoffee_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateMeetingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMeetingResponse) ProtoMessage() {}

func (x *CreateMeetingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_randomcoffee_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMeetingResponse.ProtoReflect.Descriptor instead.
func (*CreateMeetingResponse) Descriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{2}
}

func (x *CreateMeetingResponse) GetMeeting() *Meeting {
	if x != nil {
		return x.Meeting
	}
	return nil
}

type DeleteMeetingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteMeetingRequest) Reset() {
	*x = DeleteMeetingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_randomcoffee_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteMeetingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteMeetingRequest) ProtoMessage() {}

func (x *DeleteMeetingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_randomcoffee_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteMeetingRequest.ProtoReflect.Descriptor instead.
func (*DeleteMeetingRequest) Descriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteMeetingRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteMeetingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteMeetingResponse) Reset() {
	*x = DeleteMeetingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_randomcoffee_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteMeetingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteMeetingResponse) ProtoMessage() {}

func (x *DeleteMeetingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_randomcoffee_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteMeetingResponse.ProtoReflect.Descriptor instead.
func (*DeleteMeetingResponse) Descriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{4}
}

type UpdateMeetingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Start          *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=start,proto3" json:"start,omitempty"`
	End            *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=end,proto3" json:"end,omitempty"`
	OrganizationId int64                  `protobuf:"varint,4,opt,name=organization_id,json=organizationId,proto3" json:"organization_id,omitempty"`
}

func (x *UpdateMeetingRequest) Reset() {
	*x = UpdateMeetingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_randomcoffee_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMeetingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMeetingRequest) ProtoMessage() {}

func (x *UpdateMeetingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_randomcoffee_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMeetingRequest.ProtoReflect.Descriptor instead.
func (*UpdateMeetingRequest) Descriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateMeetingRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateMeetingRequest) GetStart() *timestamppb.Timestamp {
	if x != nil {
		return x.Start
	}
	return nil
}

func (x *UpdateMeetingRequest) GetEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.End
	}
	return nil
}

func (x *UpdateMeetingRequest) GetOrganizationId() int64 {
	if x != nil {
		return x.OrganizationId
	}
	return 0
}

type UpdateMeetingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateMeetingResponse) Reset() {
	*x = UpdateMeetingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_randomcoffee_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMeetingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMeetingResponse) ProtoMessage() {}

func (x *UpdateMeetingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_randomcoffee_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMeetingResponse.ProtoReflect.Descriptor instead.
func (*UpdateMeetingResponse) Descriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{6}
}

type ListMeetingsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListMeetingsRequest) Reset() {
	*x = ListMeetingsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_randomcoffee_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMeetingsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMeetingsRequest) ProtoMessage() {}

func (x *ListMeetingsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_randomcoffee_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMeetingsRequest.ProtoReflect.Descriptor instead.
func (*ListMeetingsRequest) Descriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{7}
}

type ListMeetingsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Meetings []*Meeting `protobuf:"bytes,1,rep,name=meetings,proto3" json:"meetings,omitempty"`
}

func (x *ListMeetingsResponse) Reset() {
	*x = ListMeetingsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_randomcoffee_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMeetingsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMeetingsResponse) ProtoMessage() {}

func (x *ListMeetingsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_randomcoffee_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMeetingsResponse.ProtoReflect.Descriptor instead.
func (*ListMeetingsResponse) Descriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{8}
}

func (x *ListMeetingsResponse) GetMeetings() []*Meeting {
	if x != nil {
		return x.Meetings
	}
	return nil
}

type GetMeetingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetMeetingRequest) Reset() {
	*x = GetMeetingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_randomcoffee_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMeetingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMeetingRequest) ProtoMessage() {}

func (x *GetMeetingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_randomcoffee_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMeetingRequest.ProtoReflect.Descriptor instead.
func (*GetMeetingRequest) Descriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{9}
}

func (x *GetMeetingRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetMeetingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Meeting *Meeting `protobuf:"bytes,1,opt,name=meeting,proto3" json:"meeting,omitempty"`
}

func (x *GetMeetingResponse) Reset() {
	*x = GetMeetingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_randomcoffee_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMeetingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMeetingResponse) ProtoMessage() {}

func (x *GetMeetingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_randomcoffee_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMeetingResponse.ProtoReflect.Descriptor instead.
func (*GetMeetingResponse) Descriptor() ([]byte, []int) {
	return file_api_randomcoffee_proto_rawDescGZIP(), []int{10}
}

func (x *GetMeetingResponse) GetMeeting() *Meeting {
	if x != nil {
		return x.Meeting
	}
	return nil
}

var File_api_randomcoffee_proto protoreflect.FileDescriptor

var file_api_randomcoffee_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66,
	0x65, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79,
	0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66, 0x65, 0x65, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd0, 0x02, 0x0a,
	0x07, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x30, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x05, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x12, 0x2c, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x03, 0x65, 0x6e, 0x64,
	0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x6f, 0x72, 0x67, 0x61, 0x6e,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x3a, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x73, 0x77, 0x69, 0x70,
	0x6c, 0x79, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66, 0x65, 0x65, 0x2e,
	0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x38, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22,
	0x9f, 0x01, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x30, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x2c, 0x0a, 0x03, 0x65, 0x6e,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61,
	0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0e, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x22, 0x4f, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x65, 0x74, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x6d, 0x65,
	0x65, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x77,
	0x69, 0x70, 0x6c, 0x79, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66, 0x65,
	0x65, 0x2e, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x07, 0x6d, 0x65, 0x65, 0x74, 0x69,
	0x6e, 0x67, 0x22, 0x26, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x65, 0x65, 0x74,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x17, 0x0a, 0x15, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0xaf, 0x01, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65,
	0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x05,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x2c,
	0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x12, 0x27, 0x0a, 0x0f,
	0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x17, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d,
	0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15,
	0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x50, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x65,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a,
	0x08, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63,
	0x6f, 0x66, 0x66, 0x65, 0x65, 0x2e, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x08, 0x6d,
	0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x22, 0x23, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x4d, 0x65,
	0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x4c, 0x0a, 0x12,
	0x47, 0x65, 0x74, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x36, 0x0a, 0x07, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x72, 0x61, 0x6e,
	0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66, 0x65, 0x65, 0x2e, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e,
	0x67, 0x52, 0x07, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x2a, 0x65, 0x0a, 0x0d, 0x4d, 0x65,
	0x65, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a, 0x1a, 0x4d,
	0x45, 0x45, 0x54, 0x49, 0x4e, 0x47, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x41,
	0x57, 0x41, 0x49, 0x54, 0x49, 0x4e, 0x47, 0x5f, 0x53, 0x43, 0x48, 0x45, 0x44, 0x55, 0x4c, 0x45,
	0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x43, 0x48, 0x45, 0x44, 0x55, 0x4c, 0x49, 0x4e, 0x47,
	0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x43, 0x48, 0x45, 0x44, 0x55, 0x4c, 0x45, 0x44, 0x10,
	0x03, 0x32, 0xf3, 0x04, 0x0a, 0x0c, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x43, 0x6f, 0x66, 0x66,
	0x65, 0x65, 0x12, 0x7e, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x29, 0x2e, 0x73,
	0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66,
	0x65, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79,
	0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66, 0x65, 0x65, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x3a, 0x01, 0x2a, 0x22, 0x12,
	0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x12, 0x7e, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x29, 0x2e, 0x73,
	0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66,
	0x65, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79,
	0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66, 0x65, 0x65, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x3a, 0x01, 0x2a, 0x22, 0x12,
	0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x12, 0x7e, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x29, 0x2e, 0x73,
	0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66,
	0x65, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79,
	0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66, 0x65, 0x65, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x17, 0x3a, 0x01, 0x2a, 0x22, 0x12,
	0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x12, 0x71, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x28, 0x2e, 0x73, 0x77, 0x69,
	0x70, 0x6c, 0x79, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66, 0x65, 0x65,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x72, 0x61,
	0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66, 0x65, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d,
	0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x12, 0x0c, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x65,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x70, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x26, 0x2e, 0x73,
	0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66,
	0x65, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79, 0x2e, 0x72, 0x61,
	0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66, 0x65, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65,
	0x65, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x18, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x12, 0x12, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x65, 0x74, 0x69,
	0x6e, 0x67, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x68, 0x65, 0x2d, 0x73, 0x77, 0x69, 0x70, 0x6c, 0x79,
	0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x63, 0x6f, 0x66, 0x66, 0x65, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_randomcoffee_proto_rawDescOnce sync.Once
	file_api_randomcoffee_proto_rawDescData = file_api_randomcoffee_proto_rawDesc
)

func file_api_randomcoffee_proto_rawDescGZIP() []byte {
	file_api_randomcoffee_proto_rawDescOnce.Do(func() {
		file_api_randomcoffee_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_randomcoffee_proto_rawDescData)
	})
	return file_api_randomcoffee_proto_rawDescData
}

var file_api_randomcoffee_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_randomcoffee_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_api_randomcoffee_proto_goTypes = []interface{}{
	(MeetingStatus)(0),            // 0: swiply.randomcoffee.MeetingStatus
	(*Meeting)(nil),               // 1: swiply.randomcoffee.Meeting
	(*CreateMeetingRequest)(nil),  // 2: swiply.randomcoffee.CreateMeetingRequest
	(*CreateMeetingResponse)(nil), // 3: swiply.randomcoffee.CreateMeetingResponse
	(*DeleteMeetingRequest)(nil),  // 4: swiply.randomcoffee.DeleteMeetingRequest
	(*DeleteMeetingResponse)(nil), // 5: swiply.randomcoffee.DeleteMeetingResponse
	(*UpdateMeetingRequest)(nil),  // 6: swiply.randomcoffee.UpdateMeetingRequest
	(*UpdateMeetingResponse)(nil), // 7: swiply.randomcoffee.UpdateMeetingResponse
	(*ListMeetingsRequest)(nil),   // 8: swiply.randomcoffee.ListMeetingsRequest
	(*ListMeetingsResponse)(nil),  // 9: swiply.randomcoffee.ListMeetingsResponse
	(*GetMeetingRequest)(nil),     // 10: swiply.randomcoffee.GetMeetingRequest
	(*GetMeetingResponse)(nil),    // 11: swiply.randomcoffee.GetMeetingResponse
	(*timestamppb.Timestamp)(nil), // 12: google.protobuf.Timestamp
}
var file_api_randomcoffee_proto_depIdxs = []int32{
	12, // 0: swiply.randomcoffee.Meeting.start:type_name -> google.protobuf.Timestamp
	12, // 1: swiply.randomcoffee.Meeting.end:type_name -> google.protobuf.Timestamp
	0,  // 2: swiply.randomcoffee.Meeting.status:type_name -> swiply.randomcoffee.MeetingStatus
	12, // 3: swiply.randomcoffee.Meeting.CreatedAt:type_name -> google.protobuf.Timestamp
	12, // 4: swiply.randomcoffee.CreateMeetingRequest.start:type_name -> google.protobuf.Timestamp
	12, // 5: swiply.randomcoffee.CreateMeetingRequest.end:type_name -> google.protobuf.Timestamp
	1,  // 6: swiply.randomcoffee.CreateMeetingResponse.meeting:type_name -> swiply.randomcoffee.Meeting
	12, // 7: swiply.randomcoffee.UpdateMeetingRequest.start:type_name -> google.protobuf.Timestamp
	12, // 8: swiply.randomcoffee.UpdateMeetingRequest.end:type_name -> google.protobuf.Timestamp
	1,  // 9: swiply.randomcoffee.ListMeetingsResponse.meetings:type_name -> swiply.randomcoffee.Meeting
	1,  // 10: swiply.randomcoffee.GetMeetingResponse.meeting:type_name -> swiply.randomcoffee.Meeting
	2,  // 11: swiply.randomcoffee.RandomCoffee.Create:input_type -> swiply.randomcoffee.CreateMeetingRequest
	4,  // 12: swiply.randomcoffee.RandomCoffee.Delete:input_type -> swiply.randomcoffee.DeleteMeetingRequest
	6,  // 13: swiply.randomcoffee.RandomCoffee.Update:input_type -> swiply.randomcoffee.UpdateMeetingRequest
	8,  // 14: swiply.randomcoffee.RandomCoffee.List:input_type -> swiply.randomcoffee.ListMeetingsRequest
	10, // 15: swiply.randomcoffee.RandomCoffee.Get:input_type -> swiply.randomcoffee.GetMeetingRequest
	3,  // 16: swiply.randomcoffee.RandomCoffee.Create:output_type -> swiply.randomcoffee.CreateMeetingResponse
	5,  // 17: swiply.randomcoffee.RandomCoffee.Delete:output_type -> swiply.randomcoffee.DeleteMeetingResponse
	7,  // 18: swiply.randomcoffee.RandomCoffee.Update:output_type -> swiply.randomcoffee.UpdateMeetingResponse
	9,  // 19: swiply.randomcoffee.RandomCoffee.List:output_type -> swiply.randomcoffee.ListMeetingsResponse
	11, // 20: swiply.randomcoffee.RandomCoffee.Get:output_type -> swiply.randomcoffee.GetMeetingResponse
	16, // [16:21] is the sub-list for method output_type
	11, // [11:16] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_api_randomcoffee_proto_init() }
func file_api_randomcoffee_proto_init() {
	if File_api_randomcoffee_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_randomcoffee_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Meeting); i {
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
		file_api_randomcoffee_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateMeetingRequest); i {
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
		file_api_randomcoffee_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateMeetingResponse); i {
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
		file_api_randomcoffee_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteMeetingRequest); i {
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
		file_api_randomcoffee_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteMeetingResponse); i {
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
		file_api_randomcoffee_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMeetingRequest); i {
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
		file_api_randomcoffee_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMeetingResponse); i {
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
		file_api_randomcoffee_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListMeetingsRequest); i {
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
		file_api_randomcoffee_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListMeetingsResponse); i {
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
		file_api_randomcoffee_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMeetingRequest); i {
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
		file_api_randomcoffee_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMeetingResponse); i {
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
			RawDescriptor: file_api_randomcoffee_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_randomcoffee_proto_goTypes,
		DependencyIndexes: file_api_randomcoffee_proto_depIdxs,
		EnumInfos:         file_api_randomcoffee_proto_enumTypes,
		MessageInfos:      file_api_randomcoffee_proto_msgTypes,
	}.Build()
	File_api_randomcoffee_proto = out.File
	file_api_randomcoffee_proto_rawDesc = nil
	file_api_randomcoffee_proto_goTypes = nil
	file_api_randomcoffee_proto_depIdxs = nil
}

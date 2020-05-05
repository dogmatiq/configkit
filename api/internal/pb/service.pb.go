// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.11.4
// source: api/internal/pb/service.proto

package pb

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ListApplicationIdentitiesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListApplicationIdentitiesRequest) Reset() {
	*x = ListApplicationIdentitiesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_internal_pb_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListApplicationIdentitiesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListApplicationIdentitiesRequest) ProtoMessage() {}

func (x *ListApplicationIdentitiesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_internal_pb_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListApplicationIdentitiesRequest.ProtoReflect.Descriptor instead.
func (*ListApplicationIdentitiesRequest) Descriptor() ([]byte, []int) {
	return file_api_internal_pb_service_proto_rawDescGZIP(), []int{0}
}

type ListApplicationIdentitiesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Identities []*Identity `protobuf:"bytes,1,rep,name=identities,proto3" json:"identities,omitempty"`
}

func (x *ListApplicationIdentitiesResponse) Reset() {
	*x = ListApplicationIdentitiesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_internal_pb_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListApplicationIdentitiesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListApplicationIdentitiesResponse) ProtoMessage() {}

func (x *ListApplicationIdentitiesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_internal_pb_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListApplicationIdentitiesResponse.ProtoReflect.Descriptor instead.
func (*ListApplicationIdentitiesResponse) Descriptor() ([]byte, []int) {
	return file_api_internal_pb_service_proto_rawDescGZIP(), []int{1}
}

func (x *ListApplicationIdentitiesResponse) GetIdentities() []*Identity {
	if x != nil {
		return x.Identities
	}
	return nil
}

type ListApplicationsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListApplicationsRequest) Reset() {
	*x = ListApplicationsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_internal_pb_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListApplicationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListApplicationsRequest) ProtoMessage() {}

func (x *ListApplicationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_internal_pb_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListApplicationsRequest.ProtoReflect.Descriptor instead.
func (*ListApplicationsRequest) Descriptor() ([]byte, []int) {
	return file_api_internal_pb_service_proto_rawDescGZIP(), []int{2}
}

type ListApplicationsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Applications []*Application `protobuf:"bytes,1,rep,name=applications,proto3" json:"applications,omitempty"`
}

func (x *ListApplicationsResponse) Reset() {
	*x = ListApplicationsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_internal_pb_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListApplicationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListApplicationsResponse) ProtoMessage() {}

func (x *ListApplicationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_internal_pb_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListApplicationsResponse.ProtoReflect.Descriptor instead.
func (*ListApplicationsResponse) Descriptor() ([]byte, []int) {
	return file_api_internal_pb_service_proto_rawDescGZIP(), []int{3}
}

func (x *ListApplicationsResponse) GetApplications() []*Application {
	if x != nil {
		return x.Applications
	}
	return nil
}

type WatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *WatchRequest) Reset() {
	*x = WatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_internal_pb_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WatchRequest) ProtoMessage() {}

func (x *WatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_internal_pb_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WatchRequest.ProtoReflect.Descriptor instead.
func (*WatchRequest) Descriptor() ([]byte, []int) {
	return file_api_internal_pb_service_proto_rawDescGZIP(), []int{4}
}

type WatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *WatchResponse) Reset() {
	*x = WatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_internal_pb_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WatchResponse) ProtoMessage() {}

func (x *WatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_internal_pb_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WatchResponse.ProtoReflect.Descriptor instead.
func (*WatchResponse) Descriptor() ([]byte, []int) {
	return file_api_internal_pb_service_proto_rawDescGZIP(), []int{5}
}

var File_api_internal_pb_service_proto protoreflect.FileDescriptor

var file_api_internal_pb_service_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x61, 0x70, 0x69, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70,
	0x62, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0f, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31,
	0x1a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70,
	0x62, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22, 0x0a,
	0x20, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x5e, 0x0a, 0x21, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x22, 0x19, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x5c, 0x0a, 0x18,
	0x4c, 0x69, 0x73, 0x74, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x0c, 0x61, 0x70, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31,
	0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x61, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x0e, 0x0a, 0x0c, 0x57, 0x61,
	0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x0f, 0x0a, 0x0d, 0x57, 0x61,
	0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xc0, 0x02, 0x0a, 0x06,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x82, 0x01, 0x0a, 0x19, 0x4c, 0x69, 0x73, 0x74, 0x41,
	0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x69, 0x65, 0x73, 0x12, 0x31, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x32, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x67, 0x0a, 0x10, 0x4c,
	0x69, 0x73, 0x74, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x28, 0x2e, 0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76,
	0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x64, 0x6f, 0x67, 0x6d,
	0x61, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x05, 0x57, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1d, 0x2e,
	0x64, 0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x64,
	0x6f, 0x67, 0x6d, 0x61, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x57,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x2f,
	0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x6f, 0x67,
	0x6d, 0x61, 0x74, 0x69, 0x71, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6b, 0x69, 0x74, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_internal_pb_service_proto_rawDescOnce sync.Once
	file_api_internal_pb_service_proto_rawDescData = file_api_internal_pb_service_proto_rawDesc
)

func file_api_internal_pb_service_proto_rawDescGZIP() []byte {
	file_api_internal_pb_service_proto_rawDescOnce.Do(func() {
		file_api_internal_pb_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_internal_pb_service_proto_rawDescData)
	})
	return file_api_internal_pb_service_proto_rawDescData
}

var file_api_internal_pb_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_internal_pb_service_proto_goTypes = []interface{}{
	(*ListApplicationIdentitiesRequest)(nil),  // 0: dogma.config.v1.ListApplicationIdentitiesRequest
	(*ListApplicationIdentitiesResponse)(nil), // 1: dogma.config.v1.ListApplicationIdentitiesResponse
	(*ListApplicationsRequest)(nil),           // 2: dogma.config.v1.ListApplicationsRequest
	(*ListApplicationsResponse)(nil),          // 3: dogma.config.v1.ListApplicationsResponse
	(*WatchRequest)(nil),                      // 4: dogma.config.v1.WatchRequest
	(*WatchResponse)(nil),                     // 5: dogma.config.v1.WatchResponse
	(*Identity)(nil),                          // 6: dogma.config.v1.Identity
	(*Application)(nil),                       // 7: dogma.config.v1.Application
}
var file_api_internal_pb_service_proto_depIdxs = []int32{
	6, // 0: dogma.config.v1.ListApplicationIdentitiesResponse.identities:type_name -> dogma.config.v1.Identity
	7, // 1: dogma.config.v1.ListApplicationsResponse.applications:type_name -> dogma.config.v1.Application
	0, // 2: dogma.config.v1.Config.ListApplicationIdentities:input_type -> dogma.config.v1.ListApplicationIdentitiesRequest
	2, // 3: dogma.config.v1.Config.ListApplications:input_type -> dogma.config.v1.ListApplicationsRequest
	4, // 4: dogma.config.v1.Config.Watch:input_type -> dogma.config.v1.WatchRequest
	1, // 5: dogma.config.v1.Config.ListApplicationIdentities:output_type -> dogma.config.v1.ListApplicationIdentitiesResponse
	3, // 6: dogma.config.v1.Config.ListApplications:output_type -> dogma.config.v1.ListApplicationsResponse
	5, // 7: dogma.config.v1.Config.Watch:output_type -> dogma.config.v1.WatchResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_internal_pb_service_proto_init() }
func file_api_internal_pb_service_proto_init() {
	if File_api_internal_pb_service_proto != nil {
		return
	}
	file_api_internal_pb_types_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_internal_pb_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListApplicationIdentitiesRequest); i {
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
		file_api_internal_pb_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListApplicationIdentitiesResponse); i {
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
		file_api_internal_pb_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListApplicationsRequest); i {
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
		file_api_internal_pb_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListApplicationsResponse); i {
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
		file_api_internal_pb_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WatchRequest); i {
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
		file_api_internal_pb_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WatchResponse); i {
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
			RawDescriptor: file_api_internal_pb_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_internal_pb_service_proto_goTypes,
		DependencyIndexes: file_api_internal_pb_service_proto_depIdxs,
		MessageInfos:      file_api_internal_pb_service_proto_msgTypes,
	}.Build()
	File_api_internal_pb_service_proto = out.File
	file_api_internal_pb_service_proto_rawDesc = nil
	file_api_internal_pb_service_proto_goTypes = nil
	file_api_internal_pb_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ConfigClient is the client API for Config service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConfigClient interface {
	// ListApplicationIdentities returns the identity of all applications.
	ListApplicationIdentities(ctx context.Context, in *ListApplicationIdentitiesRequest, opts ...grpc.CallOption) (*ListApplicationIdentitiesResponse, error)
	// ListApplications returns the full configuration of all applications.
	ListApplications(ctx context.Context, in *ListApplicationsRequest, opts ...grpc.CallOption) (*ListApplicationsResponse, error)
	// Watch opens a stream that sends no content, simply to detect when the
	// server goes away.
	Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (Config_WatchClient, error)
}

type configClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigClient(cc grpc.ClientConnInterface) ConfigClient {
	return &configClient{cc}
}

func (c *configClient) ListApplicationIdentities(ctx context.Context, in *ListApplicationIdentitiesRequest, opts ...grpc.CallOption) (*ListApplicationIdentitiesResponse, error) {
	out := new(ListApplicationIdentitiesResponse)
	err := c.cc.Invoke(ctx, "/dogma.config.v1.Config/ListApplicationIdentities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) ListApplications(ctx context.Context, in *ListApplicationsRequest, opts ...grpc.CallOption) (*ListApplicationsResponse, error) {
	out := new(ListApplicationsResponse)
	err := c.cc.Invoke(ctx, "/dogma.config.v1.Config/ListApplications", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (Config_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Config_serviceDesc.Streams[0], "/dogma.config.v1.Config/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &configWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Config_WatchClient interface {
	Recv() (*WatchResponse, error)
	grpc.ClientStream
}

type configWatchClient struct {
	grpc.ClientStream
}

func (x *configWatchClient) Recv() (*WatchResponse, error) {
	m := new(WatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConfigServer is the server API for Config service.
type ConfigServer interface {
	// ListApplicationIdentities returns the identity of all applications.
	ListApplicationIdentities(context.Context, *ListApplicationIdentitiesRequest) (*ListApplicationIdentitiesResponse, error)
	// ListApplications returns the full configuration of all applications.
	ListApplications(context.Context, *ListApplicationsRequest) (*ListApplicationsResponse, error)
	// Watch opens a stream that sends no content, simply to detect when the
	// server goes away.
	Watch(*WatchRequest, Config_WatchServer) error
}

// UnimplementedConfigServer can be embedded to have forward compatible implementations.
type UnimplementedConfigServer struct {
}

func (*UnimplementedConfigServer) ListApplicationIdentities(context.Context, *ListApplicationIdentitiesRequest) (*ListApplicationIdentitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListApplicationIdentities not implemented")
}
func (*UnimplementedConfigServer) ListApplications(context.Context, *ListApplicationsRequest) (*ListApplicationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListApplications not implemented")
}
func (*UnimplementedConfigServer) Watch(*WatchRequest, Config_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}

func RegisterConfigServer(s *grpc.Server, srv ConfigServer) {
	s.RegisterService(&_Config_serviceDesc, srv)
}

func _Config_ListApplicationIdentities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListApplicationIdentitiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).ListApplicationIdentities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dogma.config.v1.Config/ListApplicationIdentities",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).ListApplicationIdentities(ctx, req.(*ListApplicationIdentitiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_ListApplications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListApplicationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServer).ListApplications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dogma.config.v1.Config/ListApplications",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).ListApplications(ctx, req.(*ListApplicationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Config_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConfigServer).Watch(m, &configWatchServer{stream})
}

type Config_WatchServer interface {
	Send(*WatchResponse) error
	grpc.ServerStream
}

type configWatchServer struct {
	grpc.ServerStream
}

func (x *configWatchServer) Send(m *WatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _Config_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dogma.config.v1.Config",
	HandlerType: (*ConfigServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListApplicationIdentities",
			Handler:    _Config_ListApplicationIdentities_Handler,
		},
		{
			MethodName: "ListApplications",
			Handler:    _Config_ListApplications_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _Config_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/internal/pb/service.proto",
}

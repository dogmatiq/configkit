// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/internal/pb/service.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ListApplicationIdentitiesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListApplicationIdentitiesRequest) Reset()         { *m = ListApplicationIdentitiesRequest{} }
func (m *ListApplicationIdentitiesRequest) String() string { return proto.CompactTextString(m) }
func (*ListApplicationIdentitiesRequest) ProtoMessage()    {}
func (*ListApplicationIdentitiesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbf1bd0c00876e46, []int{0}
}

func (m *ListApplicationIdentitiesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListApplicationIdentitiesRequest.Unmarshal(m, b)
}
func (m *ListApplicationIdentitiesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListApplicationIdentitiesRequest.Marshal(b, m, deterministic)
}
func (m *ListApplicationIdentitiesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListApplicationIdentitiesRequest.Merge(m, src)
}
func (m *ListApplicationIdentitiesRequest) XXX_Size() int {
	return xxx_messageInfo_ListApplicationIdentitiesRequest.Size(m)
}
func (m *ListApplicationIdentitiesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListApplicationIdentitiesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListApplicationIdentitiesRequest proto.InternalMessageInfo

type ListApplicationIdentitiesResponse struct {
	Identities           []*Identity `protobuf:"bytes,1,rep,name=identities,proto3" json:"identities,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *ListApplicationIdentitiesResponse) Reset()         { *m = ListApplicationIdentitiesResponse{} }
func (m *ListApplicationIdentitiesResponse) String() string { return proto.CompactTextString(m) }
func (*ListApplicationIdentitiesResponse) ProtoMessage()    {}
func (*ListApplicationIdentitiesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbf1bd0c00876e46, []int{1}
}

func (m *ListApplicationIdentitiesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListApplicationIdentitiesResponse.Unmarshal(m, b)
}
func (m *ListApplicationIdentitiesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListApplicationIdentitiesResponse.Marshal(b, m, deterministic)
}
func (m *ListApplicationIdentitiesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListApplicationIdentitiesResponse.Merge(m, src)
}
func (m *ListApplicationIdentitiesResponse) XXX_Size() int {
	return xxx_messageInfo_ListApplicationIdentitiesResponse.Size(m)
}
func (m *ListApplicationIdentitiesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListApplicationIdentitiesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListApplicationIdentitiesResponse proto.InternalMessageInfo

func (m *ListApplicationIdentitiesResponse) GetIdentities() []*Identity {
	if m != nil {
		return m.Identities
	}
	return nil
}

type ListApplicationsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListApplicationsRequest) Reset()         { *m = ListApplicationsRequest{} }
func (m *ListApplicationsRequest) String() string { return proto.CompactTextString(m) }
func (*ListApplicationsRequest) ProtoMessage()    {}
func (*ListApplicationsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbf1bd0c00876e46, []int{2}
}

func (m *ListApplicationsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListApplicationsRequest.Unmarshal(m, b)
}
func (m *ListApplicationsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListApplicationsRequest.Marshal(b, m, deterministic)
}
func (m *ListApplicationsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListApplicationsRequest.Merge(m, src)
}
func (m *ListApplicationsRequest) XXX_Size() int {
	return xxx_messageInfo_ListApplicationsRequest.Size(m)
}
func (m *ListApplicationsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListApplicationsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListApplicationsRequest proto.InternalMessageInfo

type ListApplicationsResponse struct {
	Applications         []*Application `protobuf:"bytes,1,rep,name=applications,proto3" json:"applications,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ListApplicationsResponse) Reset()         { *m = ListApplicationsResponse{} }
func (m *ListApplicationsResponse) String() string { return proto.CompactTextString(m) }
func (*ListApplicationsResponse) ProtoMessage()    {}
func (*ListApplicationsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cbf1bd0c00876e46, []int{3}
}

func (m *ListApplicationsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListApplicationsResponse.Unmarshal(m, b)
}
func (m *ListApplicationsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListApplicationsResponse.Marshal(b, m, deterministic)
}
func (m *ListApplicationsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListApplicationsResponse.Merge(m, src)
}
func (m *ListApplicationsResponse) XXX_Size() int {
	return xxx_messageInfo_ListApplicationsResponse.Size(m)
}
func (m *ListApplicationsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListApplicationsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListApplicationsResponse proto.InternalMessageInfo

func (m *ListApplicationsResponse) GetApplications() []*Application {
	if m != nil {
		return m.Applications
	}
	return nil
}

func init() {
	proto.RegisterType((*ListApplicationIdentitiesRequest)(nil), "dogma.configkit.v1.ListApplicationIdentitiesRequest")
	proto.RegisterType((*ListApplicationIdentitiesResponse)(nil), "dogma.configkit.v1.ListApplicationIdentitiesResponse")
	proto.RegisterType((*ListApplicationsRequest)(nil), "dogma.configkit.v1.ListApplicationsRequest")
	proto.RegisterType((*ListApplicationsResponse)(nil), "dogma.configkit.v1.ListApplicationsResponse")
}

func init() { proto.RegisterFile("api/internal/pb/service.proto", fileDescriptor_cbf1bd0c00876e46) }

var fileDescriptor_cbf1bd0c00876e46 = []byte{
	// 267 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x31, 0x4f, 0xc3, 0x30,
	0x10, 0x85, 0x55, 0x90, 0x3a, 0x1c, 0x0c, 0xc8, 0x0b, 0x6d, 0x01, 0x51, 0x32, 0x21, 0x01, 0xb6,
	0x28, 0xb0, 0xb1, 0x40, 0x27, 0x24, 0xa6, 0x8c, 0x2c, 0xc8, 0x49, 0x8f, 0x70, 0xa2, 0xb1, 0xdd,
	0xf8, 0x5a, 0xa9, 0x2b, 0x13, 0x3f, 0x1b, 0x29, 0x8d, 0xa2, 0xe0, 0x16, 0x48, 0xd7, 0xcb, 0xfb,
	0xee, 0xcb, 0xb3, 0x0e, 0x4e, 0xb4, 0x23, 0x45, 0x86, 0xb1, 0x30, 0x7a, 0xaa, 0x5c, 0xa2, 0x3c,
	0x16, 0x0b, 0x4a, 0x51, 0xba, 0xc2, 0xb2, 0x15, 0x62, 0x62, 0xb3, 0x5c, 0xcb, 0xd4, 0x9a, 0x37,
	0xca, 0x3e, 0x88, 0xe5, 0xe2, 0x7a, 0x70, 0x14, 0x22, 0xbc, 0x74, 0xe8, 0x57, 0x40, 0x14, 0xc1,
	0xf0, 0x99, 0x3c, 0x3f, 0x38, 0x37, 0xa5, 0x54, 0x33, 0x59, 0xf3, 0x34, 0x41, 0xc3, 0xc4, 0x84,
	0x3e, 0xc6, 0xd9, 0x1c, 0x3d, 0x47, 0x1a, 0xce, 0xfe, 0xc8, 0x78, 0x67, 0x8d, 0x47, 0x71, 0x0f,
	0x40, 0xf5, 0xb4, 0xd7, 0x19, 0xee, 0x9e, 0xef, 0x8d, 0x8e, 0xe5, 0xfa, 0xef, 0xc8, 0x8a, 0x5d,
	0xc6, 0x8d, 0x7c, 0xd4, 0x87, 0xc3, 0x40, 0x51, 0xdb, 0x5f, 0xa1, 0xb7, 0xfe, 0xa9, 0x92, 0x8e,
	0x61, 0x5f, 0x37, 0xe6, 0x95, 0xf6, 0x74, 0x93, 0xb6, 0xc1, 0xc7, 0x3f, 0xa0, 0xd1, 0xe7, 0x0e,
	0x74, 0xc7, 0x65, 0x54, 0x7c, 0x75, 0xa0, 0xff, 0x6b, 0x55, 0x71, 0xbb, 0x69, 0xef, 0x7f, 0xaf,
	0x37, 0xb8, 0xdb, 0x92, 0xaa, 0xaa, 0xe5, 0x70, 0x10, 0xd6, 0x16, 0x17, 0x2d, 0x56, 0xd5, 0xde,
	0xcb, 0x76, 0xe1, 0x95, 0xee, 0x51, 0xbd, 0x5c, 0x65, 0xc4, 0xef, 0xf3, 0x44, 0xa6, 0x36, 0x57,
	0x25, 0xc9, 0x34, 0x53, 0x35, 0xac, 0x82, 0x23, 0x4a, 0xba, 0xe5, 0xfd, 0xdc, 0x7c, 0x07, 0x00,
	0x00, 0xff, 0xff, 0xf4, 0x0d, 0x6d, 0x44, 0x91, 0x02, 0x00, 0x00,
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
}

type configClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigClient(cc grpc.ClientConnInterface) ConfigClient {
	return &configClient{cc}
}

func (c *configClient) ListApplicationIdentities(ctx context.Context, in *ListApplicationIdentitiesRequest, opts ...grpc.CallOption) (*ListApplicationIdentitiesResponse, error) {
	out := new(ListApplicationIdentitiesResponse)
	err := c.cc.Invoke(ctx, "/dogma.configkit.v1.Config/ListApplicationIdentities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configClient) ListApplications(ctx context.Context, in *ListApplicationsRequest, opts ...grpc.CallOption) (*ListApplicationsResponse, error) {
	out := new(ListApplicationsResponse)
	err := c.cc.Invoke(ctx, "/dogma.configkit.v1.Config/ListApplications", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigServer is the server API for Config service.
type ConfigServer interface {
	// ListApplicationIdentities returns the identity of all applications.
	ListApplicationIdentities(context.Context, *ListApplicationIdentitiesRequest) (*ListApplicationIdentitiesResponse, error)
	// ListApplications returns the full configuration of all applications.
	ListApplications(context.Context, *ListApplicationsRequest) (*ListApplicationsResponse, error)
}

// UnimplementedConfigServer can be embedded to have forward compatible implementations.
type UnimplementedConfigServer struct {
}

func (*UnimplementedConfigServer) ListApplicationIdentities(ctx context.Context, req *ListApplicationIdentitiesRequest) (*ListApplicationIdentitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListApplicationIdentities not implemented")
}
func (*UnimplementedConfigServer) ListApplications(ctx context.Context, req *ListApplicationsRequest) (*ListApplicationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListApplications not implemented")
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
		FullMethod: "/dogma.configkit.v1.Config/ListApplicationIdentities",
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
		FullMethod: "/dogma.configkit.v1.Config/ListApplications",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServer).ListApplications(ctx, req.(*ListApplicationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Config_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dogma.configkit.v1.Config",
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
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/internal/pb/service.proto",
}

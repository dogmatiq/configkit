// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/internal/pb/types.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// MessageRole is a protocol buffers representation of the
// configkit/message.Role enumeration.
type MessageRole int32

const (
	MessageRole_UNKNOWN_MESSAGE_ROLE MessageRole = 0
	MessageRole_COMMAND              MessageRole = 1
	MessageRole_EVENT                MessageRole = 2
	MessageRole_TIMEOUT              MessageRole = 3
)

var MessageRole_name = map[int32]string{
	0: "UNKNOWN_MESSAGE_ROLE",
	1: "COMMAND",
	2: "EVENT",
	3: "TIMEOUT",
}

var MessageRole_value = map[string]int32{
	"UNKNOWN_MESSAGE_ROLE": 0,
	"COMMAND":              1,
	"EVENT":                2,
	"TIMEOUT":              3,
}

func (x MessageRole) String() string {
	return proto.EnumName(MessageRole_name, int32(x))
}

func (MessageRole) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_565b712cd72dd24a, []int{0}
}

// HandlerType is a protocol buffers representation of the configkit.HandlerType
// enumeration.
type HandlerType int32

const (
	HandlerType_UNKNOWN_HANDLER_TYPE HandlerType = 0
	HandlerType_AGGREGATE            HandlerType = 1
	HandlerType_PROCESS              HandlerType = 2
	HandlerType_INTEGRATION          HandlerType = 3
	HandlerType_PROJECTION           HandlerType = 4
)

var HandlerType_name = map[int32]string{
	0: "UNKNOWN_HANDLER_TYPE",
	1: "AGGREGATE",
	2: "PROCESS",
	3: "INTEGRATION",
	4: "PROJECTION",
}

var HandlerType_value = map[string]int32{
	"UNKNOWN_HANDLER_TYPE": 0,
	"AGGREGATE":            1,
	"PROCESS":              2,
	"INTEGRATION":          3,
	"PROJECTION":           4,
}

func (x HandlerType) String() string {
	return proto.EnumName(HandlerType_name, int32(x))
}

func (HandlerType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_565b712cd72dd24a, []int{1}
}

// Identity is a protocol buffers representation of the configkit.Identity type.
type Identity struct {
	// Name is the entity's unique name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Key is the entity's immutable, unique key.
	Key                  string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Identity) Reset()         { *m = Identity{} }
func (m *Identity) String() string { return proto.CompactTextString(m) }
func (*Identity) ProtoMessage()    {}
func (*Identity) Descriptor() ([]byte, []int) {
	return fileDescriptor_565b712cd72dd24a, []int{0}
}

func (m *Identity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Identity.Unmarshal(m, b)
}
func (m *Identity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Identity.Marshal(b, m, deterministic)
}
func (m *Identity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Identity.Merge(m, src)
}
func (m *Identity) XXX_Size() int {
	return xxx_messageInfo_Identity.Size(m)
}
func (m *Identity) XXX_DiscardUnknown() {
	xxx_messageInfo_Identity.DiscardUnknown(m)
}

var xxx_messageInfo_Identity proto.InternalMessageInfo

func (m *Identity) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Identity) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

// Application is a protocol buffers representation of the configkit.Application
// interface.
type Application struct {
	// Identity is the application's identity.
	Identity *Identity `protobuf:"bytes,1,opt,name=identity,proto3" json:"identity,omitempty"`
	// TypeName is the fully-qualified name of the Go type used to implement the
	// application.
	TypeName string `protobuf:"bytes,2,opt,name=type_name,json=typeName,proto3" json:"type_name,omitempty"`
	// Messages is an ordered-sequence of message name / role pairs.
	//
	// This directly correlates to the configkit.Application.MessageNames().Roles
	// value. The produced/consumed message names are not encoded directly in the
	// application, but rather rebuilt from the handlers when the application is
	// unmarshaled.
	Messages []*NameRole `protobuf:"bytes,3,rep,name=messages,proto3" json:"messages,omitempty"`
	// Handlers is the set of handlers within the application.
	Handlers             []*Handler `protobuf:"bytes,4,rep,name=handlers,proto3" json:"handlers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Application) Reset()         { *m = Application{} }
func (m *Application) String() string { return proto.CompactTextString(m) }
func (*Application) ProtoMessage()    {}
func (*Application) Descriptor() ([]byte, []int) {
	return fileDescriptor_565b712cd72dd24a, []int{1}
}

func (m *Application) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Application.Unmarshal(m, b)
}
func (m *Application) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Application.Marshal(b, m, deterministic)
}
func (m *Application) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Application.Merge(m, src)
}
func (m *Application) XXX_Size() int {
	return xxx_messageInfo_Application.Size(m)
}
func (m *Application) XXX_DiscardUnknown() {
	xxx_messageInfo_Application.DiscardUnknown(m)
}

var xxx_messageInfo_Application proto.InternalMessageInfo

func (m *Application) GetIdentity() *Identity {
	if m != nil {
		return m.Identity
	}
	return nil
}

func (m *Application) GetTypeName() string {
	if m != nil {
		return m.TypeName
	}
	return ""
}

func (m *Application) GetMessages() []*NameRole {
	if m != nil {
		return m.Messages
	}
	return nil
}

func (m *Application) GetHandlers() []*Handler {
	if m != nil {
		return m.Handlers
	}
	return nil
}

// NameRole is a 2-tuple containing a message name and its role.
type NameRole struct {
	// Name is the fully-qualified message name.
	Name []byte `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Role is the role this message plays within the application.
	Role                 MessageRole `protobuf:"varint,2,opt,name=role,proto3,enum=dogma.configkit.v1.MessageRole" json:"role,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *NameRole) Reset()         { *m = NameRole{} }
func (m *NameRole) String() string { return proto.CompactTextString(m) }
func (*NameRole) ProtoMessage()    {}
func (*NameRole) Descriptor() ([]byte, []int) {
	return fileDescriptor_565b712cd72dd24a, []int{2}
}

func (m *NameRole) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NameRole.Unmarshal(m, b)
}
func (m *NameRole) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NameRole.Marshal(b, m, deterministic)
}
func (m *NameRole) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NameRole.Merge(m, src)
}
func (m *NameRole) XXX_Size() int {
	return xxx_messageInfo_NameRole.Size(m)
}
func (m *NameRole) XXX_DiscardUnknown() {
	xxx_messageInfo_NameRole.DiscardUnknown(m)
}

var xxx_messageInfo_NameRole proto.InternalMessageInfo

func (m *NameRole) GetName() []byte {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *NameRole) GetRole() MessageRole {
	if m != nil {
		return m.Role
	}
	return MessageRole_UNKNOWN_MESSAGE_ROLE
}

// Handler is a protocol buffers representation of the configkit.Handler
// interface.
type Handler struct {
	// Identity is the handler's identity.
	Identity *Identity `protobuf:"bytes,1,opt,name=identity,proto3" json:"identity,omitempty"`
	// TypeName is the fully-qualified name of the Go type used to implement the
	// handler.
	TypeName string `protobuf:"bytes,2,opt,name=type_name,json=typeName,proto3" json:"type_name,omitempty"`
	// Type is the handler's type.
	Type HandlerType `protobuf:"varint,3,opt,name=type,proto3,enum=dogma.configkit.v1.HandlerType" json:"type,omitempty"`
	// Produced is a list of the messages produced by this handler.
	// Each value is the index of a MessageRolePair within the application.
	Produced []uint32 `protobuf:"varint,4,rep,packed,name=produced,proto3" json:"produced,omitempty"`
	// Consumed is a list of the messages consumed by this handler.
	// Each value is the index of a MessageRolePair within the application.
	Consumed             []uint32 `protobuf:"varint,5,rep,packed,name=consumed,proto3" json:"consumed,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Handler) Reset()         { *m = Handler{} }
func (m *Handler) String() string { return proto.CompactTextString(m) }
func (*Handler) ProtoMessage()    {}
func (*Handler) Descriptor() ([]byte, []int) {
	return fileDescriptor_565b712cd72dd24a, []int{3}
}

func (m *Handler) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Handler.Unmarshal(m, b)
}
func (m *Handler) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Handler.Marshal(b, m, deterministic)
}
func (m *Handler) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Handler.Merge(m, src)
}
func (m *Handler) XXX_Size() int {
	return xxx_messageInfo_Handler.Size(m)
}
func (m *Handler) XXX_DiscardUnknown() {
	xxx_messageInfo_Handler.DiscardUnknown(m)
}

var xxx_messageInfo_Handler proto.InternalMessageInfo

func (m *Handler) GetIdentity() *Identity {
	if m != nil {
		return m.Identity
	}
	return nil
}

func (m *Handler) GetTypeName() string {
	if m != nil {
		return m.TypeName
	}
	return ""
}

func (m *Handler) GetType() HandlerType {
	if m != nil {
		return m.Type
	}
	return HandlerType_UNKNOWN_HANDLER_TYPE
}

func (m *Handler) GetProduced() []uint32 {
	if m != nil {
		return m.Produced
	}
	return nil
}

func (m *Handler) GetConsumed() []uint32 {
	if m != nil {
		return m.Consumed
	}
	return nil
}

func init() {
	proto.RegisterEnum("dogma.configkit.v1.MessageRole", MessageRole_name, MessageRole_value)
	proto.RegisterEnum("dogma.configkit.v1.HandlerType", HandlerType_name, HandlerType_value)
	proto.RegisterType((*Identity)(nil), "dogma.configkit.v1.Identity")
	proto.RegisterType((*Application)(nil), "dogma.configkit.v1.Application")
	proto.RegisterType((*NameRole)(nil), "dogma.configkit.v1.NameRole")
	proto.RegisterType((*Handler)(nil), "dogma.configkit.v1.Handler")
}

func init() { proto.RegisterFile("api/internal/pb/types.proto", fileDescriptor_565b712cd72dd24a) }

var fileDescriptor_565b712cd72dd24a = []byte{
	// 460 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x53, 0x41, 0x8f, 0x93, 0x40,
	0x14, 0x96, 0xc2, 0xba, 0xf4, 0xe1, 0xae, 0x64, 0xe2, 0x81, 0x58, 0x13, 0x37, 0x3d, 0x99, 0x4d,
	0x04, 0xdd, 0x1e, 0xdc, 0x2b, 0x76, 0x27, 0x6c, 0xb5, 0x40, 0x33, 0xb0, 0x1a, 0xbd, 0x34, 0x14,
	0xc6, 0xee, 0x64, 0x81, 0x41, 0xa0, 0x26, 0xfd, 0x91, 0x9e, 0xfd, 0x3b, 0x66, 0x86, 0x96, 0x54,
	0xdd, 0x3d, 0x7a, 0x7b, 0xef, 0xcd, 0xf7, 0xbe, 0xef, 0x9b, 0x2f, 0x79, 0x30, 0x4a, 0x2a, 0xe6,
	0xb0, 0xb2, 0xa5, 0x75, 0x99, 0xe4, 0x4e, 0xb5, 0x72, 0xda, 0x6d, 0x45, 0x1b, 0xbb, 0xaa, 0x79,
	0xcb, 0x11, 0xca, 0xf8, 0xba, 0x48, 0xec, 0x94, 0x97, 0xdf, 0xd8, 0xfa, 0x8e, 0xb5, 0xf6, 0x8f,
	0xb7, 0xe3, 0x37, 0xa0, 0xcf, 0x32, 0x5a, 0xb6, 0xac, 0xdd, 0x22, 0x04, 0x5a, 0x99, 0x14, 0xd4,
	0x52, 0xce, 0x94, 0x57, 0x43, 0x22, 0x6b, 0x64, 0x82, 0x7a, 0x47, 0xb7, 0xd6, 0x40, 0x8e, 0x44,
	0x39, 0xfe, 0xa5, 0x80, 0xe1, 0x56, 0x55, 0xce, 0xd2, 0xa4, 0x65, 0xbc, 0x44, 0x97, 0xa0, 0xb3,
	0x1d, 0x83, 0xdc, 0x34, 0x2e, 0x5e, 0xd8, 0xff, 0x0a, 0xd9, 0x7b, 0x15, 0xd2, 0xa3, 0xd1, 0x08,
	0x86, 0xc2, 0xde, 0x52, 0x8a, 0x76, 0x0a, 0xba, 0x18, 0x04, 0x42, 0xf8, 0x12, 0xf4, 0x82, 0x36,
	0x4d, 0xb2, 0xa6, 0x8d, 0xa5, 0x9e, 0xa9, 0x0f, 0xd1, 0x0a, 0x2c, 0xe1, 0x39, 0x25, 0x3d, 0x1a,
	0xbd, 0x03, 0xfd, 0x36, 0x29, 0xb3, 0x9c, 0xd6, 0x8d, 0xa5, 0xc9, 0xcd, 0xd1, 0x7d, 0x9b, 0xd7,
	0x1d, 0x86, 0xf4, 0xe0, 0x71, 0x04, 0xfa, 0x9e, 0xee, 0x8f, 0x2c, 0x9e, 0xec, 0xb2, 0x98, 0x80,
	0x56, 0xf3, 0xbc, 0xb3, 0x7a, 0x7a, 0xf1, 0xf2, 0x3e, 0x52, 0xbf, 0x33, 0x21, 0x1d, 0x49, 0xf0,
	0xf8, 0xa7, 0x02, 0xc7, 0x3b, 0xa9, 0xff, 0x15, 0xd5, 0x04, 0x34, 0x51, 0x5b, 0xea, 0xc3, 0xbe,
	0x76, 0x0e, 0xe2, 0x6d, 0x45, 0x89, 0x04, 0xa3, 0xe7, 0xa0, 0x57, 0x35, 0xcf, 0x36, 0x29, 0xcd,
	0x64, 0x4a, 0x27, 0xa4, 0xef, 0xc5, 0x5b, 0xca, 0xcb, 0x66, 0x53, 0xd0, 0xcc, 0x3a, 0xea, 0xde,
	0xf6, 0xfd, 0xf9, 0x1c, 0x8c, 0x83, 0x4f, 0x22, 0x0b, 0x9e, 0xdd, 0x04, 0x1f, 0x83, 0xf0, 0x73,
	0xb0, 0xf4, 0x71, 0x14, 0xb9, 0x1e, 0x5e, 0x92, 0x70, 0x8e, 0xcd, 0x47, 0xc8, 0x80, 0xe3, 0x69,
	0xe8, 0xfb, 0x6e, 0x70, 0x65, 0x2a, 0x68, 0x08, 0x47, 0xf8, 0x13, 0x0e, 0x62, 0x73, 0x20, 0xe6,
	0xf1, 0xcc, 0xc7, 0xe1, 0x4d, 0x6c, 0xaa, 0xe7, 0x19, 0x18, 0x07, 0xd6, 0x0e, 0xd9, 0xae, 0xdd,
	0xe0, 0x6a, 0x8e, 0xc9, 0x32, 0xfe, 0xb2, 0x10, 0x6c, 0x27, 0x30, 0x74, 0x3d, 0x8f, 0x60, 0xcf,
	0x8d, 0xb1, 0xa9, 0x08, 0x92, 0x05, 0x09, 0xa7, 0x38, 0x8a, 0xcc, 0x01, 0x7a, 0x0a, 0xc6, 0x2c,
	0x88, 0xb1, 0x47, 0xdc, 0x78, 0x16, 0x06, 0xa6, 0x8a, 0x4e, 0x01, 0x16, 0x24, 0xfc, 0x80, 0xa7,
	0xb2, 0xd7, 0xde, 0x3b, 0x5f, 0x5f, 0xaf, 0x59, 0x7b, 0xbb, 0x59, 0xd9, 0x29, 0x2f, 0x1c, 0x19,
	0x4f, 0xcb, 0xbe, 0x3b, 0x7d, 0x42, 0xce, 0x5f, 0x57, 0xb3, 0x7a, 0x2c, 0x0f, 0x66, 0xf2, 0x3b,
	0x00, 0x00, 0xff, 0xff, 0x9c, 0x60, 0xa5, 0x67, 0x4f, 0x03, 0x00, 0x00,
}

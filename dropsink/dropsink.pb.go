// Code generated by protoc-gen-go. DO NOT EDIT.
// source: dropsink.proto

package dropsink

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

type DropEntry struct {
	// TODO native Timestamp is available
	Timestamp            int64             `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Fields               map[string]string `protobuf:"bytes,3,rep,name=fields,proto3" json:"fields,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *DropEntry) Reset()         { *m = DropEntry{} }
func (m *DropEntry) String() string { return proto.CompactTextString(m) }
func (*DropEntry) ProtoMessage()    {}
func (*DropEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_1221e8e7975cbdf6, []int{0}
}

func (m *DropEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DropEntry.Unmarshal(m, b)
}
func (m *DropEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DropEntry.Marshal(b, m, deterministic)
}
func (m *DropEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DropEntry.Merge(m, src)
}
func (m *DropEntry) XXX_Size() int {
	return xxx_messageInfo_DropEntry.Size(m)
}
func (m *DropEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_DropEntry.DiscardUnknown(m)
}

var xxx_messageInfo_DropEntry proto.InternalMessageInfo

func (m *DropEntry) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *DropEntry) GetFields() map[string]string {
	if m != nil {
		return m.Fields
	}
	return nil
}

type Void struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Void) Reset()         { *m = Void{} }
func (m *Void) String() string { return proto.CompactTextString(m) }
func (*Void) ProtoMessage()    {}
func (*Void) Descriptor() ([]byte, []int) {
	return fileDescriptor_1221e8e7975cbdf6, []int{1}
}

func (m *Void) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Void.Unmarshal(m, b)
}
func (m *Void) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Void.Marshal(b, m, deterministic)
}
func (m *Void) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Void.Merge(m, src)
}
func (m *Void) XXX_Size() int {
	return xxx_messageInfo_Void.Size(m)
}
func (m *Void) XXX_DiscardUnknown() {
	xxx_messageInfo_Void.DiscardUnknown(m)
}

var xxx_messageInfo_Void proto.InternalMessageInfo

func init() {
	proto.RegisterType((*DropEntry)(nil), "dropsink.DropEntry")
	proto.RegisterMapType((map[string]string)(nil), "dropsink.DropEntry.FieldsEntry")
	proto.RegisterType((*Void)(nil), "dropsink.Void")
}

func init() { proto.RegisterFile("dropsink.proto", fileDescriptor_1221e8e7975cbdf6) }

var fileDescriptor_1221e8e7975cbdf6 = []byte{
	// 191 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0x29, 0xca, 0x2f,
	0x28, 0xce, 0xcc, 0xcb, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x80, 0xf1, 0x95, 0xe6,
	0x32, 0x72, 0x71, 0xba, 0x00, 0x39, 0xae, 0x79, 0x25, 0x45, 0x95, 0x42, 0x32, 0x5c, 0x9c, 0x25,
	0x99, 0xb9, 0xa9, 0xc5, 0x25, 0x89, 0xb9, 0x05, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xcc, 0x41, 0x08,
	0x01, 0x21, 0x73, 0x2e, 0xb6, 0xb4, 0xcc, 0xd4, 0x9c, 0x94, 0x62, 0x09, 0x66, 0x05, 0x66, 0x0d,
	0x6e, 0x23, 0x79, 0x3d, 0xb8, 0xb1, 0x70, 0x23, 0xf4, 0xdc, 0xc0, 0x2a, 0xc0, 0xec, 0x20, 0xa8,
	0x72, 0x29, 0x4b, 0x2e, 0x6e, 0x24, 0x61, 0x21, 0x01, 0x2e, 0xe6, 0xec, 0xd4, 0x4a, 0xb0, 0xf9,
	0x9c, 0x41, 0x20, 0xa6, 0x90, 0x08, 0x17, 0x6b, 0x59, 0x62, 0x4e, 0x69, 0xaa, 0x04, 0x13, 0x58,
	0x0c, 0xc2, 0xb1, 0x62, 0xb2, 0x60, 0x54, 0x62, 0xe3, 0x62, 0x09, 0xcb, 0xcf, 0x4c, 0x31, 0xb2,
	0xe4, 0xe2, 0x00, 0xd9, 0x11, 0x0c, 0xb4, 0x4c, 0x48, 0x97, 0x8b, 0x25, 0xa0, 0xb4, 0x38, 0x43,
	0x48, 0x18, 0x8b, 0xfd, 0x52, 0x7c, 0x08, 0x41, 0x90, 0x46, 0x25, 0x86, 0x24, 0x36, 0xb0, 0x9f,
	0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x4f, 0xba, 0x74, 0x3b, 0x05, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DropSinkClient is the client API for DropSink service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DropSinkClient interface {
	Push(ctx context.Context, in *DropEntry, opts ...grpc.CallOption) (*Void, error)
}

type dropSinkClient struct {
	cc *grpc.ClientConn
}

func NewDropSinkClient(cc *grpc.ClientConn) DropSinkClient {
	return &dropSinkClient{cc}
}

func (c *dropSinkClient) Push(ctx context.Context, in *DropEntry, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/dropsink.DropSink/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DropSinkServer is the server API for DropSink service.
type DropSinkServer interface {
	Push(context.Context, *DropEntry) (*Void, error)
}

// UnimplementedDropSinkServer can be embedded to have forward compatible implementations.
type UnimplementedDropSinkServer struct {
}

func (*UnimplementedDropSinkServer) Push(ctx context.Context, req *DropEntry) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}

func RegisterDropSinkServer(s *grpc.Server, srv DropSinkServer) {
	s.RegisterService(&_DropSink_serviceDesc, srv)
}

func _DropSink_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DropEntry)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DropSinkServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dropsink.DropSink/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DropSinkServer).Push(ctx, req.(*DropEntry))
	}
	return interceptor(ctx, in, info, handler)
}

var _DropSink_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dropsink.DropSink",
	HandlerType: (*DropSinkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _DropSink_Push_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dropsink.proto",
}
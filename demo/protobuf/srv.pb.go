// Code generated by protoc-gen-go. DO NOT EDIT.
// source: srv.proto

package example

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

type StringMessage struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StringMessage) Reset()         { *m = StringMessage{} }
func (m *StringMessage) String() string { return proto.CompactTextString(m) }
func (*StringMessage) ProtoMessage()    {}
func (*StringMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_2bbe8325d22c1a26, []int{0}
}

func (m *StringMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StringMessage.Unmarshal(m, b)
}
func (m *StringMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StringMessage.Marshal(b, m, deterministic)
}
func (m *StringMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StringMessage.Merge(m, src)
}
func (m *StringMessage) XXX_Size() int {
	return xxx_messageInfo_StringMessage.Size(m)
}
func (m *StringMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_StringMessage.DiscardUnknown(m)
}

var xxx_messageInfo_StringMessage proto.InternalMessageInfo

func (m *StringMessage) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*StringMessage)(nil), "example.StringMessage")
}

func init() {
	proto.RegisterFile("srv.proto", fileDescriptor_2bbe8325d22c1a26)
}

var fileDescriptor_2bbe8325d22c1a26 = []byte{
	// 120 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2c, 0x2e, 0x2a, 0xd3,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4f, 0xad, 0x48, 0xcc, 0x2d, 0xc8, 0x49, 0x55, 0x52,
	0xe5, 0xe2, 0x0d, 0x2e, 0x29, 0xca, 0xcc, 0x4b, 0xf7, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x15,
	0x12, 0xe1, 0x62, 0x2d, 0x4b, 0xcc, 0x29, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82,
	0x70, 0x8c, 0xdc, 0xb9, 0xb8, 0x23, 0xf3, 0x4b, 0x8b, 0x82, 0x53, 0x8b, 0xca, 0x32, 0x93, 0x53,
	0x85, 0x2c, 0xb8, 0x58, 0x5c, 0x93, 0x33, 0xf2, 0x85, 0xc4, 0xf4, 0xa0, 0xe6, 0xe8, 0xa1, 0x18,
	0x22, 0x85, 0x43, 0x5c, 0x89, 0x21, 0x89, 0x0d, 0x6c, 0xbf, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff,
	0x99, 0xa8, 0x18, 0xb2, 0x8c, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// YourServiceClient is the client API for YourService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type YourServiceClient interface {
	Echo(ctx context.Context, in *StringMessage, opts ...grpc.CallOption) (*StringMessage, error)
}

type yourServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewYourServiceClient(cc grpc.ClientConnInterface) YourServiceClient {
	return &yourServiceClient{cc}
}

func (c *yourServiceClient) Echo(ctx context.Context, in *StringMessage, opts ...grpc.CallOption) (*StringMessage, error) {
	out := new(StringMessage)
	err := c.cc.Invoke(ctx, "/example.YourService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// YourServiceServer is the server API for YourService service.
type YourServiceServer interface {
	Echo(context.Context, *StringMessage) (*StringMessage, error)
}

// UnimplementedYourServiceServer can be embedded to have forward compatible implementations.
type UnimplementedYourServiceServer struct {
}

func (*UnimplementedYourServiceServer) Echo(ctx context.Context, req *StringMessage) (*StringMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}

func RegisterYourServiceServer(s *grpc.Server, srv YourServiceServer) {
	s.RegisterService(&_YourService_serviceDesc, srv)
}

func _YourService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StringMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(YourServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/example.YourService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(YourServiceServer).Echo(ctx, req.(*StringMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _YourService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "example.YourService",
	HandlerType: (*YourServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _YourService_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "srv.proto",
}

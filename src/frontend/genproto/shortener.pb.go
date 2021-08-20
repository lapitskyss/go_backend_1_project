// Code generated by protoc-gen-go. DO NOT EDIT.
// source: shortener.proto

package shortener

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

type GetLinkRequest struct {
	Hash                 string   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	UserAgent            string   `protobuf:"bytes,2,opt,name=user_agent,json=userAgent,proto3" json:"user_agent,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetLinkRequest) Reset()         { *m = GetLinkRequest{} }
func (m *GetLinkRequest) String() string { return proto.CompactTextString(m) }
func (*GetLinkRequest) ProtoMessage()    {}
func (*GetLinkRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{0}
}

func (m *GetLinkRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetLinkRequest.Unmarshal(m, b)
}
func (m *GetLinkRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetLinkRequest.Marshal(b, m, deterministic)
}
func (m *GetLinkRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetLinkRequest.Merge(m, src)
}
func (m *GetLinkRequest) XXX_Size() int {
	return xxx_messageInfo_GetLinkRequest.Size(m)
}
func (m *GetLinkRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetLinkRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetLinkRequest proto.InternalMessageInfo

func (m *GetLinkRequest) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *GetLinkRequest) GetUserAgent() string {
	if m != nil {
		return m.UserAgent
	}
	return ""
}

type Link struct {
	Url                  string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Hash                 string   `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Link) Reset()         { *m = Link{} }
func (m *Link) String() string { return proto.CompactTextString(m) }
func (*Link) ProtoMessage()    {}
func (*Link) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a64040fb43d257f, []int{1}
}

func (m *Link) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Link.Unmarshal(m, b)
}
func (m *Link) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Link.Marshal(b, m, deterministic)
}
func (m *Link) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Link.Merge(m, src)
}
func (m *Link) XXX_Size() int {
	return xxx_messageInfo_Link.Size(m)
}
func (m *Link) XXX_DiscardUnknown() {
	xxx_messageInfo_Link.DiscardUnknown(m)
}

var xxx_messageInfo_Link proto.InternalMessageInfo

func (m *Link) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Link) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func init() {
	proto.RegisterType((*GetLinkRequest)(nil), "shortener.GetLinkRequest")
	proto.RegisterType((*Link)(nil), "shortener.Link")
}

func init() { proto.RegisterFile("shortener.proto", fileDescriptor_6a64040fb43d257f) }

var fileDescriptor_6a64040fb43d257f = []byte{
	// 165 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0xce, 0xc8, 0x2f,
	0x2a, 0x49, 0xcd, 0x4b, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x84, 0x0b, 0x28,
	0x39, 0x73, 0xf1, 0xb9, 0xa7, 0x96, 0xf8, 0x64, 0xe6, 0x65, 0x07, 0xa5, 0x16, 0x96, 0xa6, 0x16,
	0x97, 0x08, 0x09, 0x71, 0xb1, 0x64, 0x24, 0x16, 0x67, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06,
	0x81, 0xd9, 0x42, 0xb2, 0x5c, 0x5c, 0xa5, 0xc5, 0xa9, 0x45, 0xf1, 0x89, 0xe9, 0xa9, 0x79, 0x25,
	0x12, 0x4c, 0x60, 0x19, 0x4e, 0x90, 0x88, 0x23, 0x48, 0x40, 0x49, 0x87, 0x8b, 0x05, 0x64, 0x82,
	0x90, 0x00, 0x17, 0x73, 0x69, 0x51, 0x0e, 0x54, 0x27, 0x88, 0x09, 0x37, 0x8c, 0x09, 0x61, 0x98,
	0x91, 0x1b, 0x17, 0x37, 0x48, 0x75, 0x70, 0x6a, 0x51, 0x59, 0x66, 0x72, 0xaa, 0x90, 0x39, 0x17,
	0x3b, 0xd4, 0x05, 0x42, 0x92, 0x7a, 0x08, 0x97, 0xa2, 0xba, 0x4a, 0x8a, 0x1f, 0x49, 0x0a, 0x24,
	0xae, 0xc4, 0x90, 0xc4, 0x06, 0xf6, 0x8c, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xac, 0x63, 0x60,
	0x91, 0xdf, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LinkServiceClient is the client API for LinkService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LinkServiceClient interface {
	GetLink(ctx context.Context, in *GetLinkRequest, opts ...grpc.CallOption) (*Link, error)
}

type linkServiceClient struct {
	cc *grpc.ClientConn
}

func NewLinkServiceClient(cc *grpc.ClientConn) LinkServiceClient {
	return &linkServiceClient{cc}
}

func (c *linkServiceClient) GetLink(ctx context.Context, in *GetLinkRequest, opts ...grpc.CallOption) (*Link, error) {
	out := new(Link)
	err := c.cc.Invoke(ctx, "/shortener.LinkService/GetLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LinkServiceServer is the server API for LinkService service.
type LinkServiceServer interface {
	GetLink(context.Context, *GetLinkRequest) (*Link, error)
}

// UnimplementedLinkServiceServer can be embedded to have forward compatible implementations.
type UnimplementedLinkServiceServer struct {
}

func (*UnimplementedLinkServiceServer) GetLink(ctx context.Context, req *GetLinkRequest) (*Link, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLink not implemented")
}

func RegisterLinkServiceServer(s *grpc.Server, srv LinkServiceServer) {
	s.RegisterService(&_LinkService_serviceDesc, srv)
}

func _LinkService_GetLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkServiceServer).GetLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.LinkService/GetLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkServiceServer).GetLink(ctx, req.(*GetLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LinkService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.LinkService",
	HandlerType: (*LinkServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLink",
			Handler:    _LinkService_GetLink_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shortener.proto",
}
// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: say_hello.proto

package hello

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

// HelloClient is the client API for Hello service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HelloClient interface {
	SayHello(ctx context.Context, in *HelloMsg, opts ...grpc.CallOption) (*HelloRes, error)
	AddNode(ctx context.Context, in *AddNodeRequest, opts ...grpc.CallOption) (*AddNodeResponse, error)
}

type helloClient struct {
	cc grpc.ClientConnInterface
}

func NewHelloClient(cc grpc.ClientConnInterface) HelloClient {
	return &helloClient{cc}
}

func (c *helloClient) SayHello(ctx context.Context, in *HelloMsg, opts ...grpc.CallOption) (*HelloRes, error) {
	out := new(HelloRes)
	err := c.cc.Invoke(ctx, "/hello.Hello/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *helloClient) AddNode(ctx context.Context, in *AddNodeRequest, opts ...grpc.CallOption) (*AddNodeResponse, error) {
	out := new(AddNodeResponse)
	err := c.cc.Invoke(ctx, "/hello.Hello/AddNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HelloServer is the server API for Hello service.
// All implementations must embed UnimplementedHelloServer
// for forward compatibility
type HelloServer interface {
	SayHello(context.Context, *HelloMsg) (*HelloRes, error)
	AddNode(context.Context, *AddNodeRequest) (*AddNodeResponse, error)
	mustEmbedUnimplementedHelloServer()
}

// UnimplementedHelloServer must be embedded to have forward compatible implementations.
type UnimplementedHelloServer struct {
}

func (UnimplementedHelloServer) SayHello(context.Context, *HelloMsg) (*HelloRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedHelloServer) AddNode(context.Context, *AddNodeRequest) (*AddNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNode not implemented")
}
func (UnimplementedHelloServer) mustEmbedUnimplementedHelloServer() {}

// UnsafeHelloServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HelloServer will
// result in compilation errors.
type UnsafeHelloServer interface {
	mustEmbedUnimplementedHelloServer()
}

func RegisterHelloServer(s grpc.ServiceRegistrar, srv HelloServer) {
	s.RegisterService(&Hello_ServiceDesc, srv)
}

func _Hello_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hello.Hello/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloServer).SayHello(ctx, req.(*HelloMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hello_AddNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloServer).AddNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hello.Hello/AddNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloServer).AddNode(ctx, req.(*AddNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Hello_ServiceDesc is the grpc.ServiceDesc for Hello service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Hello_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hello.Hello",
	HandlerType: (*HelloServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Hello_SayHello_Handler,
		},
		{
			MethodName: "AddNode",
			Handler:    _Hello_AddNode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "say_hello.proto",
}

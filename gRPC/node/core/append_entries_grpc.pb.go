// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: node/core/append_entries.proto

package gRPC

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

// AppendEntriesServiceClient is the client API for AppendEntriesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AppendEntriesServiceClient interface {
	AppendEntries(ctx context.Context, in *AppendEntriesArgs, opts ...grpc.CallOption) (*AppendEntriesReply, error)
}

type appendEntriesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAppendEntriesServiceClient(cc grpc.ClientConnInterface) AppendEntriesServiceClient {
	return &appendEntriesServiceClient{cc}
}

func (c *appendEntriesServiceClient) AppendEntries(ctx context.Context, in *AppendEntriesArgs, opts ...grpc.CallOption) (*AppendEntriesReply, error) {
	out := new(AppendEntriesReply)
	err := c.cc.Invoke(ctx, "/core.AppendEntriesService/AppendEntries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AppendEntriesServiceServer is the server API for AppendEntriesService service.
// All implementations must embed UnimplementedAppendEntriesServiceServer
// for forward compatibility
type AppendEntriesServiceServer interface {
	AppendEntries(context.Context, *AppendEntriesArgs) (*AppendEntriesReply, error)
	mustEmbedUnimplementedAppendEntriesServiceServer()
}

// UnimplementedAppendEntriesServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAppendEntriesServiceServer struct {
}

func (UnimplementedAppendEntriesServiceServer) AppendEntries(context.Context, *AppendEntriesArgs) (*AppendEntriesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (UnimplementedAppendEntriesServiceServer) mustEmbedUnimplementedAppendEntriesServiceServer() {}

// UnsafeAppendEntriesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppendEntriesServiceServer will
// result in compilation errors.
type UnsafeAppendEntriesServiceServer interface {
	mustEmbedUnimplementedAppendEntriesServiceServer()
}

func RegisterAppendEntriesServiceServer(s grpc.ServiceRegistrar, srv AppendEntriesServiceServer) {
	s.RegisterService(&AppendEntriesService_ServiceDesc, srv)
}

func _AppendEntriesService_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppendEntriesServiceServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.AppendEntriesService/AppendEntries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppendEntriesServiceServer).AppendEntries(ctx, req.(*AppendEntriesArgs))
	}
	return interceptor(ctx, in, info, handler)
}

// AppendEntriesService_ServiceDesc is the grpc.ServiceDesc for AppendEntriesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AppendEntriesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "core.AppendEntriesService",
	HandlerType: (*AppendEntriesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppendEntries",
			Handler:    _AppendEntriesService_AppendEntries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node/core/append_entries.proto",
}

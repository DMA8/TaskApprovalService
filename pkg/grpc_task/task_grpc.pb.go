// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.13.0
// source: proto/task.proto

package grpc_task

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

// GrpcTaskClient is the client API for GrpcTask service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GrpcTaskClient interface {
	PushTask(ctx context.Context, in *TaskMessage, opts ...grpc.CallOption) (*TaskResponse, error)
	PushMail(ctx context.Context, in *Mail, opts ...grpc.CallOption) (*TaskResponse, error)
}

type grpcTaskClient struct {
	cc grpc.ClientConnInterface
}

func NewGrpcTaskClient(cc grpc.ClientConnInterface) GrpcTaskClient {
	return &grpcTaskClient{cc}
}

func (c *grpcTaskClient) PushTask(ctx context.Context, in *TaskMessage, opts ...grpc.CallOption) (*TaskResponse, error) {
	out := new(TaskResponse)
	err := c.cc.Invoke(ctx, "/grpcTask.GrpcTask/PushTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcTaskClient) PushMail(ctx context.Context, in *Mail, opts ...grpc.CallOption) (*TaskResponse, error) {
	out := new(TaskResponse)
	err := c.cc.Invoke(ctx, "/grpcTask.GrpcTask/PushMail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GrpcTaskServer is the server API for GrpcTask service.
// All implementations must embed UnimplementedGrpcTaskServer
// for forward compatibility
type GrpcTaskServer interface {
	PushTask(context.Context, *TaskMessage) (*TaskResponse, error)
	PushMail(context.Context, *Mail) (*TaskResponse, error)
	mustEmbedUnimplementedGrpcTaskServer()
}

// UnimplementedGrpcTaskServer must be embedded to have forward compatible implementations.
type UnimplementedGrpcTaskServer struct {
}

func (UnimplementedGrpcTaskServer) PushTask(context.Context, *TaskMessage) (*TaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushTask not implemented")
}
func (UnimplementedGrpcTaskServer) PushMail(context.Context, *Mail) (*TaskResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushMail not implemented")
}
func (UnimplementedGrpcTaskServer) mustEmbedUnimplementedGrpcTaskServer() {}

// UnsafeGrpcTaskServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GrpcTaskServer will
// result in compilation errors.
type UnsafeGrpcTaskServer interface {
	mustEmbedUnimplementedGrpcTaskServer()
}

func RegisterGrpcTaskServer(s grpc.ServiceRegistrar, srv GrpcTaskServer) {
	s.RegisterService(&GrpcTask_ServiceDesc, srv)
}

func _GrpcTask_PushTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcTaskServer).PushTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpcTask.GrpcTask/PushTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcTaskServer).PushTask(ctx, req.(*TaskMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcTask_PushMail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Mail)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcTaskServer).PushMail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpcTask.GrpcTask/PushMail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcTaskServer).PushMail(ctx, req.(*Mail))
	}
	return interceptor(ctx, in, info, handler)
}

// GrpcTask_ServiceDesc is the grpc.ServiceDesc for GrpcTask service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GrpcTask_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpcTask.GrpcTask",
	HandlerType: (*GrpcTaskServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushTask",
			Handler:    _GrpcTask_PushTask_Handler,
		},
		{
			MethodName: "PushMail",
			Handler:    _GrpcTask_PushMail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/task.proto",
}

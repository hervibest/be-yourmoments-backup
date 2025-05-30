// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: pb/ai/ai.proto

package aipb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AiService_ProcessPhoto_FullMethodName     = "/ai.AiService/ProcessPhoto"
	AiService_ProcessFacecam_FullMethodName   = "/ai.AiService/ProcessFacecam"
	AiService_ProcessBulkPhoto_FullMethodName = "/ai.AiService/ProcessBulkPhoto"
)

// AiServiceClient is the client API for AiService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AiServiceClient interface {
	ProcessPhoto(ctx context.Context, in *ProcessPhotoRequest, opts ...grpc.CallOption) (*ProcessPhotoResponse, error)
	ProcessFacecam(ctx context.Context, in *ProcessFacecamRequest, opts ...grpc.CallOption) (*ProcessFacecamResponse, error)
	ProcessBulkPhoto(ctx context.Context, in *ProcessBulkPhotoRequest, opts ...grpc.CallOption) (*ProcessBulkPhotoResponse, error)
}

type aiServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAiServiceClient(cc grpc.ClientConnInterface) AiServiceClient {
	return &aiServiceClient{cc}
}

func (c *aiServiceClient) ProcessPhoto(ctx context.Context, in *ProcessPhotoRequest, opts ...grpc.CallOption) (*ProcessPhotoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProcessPhotoResponse)
	err := c.cc.Invoke(ctx, AiService_ProcessPhoto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aiServiceClient) ProcessFacecam(ctx context.Context, in *ProcessFacecamRequest, opts ...grpc.CallOption) (*ProcessFacecamResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProcessFacecamResponse)
	err := c.cc.Invoke(ctx, AiService_ProcessFacecam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aiServiceClient) ProcessBulkPhoto(ctx context.Context, in *ProcessBulkPhotoRequest, opts ...grpc.CallOption) (*ProcessBulkPhotoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProcessBulkPhotoResponse)
	err := c.cc.Invoke(ctx, AiService_ProcessBulkPhoto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AiServiceServer is the server API for AiService service.
// All implementations must embed UnimplementedAiServiceServer
// for forward compatibility.
type AiServiceServer interface {
	ProcessPhoto(context.Context, *ProcessPhotoRequest) (*ProcessPhotoResponse, error)
	ProcessFacecam(context.Context, *ProcessFacecamRequest) (*ProcessFacecamResponse, error)
	ProcessBulkPhoto(context.Context, *ProcessBulkPhotoRequest) (*ProcessBulkPhotoResponse, error)
	mustEmbedUnimplementedAiServiceServer()
}

// UnimplementedAiServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAiServiceServer struct{}

func (UnimplementedAiServiceServer) ProcessPhoto(context.Context, *ProcessPhotoRequest) (*ProcessPhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessPhoto not implemented")
}
func (UnimplementedAiServiceServer) ProcessFacecam(context.Context, *ProcessFacecamRequest) (*ProcessFacecamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessFacecam not implemented")
}
func (UnimplementedAiServiceServer) ProcessBulkPhoto(context.Context, *ProcessBulkPhotoRequest) (*ProcessBulkPhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessBulkPhoto not implemented")
}
func (UnimplementedAiServiceServer) mustEmbedUnimplementedAiServiceServer() {}
func (UnimplementedAiServiceServer) testEmbeddedByValue()                   {}

// UnsafeAiServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AiServiceServer will
// result in compilation errors.
type UnsafeAiServiceServer interface {
	mustEmbedUnimplementedAiServiceServer()
}

func RegisterAiServiceServer(s grpc.ServiceRegistrar, srv AiServiceServer) {
	// If the following call pancis, it indicates UnimplementedAiServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AiService_ServiceDesc, srv)
}

func _AiService_ProcessPhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessPhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiServiceServer).ProcessPhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AiService_ProcessPhoto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiServiceServer).ProcessPhoto(ctx, req.(*ProcessPhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiService_ProcessFacecam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessFacecamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiServiceServer).ProcessFacecam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AiService_ProcessFacecam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiServiceServer).ProcessFacecam(ctx, req.(*ProcessFacecamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AiService_ProcessBulkPhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessBulkPhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiServiceServer).ProcessBulkPhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AiService_ProcessBulkPhoto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiServiceServer).ProcessBulkPhoto(ctx, req.(*ProcessBulkPhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AiService_ServiceDesc is the grpc.ServiceDesc for AiService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AiService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ai.AiService",
	HandlerType: (*AiServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessPhoto",
			Handler:    _AiService_ProcessPhoto_Handler,
		},
		{
			MethodName: "ProcessFacecam",
			Handler:    _AiService_ProcessFacecam_Handler,
		},
		{
			MethodName: "ProcessBulkPhoto",
			Handler:    _AiService_ProcessBulkPhoto_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/ai/ai.proto",
}

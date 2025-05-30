// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: pb/photo/photo.proto

package photopb

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
	PhotoService_UpdatePhotographerPhoto_FullMethodName     = "/photo.PhotoService/UpdatePhotographerPhoto"
	PhotoService_UpdateFaceRecogPhoto_FullMethodName        = "/photo.PhotoService/UpdateFaceRecogPhoto"
	PhotoService_CreatePhoto_FullMethodName                 = "/photo.PhotoService/CreatePhoto"
	PhotoService_CreateUserSimilarFacecam_FullMethodName    = "/photo.PhotoService/CreateUserSimilarFacecam"
	PhotoService_CreateFacecam_FullMethodName               = "/photo.PhotoService/CreateFacecam"
	PhotoService_UpdatePhotoDetail_FullMethodName           = "/photo.PhotoService/UpdatePhotoDetail"
	PhotoService_CreateUserSimilar_FullMethodName           = "/photo.PhotoService/CreateUserSimilar"
	PhotoService_CreateCreator_FullMethodName               = "/photo.PhotoService/CreateCreator"
	PhotoService_GetCreator_FullMethodName                  = "/photo.PhotoService/GetCreator"
	PhotoService_CalculatePhotoPrice_FullMethodName         = "/photo.PhotoService/CalculatePhotoPrice"
	PhotoService_OwnerOwnPhotos_FullMethodName              = "/photo.PhotoService/OwnerOwnPhotos"
	PhotoService_CreateBulkPhoto_FullMethodName             = "/photo.PhotoService/CreateBulkPhoto"
	PhotoService_CreateBulkUserSimilarPhotos_FullMethodName = "/photo.PhotoService/CreateBulkUserSimilarPhotos"
	PhotoService_GetPhotoWithDetails_FullMethodName         = "/photo.PhotoService/GetPhotoWithDetails"
	PhotoService_CancelPhotos_FullMethodName                = "/photo.PhotoService/CancelPhotos"
)

// PhotoServiceClient is the client API for PhotoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PhotoServiceClient interface {
	UpdatePhotographerPhoto(ctx context.Context, in *UpdatePhotographerPhotoRequest, opts ...grpc.CallOption) (*UpdatePhotographerPhotoResponse, error)
	UpdateFaceRecogPhoto(ctx context.Context, in *UpdateFaceRecogPhotoRequest, opts ...grpc.CallOption) (*UpdateFaceRecogPhotoResponse, error)
	CreatePhoto(ctx context.Context, in *CreatePhotoRequest, opts ...grpc.CallOption) (*CreatePhotoResponse, error)
	CreateUserSimilarFacecam(ctx context.Context, in *CreateUserSimilarFacecamRequest, opts ...grpc.CallOption) (*CreateUserSimilarFacecamResponse, error)
	CreateFacecam(ctx context.Context, in *CreateFacecamRequest, opts ...grpc.CallOption) (*CreateFacecamResponse, error)
	UpdatePhotoDetail(ctx context.Context, in *UpdatePhotoDetailRequest, opts ...grpc.CallOption) (*UpdatePhotoDetailResponse, error)
	CreateUserSimilar(ctx context.Context, in *CreateUserSimilarPhotoRequest, opts ...grpc.CallOption) (*CreateUserSimilarPhotoResponse, error)
	CreateCreator(ctx context.Context, in *CreateCreatorRequest, opts ...grpc.CallOption) (*CreateCreatorResponse, error)
	GetCreator(ctx context.Context, in *GetCreatorRequest, opts ...grpc.CallOption) (*GetCreatorResponse, error)
	CalculatePhotoPrice(ctx context.Context, in *CalculatePhotoPriceRequest, opts ...grpc.CallOption) (*CalculatePhotoPriceResponse, error)
	OwnerOwnPhotos(ctx context.Context, in *OwnerOwnPhotosRequest, opts ...grpc.CallOption) (*OwnerOwnPhotosResponse, error)
	CreateBulkPhoto(ctx context.Context, in *CreateBulkPhotoRequest, opts ...grpc.CallOption) (*CreateBulkPhotoResponse, error)
	CreateBulkUserSimilarPhotos(ctx context.Context, in *CreateBulkUserSimilarPhotoRequest, opts ...grpc.CallOption) (*CreateBulkUserSimilarPhotoResponse, error)
	GetPhotoWithDetails(ctx context.Context, in *GetPhotoWithDetailsRequest, opts ...grpc.CallOption) (*GetPhotoWithDetailsResponse, error)
	CancelPhotos(ctx context.Context, in *CancelPhotosRequest, opts ...grpc.CallOption) (*CancelPhotosResponse, error)
}

type photoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPhotoServiceClient(cc grpc.ClientConnInterface) PhotoServiceClient {
	return &photoServiceClient{cc}
}

func (c *photoServiceClient) UpdatePhotographerPhoto(ctx context.Context, in *UpdatePhotographerPhotoRequest, opts ...grpc.CallOption) (*UpdatePhotographerPhotoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdatePhotographerPhotoResponse)
	err := c.cc.Invoke(ctx, PhotoService_UpdatePhotographerPhoto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) UpdateFaceRecogPhoto(ctx context.Context, in *UpdateFaceRecogPhotoRequest, opts ...grpc.CallOption) (*UpdateFaceRecogPhotoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateFaceRecogPhotoResponse)
	err := c.cc.Invoke(ctx, PhotoService_UpdateFaceRecogPhoto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) CreatePhoto(ctx context.Context, in *CreatePhotoRequest, opts ...grpc.CallOption) (*CreatePhotoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreatePhotoResponse)
	err := c.cc.Invoke(ctx, PhotoService_CreatePhoto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) CreateUserSimilarFacecam(ctx context.Context, in *CreateUserSimilarFacecamRequest, opts ...grpc.CallOption) (*CreateUserSimilarFacecamResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateUserSimilarFacecamResponse)
	err := c.cc.Invoke(ctx, PhotoService_CreateUserSimilarFacecam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) CreateFacecam(ctx context.Context, in *CreateFacecamRequest, opts ...grpc.CallOption) (*CreateFacecamResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateFacecamResponse)
	err := c.cc.Invoke(ctx, PhotoService_CreateFacecam_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) UpdatePhotoDetail(ctx context.Context, in *UpdatePhotoDetailRequest, opts ...grpc.CallOption) (*UpdatePhotoDetailResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdatePhotoDetailResponse)
	err := c.cc.Invoke(ctx, PhotoService_UpdatePhotoDetail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) CreateUserSimilar(ctx context.Context, in *CreateUserSimilarPhotoRequest, opts ...grpc.CallOption) (*CreateUserSimilarPhotoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateUserSimilarPhotoResponse)
	err := c.cc.Invoke(ctx, PhotoService_CreateUserSimilar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) CreateCreator(ctx context.Context, in *CreateCreatorRequest, opts ...grpc.CallOption) (*CreateCreatorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCreatorResponse)
	err := c.cc.Invoke(ctx, PhotoService_CreateCreator_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) GetCreator(ctx context.Context, in *GetCreatorRequest, opts ...grpc.CallOption) (*GetCreatorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCreatorResponse)
	err := c.cc.Invoke(ctx, PhotoService_GetCreator_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) CalculatePhotoPrice(ctx context.Context, in *CalculatePhotoPriceRequest, opts ...grpc.CallOption) (*CalculatePhotoPriceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CalculatePhotoPriceResponse)
	err := c.cc.Invoke(ctx, PhotoService_CalculatePhotoPrice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) OwnerOwnPhotos(ctx context.Context, in *OwnerOwnPhotosRequest, opts ...grpc.CallOption) (*OwnerOwnPhotosResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OwnerOwnPhotosResponse)
	err := c.cc.Invoke(ctx, PhotoService_OwnerOwnPhotos_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) CreateBulkPhoto(ctx context.Context, in *CreateBulkPhotoRequest, opts ...grpc.CallOption) (*CreateBulkPhotoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateBulkPhotoResponse)
	err := c.cc.Invoke(ctx, PhotoService_CreateBulkPhoto_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) CreateBulkUserSimilarPhotos(ctx context.Context, in *CreateBulkUserSimilarPhotoRequest, opts ...grpc.CallOption) (*CreateBulkUserSimilarPhotoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateBulkUserSimilarPhotoResponse)
	err := c.cc.Invoke(ctx, PhotoService_CreateBulkUserSimilarPhotos_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) GetPhotoWithDetails(ctx context.Context, in *GetPhotoWithDetailsRequest, opts ...grpc.CallOption) (*GetPhotoWithDetailsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPhotoWithDetailsResponse)
	err := c.cc.Invoke(ctx, PhotoService_GetPhotoWithDetails_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *photoServiceClient) CancelPhotos(ctx context.Context, in *CancelPhotosRequest, opts ...grpc.CallOption) (*CancelPhotosResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CancelPhotosResponse)
	err := c.cc.Invoke(ctx, PhotoService_CancelPhotos_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PhotoServiceServer is the server API for PhotoService service.
// All implementations must embed UnimplementedPhotoServiceServer
// for forward compatibility.
type PhotoServiceServer interface {
	UpdatePhotographerPhoto(context.Context, *UpdatePhotographerPhotoRequest) (*UpdatePhotographerPhotoResponse, error)
	UpdateFaceRecogPhoto(context.Context, *UpdateFaceRecogPhotoRequest) (*UpdateFaceRecogPhotoResponse, error)
	CreatePhoto(context.Context, *CreatePhotoRequest) (*CreatePhotoResponse, error)
	CreateUserSimilarFacecam(context.Context, *CreateUserSimilarFacecamRequest) (*CreateUserSimilarFacecamResponse, error)
	CreateFacecam(context.Context, *CreateFacecamRequest) (*CreateFacecamResponse, error)
	UpdatePhotoDetail(context.Context, *UpdatePhotoDetailRequest) (*UpdatePhotoDetailResponse, error)
	CreateUserSimilar(context.Context, *CreateUserSimilarPhotoRequest) (*CreateUserSimilarPhotoResponse, error)
	CreateCreator(context.Context, *CreateCreatorRequest) (*CreateCreatorResponse, error)
	GetCreator(context.Context, *GetCreatorRequest) (*GetCreatorResponse, error)
	CalculatePhotoPrice(context.Context, *CalculatePhotoPriceRequest) (*CalculatePhotoPriceResponse, error)
	OwnerOwnPhotos(context.Context, *OwnerOwnPhotosRequest) (*OwnerOwnPhotosResponse, error)
	CreateBulkPhoto(context.Context, *CreateBulkPhotoRequest) (*CreateBulkPhotoResponse, error)
	CreateBulkUserSimilarPhotos(context.Context, *CreateBulkUserSimilarPhotoRequest) (*CreateBulkUserSimilarPhotoResponse, error)
	GetPhotoWithDetails(context.Context, *GetPhotoWithDetailsRequest) (*GetPhotoWithDetailsResponse, error)
	CancelPhotos(context.Context, *CancelPhotosRequest) (*CancelPhotosResponse, error)
	mustEmbedUnimplementedPhotoServiceServer()
}

// UnimplementedPhotoServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPhotoServiceServer struct{}

func (UnimplementedPhotoServiceServer) UpdatePhotographerPhoto(context.Context, *UpdatePhotographerPhotoRequest) (*UpdatePhotographerPhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePhotographerPhoto not implemented")
}
func (UnimplementedPhotoServiceServer) UpdateFaceRecogPhoto(context.Context, *UpdateFaceRecogPhotoRequest) (*UpdateFaceRecogPhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFaceRecogPhoto not implemented")
}
func (UnimplementedPhotoServiceServer) CreatePhoto(context.Context, *CreatePhotoRequest) (*CreatePhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePhoto not implemented")
}
func (UnimplementedPhotoServiceServer) CreateUserSimilarFacecam(context.Context, *CreateUserSimilarFacecamRequest) (*CreateUserSimilarFacecamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserSimilarFacecam not implemented")
}
func (UnimplementedPhotoServiceServer) CreateFacecam(context.Context, *CreateFacecamRequest) (*CreateFacecamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFacecam not implemented")
}
func (UnimplementedPhotoServiceServer) UpdatePhotoDetail(context.Context, *UpdatePhotoDetailRequest) (*UpdatePhotoDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePhotoDetail not implemented")
}
func (UnimplementedPhotoServiceServer) CreateUserSimilar(context.Context, *CreateUserSimilarPhotoRequest) (*CreateUserSimilarPhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserSimilar not implemented")
}
func (UnimplementedPhotoServiceServer) CreateCreator(context.Context, *CreateCreatorRequest) (*CreateCreatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCreator not implemented")
}
func (UnimplementedPhotoServiceServer) GetCreator(context.Context, *GetCreatorRequest) (*GetCreatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCreator not implemented")
}
func (UnimplementedPhotoServiceServer) CalculatePhotoPrice(context.Context, *CalculatePhotoPriceRequest) (*CalculatePhotoPriceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculatePhotoPrice not implemented")
}
func (UnimplementedPhotoServiceServer) OwnerOwnPhotos(context.Context, *OwnerOwnPhotosRequest) (*OwnerOwnPhotosResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OwnerOwnPhotos not implemented")
}
func (UnimplementedPhotoServiceServer) CreateBulkPhoto(context.Context, *CreateBulkPhotoRequest) (*CreateBulkPhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBulkPhoto not implemented")
}
func (UnimplementedPhotoServiceServer) CreateBulkUserSimilarPhotos(context.Context, *CreateBulkUserSimilarPhotoRequest) (*CreateBulkUserSimilarPhotoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBulkUserSimilarPhotos not implemented")
}
func (UnimplementedPhotoServiceServer) GetPhotoWithDetails(context.Context, *GetPhotoWithDetailsRequest) (*GetPhotoWithDetailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPhotoWithDetails not implemented")
}
func (UnimplementedPhotoServiceServer) CancelPhotos(context.Context, *CancelPhotosRequest) (*CancelPhotosResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelPhotos not implemented")
}
func (UnimplementedPhotoServiceServer) mustEmbedUnimplementedPhotoServiceServer() {}
func (UnimplementedPhotoServiceServer) testEmbeddedByValue()                      {}

// UnsafePhotoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PhotoServiceServer will
// result in compilation errors.
type UnsafePhotoServiceServer interface {
	mustEmbedUnimplementedPhotoServiceServer()
}

func RegisterPhotoServiceServer(s grpc.ServiceRegistrar, srv PhotoServiceServer) {
	// If the following call pancis, it indicates UnimplementedPhotoServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PhotoService_ServiceDesc, srv)
}

func _PhotoService_UpdatePhotographerPhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePhotographerPhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).UpdatePhotographerPhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_UpdatePhotographerPhoto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).UpdatePhotographerPhoto(ctx, req.(*UpdatePhotographerPhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_UpdateFaceRecogPhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateFaceRecogPhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).UpdateFaceRecogPhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_UpdateFaceRecogPhoto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).UpdateFaceRecogPhoto(ctx, req.(*UpdateFaceRecogPhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_CreatePhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).CreatePhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_CreatePhoto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).CreatePhoto(ctx, req.(*CreatePhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_CreateUserSimilarFacecam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserSimilarFacecamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).CreateUserSimilarFacecam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_CreateUserSimilarFacecam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).CreateUserSimilarFacecam(ctx, req.(*CreateUserSimilarFacecamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_CreateFacecam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFacecamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).CreateFacecam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_CreateFacecam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).CreateFacecam(ctx, req.(*CreateFacecamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_UpdatePhotoDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePhotoDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).UpdatePhotoDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_UpdatePhotoDetail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).UpdatePhotoDetail(ctx, req.(*UpdatePhotoDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_CreateUserSimilar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserSimilarPhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).CreateUserSimilar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_CreateUserSimilar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).CreateUserSimilar(ctx, req.(*CreateUserSimilarPhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_CreateCreator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCreatorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).CreateCreator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_CreateCreator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).CreateCreator(ctx, req.(*CreateCreatorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_GetCreator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCreatorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).GetCreator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_GetCreator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).GetCreator(ctx, req.(*GetCreatorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_CalculatePhotoPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalculatePhotoPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).CalculatePhotoPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_CalculatePhotoPrice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).CalculatePhotoPrice(ctx, req.(*CalculatePhotoPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_OwnerOwnPhotos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OwnerOwnPhotosRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).OwnerOwnPhotos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_OwnerOwnPhotos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).OwnerOwnPhotos(ctx, req.(*OwnerOwnPhotosRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_CreateBulkPhoto_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBulkPhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).CreateBulkPhoto(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_CreateBulkPhoto_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).CreateBulkPhoto(ctx, req.(*CreateBulkPhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_CreateBulkUserSimilarPhotos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBulkUserSimilarPhotoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).CreateBulkUserSimilarPhotos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_CreateBulkUserSimilarPhotos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).CreateBulkUserSimilarPhotos(ctx, req.(*CreateBulkUserSimilarPhotoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_GetPhotoWithDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPhotoWithDetailsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).GetPhotoWithDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_GetPhotoWithDetails_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).GetPhotoWithDetails(ctx, req.(*GetPhotoWithDetailsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PhotoService_CancelPhotos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelPhotosRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PhotoServiceServer).CancelPhotos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PhotoService_CancelPhotos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PhotoServiceServer).CancelPhotos(ctx, req.(*CancelPhotosRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PhotoService_ServiceDesc is the grpc.ServiceDesc for PhotoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PhotoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "photo.PhotoService",
	HandlerType: (*PhotoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdatePhotographerPhoto",
			Handler:    _PhotoService_UpdatePhotographerPhoto_Handler,
		},
		{
			MethodName: "UpdateFaceRecogPhoto",
			Handler:    _PhotoService_UpdateFaceRecogPhoto_Handler,
		},
		{
			MethodName: "CreatePhoto",
			Handler:    _PhotoService_CreatePhoto_Handler,
		},
		{
			MethodName: "CreateUserSimilarFacecam",
			Handler:    _PhotoService_CreateUserSimilarFacecam_Handler,
		},
		{
			MethodName: "CreateFacecam",
			Handler:    _PhotoService_CreateFacecam_Handler,
		},
		{
			MethodName: "UpdatePhotoDetail",
			Handler:    _PhotoService_UpdatePhotoDetail_Handler,
		},
		{
			MethodName: "CreateUserSimilar",
			Handler:    _PhotoService_CreateUserSimilar_Handler,
		},
		{
			MethodName: "CreateCreator",
			Handler:    _PhotoService_CreateCreator_Handler,
		},
		{
			MethodName: "GetCreator",
			Handler:    _PhotoService_GetCreator_Handler,
		},
		{
			MethodName: "CalculatePhotoPrice",
			Handler:    _PhotoService_CalculatePhotoPrice_Handler,
		},
		{
			MethodName: "OwnerOwnPhotos",
			Handler:    _PhotoService_OwnerOwnPhotos_Handler,
		},
		{
			MethodName: "CreateBulkPhoto",
			Handler:    _PhotoService_CreateBulkPhoto_Handler,
		},
		{
			MethodName: "CreateBulkUserSimilarPhotos",
			Handler:    _PhotoService_CreateBulkUserSimilarPhotos_Handler,
		},
		{
			MethodName: "GetPhotoWithDetails",
			Handler:    _PhotoService_GetPhotoWithDetails_Handler,
		},
		{
			MethodName: "CancelPhotos",
			Handler:    _PhotoService_CancelPhotos_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/photo/photo.proto",
}

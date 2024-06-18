// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.26.1
// source: protos/posttag.proto

package posttag

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	PostTagService_CreatePostTag_FullMethodName  = "/forum.PostTagService/CreatePostTag"
	PostTagService_DeletePostTag_FullMethodName  = "/forum.PostTagService/DeletePostTag"
	PostTagService_GetPostsByTag_FullMethodName  = "/forum.PostTagService/GetPostsByTag"
	PostTagService_GetAllPostTags_FullMethodName = "/forum.PostTagService/GetAllPostTags"
)

// PostTagServiceClient is the client API for PostTagService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostTagServiceClient interface {
	// PostTag CRUD
	CreatePostTag(ctx context.Context, in *CreatePostTagRequest, opts ...grpc.CallOption) (*CreatePostTagResponse, error)
	DeletePostTag(ctx context.Context, in *DeletePostTagRequest, opts ...grpc.CallOption) (*DeletePostTagResponse, error)
	// Get posts associated with a tag
	GetPostsByTag(ctx context.Context, in *GetPostsByTagRequest, opts ...grpc.CallOption) (*GetPostsByTagResponse, error)
	// PostTag GetAll
	GetAllPostTags(ctx context.Context, in *GetAllPostTagsRequest, opts ...grpc.CallOption) (*GetAllPostTagsResponse, error)
}

type postTagServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPostTagServiceClient(cc grpc.ClientConnInterface) PostTagServiceClient {
	return &postTagServiceClient{cc}
}

func (c *postTagServiceClient) CreatePostTag(ctx context.Context, in *CreatePostTagRequest, opts ...grpc.CallOption) (*CreatePostTagResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreatePostTagResponse)
	err := c.cc.Invoke(ctx, PostTagService_CreatePostTag_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postTagServiceClient) DeletePostTag(ctx context.Context, in *DeletePostTagRequest, opts ...grpc.CallOption) (*DeletePostTagResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeletePostTagResponse)
	err := c.cc.Invoke(ctx, PostTagService_DeletePostTag_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postTagServiceClient) GetPostsByTag(ctx context.Context, in *GetPostsByTagRequest, opts ...grpc.CallOption) (*GetPostsByTagResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPostsByTagResponse)
	err := c.cc.Invoke(ctx, PostTagService_GetPostsByTag_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postTagServiceClient) GetAllPostTags(ctx context.Context, in *GetAllPostTagsRequest, opts ...grpc.CallOption) (*GetAllPostTagsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllPostTagsResponse)
	err := c.cc.Invoke(ctx, PostTagService_GetAllPostTags_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostTagServiceServer is the server API for PostTagService service.
// All implementations must embed UnimplementedPostTagServiceServer
// for forward compatibility
type PostTagServiceServer interface {
	// PostTag CRUD
	CreatePostTag(context.Context, *CreatePostTagRequest) (*CreatePostTagResponse, error)
	DeletePostTag(context.Context, *DeletePostTagRequest) (*DeletePostTagResponse, error)
	// Get posts associated with a tag
	GetPostsByTag(context.Context, *GetPostsByTagRequest) (*GetPostsByTagResponse, error)
	// PostTag GetAll
	GetAllPostTags(context.Context, *GetAllPostTagsRequest) (*GetAllPostTagsResponse, error)
	mustEmbedUnimplementedPostTagServiceServer()
}

// UnimplementedPostTagServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPostTagServiceServer struct {
}

func (UnimplementedPostTagServiceServer) CreatePostTag(context.Context, *CreatePostTagRequest) (*CreatePostTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePostTag not implemented")
}
func (UnimplementedPostTagServiceServer) DeletePostTag(context.Context, *DeletePostTagRequest) (*DeletePostTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePostTag not implemented")
}
func (UnimplementedPostTagServiceServer) GetPostsByTag(context.Context, *GetPostsByTagRequest) (*GetPostsByTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPostsByTag not implemented")
}
func (UnimplementedPostTagServiceServer) GetAllPostTags(context.Context, *GetAllPostTagsRequest) (*GetAllPostTagsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllPostTags not implemented")
}
func (UnimplementedPostTagServiceServer) mustEmbedUnimplementedPostTagServiceServer() {}

// UnsafePostTagServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostTagServiceServer will
// result in compilation errors.
type UnsafePostTagServiceServer interface {
	mustEmbedUnimplementedPostTagServiceServer()
}

func RegisterPostTagServiceServer(s grpc.ServiceRegistrar, srv PostTagServiceServer) {
	s.RegisterService(&PostTagService_ServiceDesc, srv)
}

func _PostTagService_CreatePostTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePostTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostTagServiceServer).CreatePostTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostTagService_CreatePostTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostTagServiceServer).CreatePostTag(ctx, req.(*CreatePostTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostTagService_DeletePostTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePostTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostTagServiceServer).DeletePostTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostTagService_DeletePostTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostTagServiceServer).DeletePostTag(ctx, req.(*DeletePostTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostTagService_GetPostsByTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostsByTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostTagServiceServer).GetPostsByTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostTagService_GetPostsByTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostTagServiceServer).GetPostsByTag(ctx, req.(*GetPostsByTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostTagService_GetAllPostTags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllPostTagsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostTagServiceServer).GetAllPostTags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostTagService_GetAllPostTags_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostTagServiceServer).GetAllPostTags(ctx, req.(*GetAllPostTagsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PostTagService_ServiceDesc is the grpc.ServiceDesc for PostTagService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostTagService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "forum.PostTagService",
	HandlerType: (*PostTagServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePostTag",
			Handler:    _PostTagService_CreatePostTag_Handler,
		},
		{
			MethodName: "DeletePostTag",
			Handler:    _PostTagService_DeletePostTag_Handler,
		},
		{
			MethodName: "GetPostsByTag",
			Handler:    _PostTagService_GetPostsByTag_Handler,
		},
		{
			MethodName: "GetAllPostTags",
			Handler:    _PostTagService_GetAllPostTags_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/posttag.proto",
}
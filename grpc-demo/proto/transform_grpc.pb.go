// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: proto/transform.proto

package proto

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

// TransformTextServiceClient is the client API for TransformTextService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransformTextServiceClient interface {
	// Unary rpc
	TransformText(ctx context.Context, in *TransformTextRequest, opts ...grpc.CallOption) (*TransformTextResponse, error)
	// BiDi rpc
	TransformAndSplitText(ctx context.Context, opts ...grpc.CallOption) (TransformTextService_TransformAndSplitTextClient, error)
}

type transformTextServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTransformTextServiceClient(cc grpc.ClientConnInterface) TransformTextServiceClient {
	return &transformTextServiceClient{cc}
}

func (c *transformTextServiceClient) TransformText(ctx context.Context, in *TransformTextRequest, opts ...grpc.CallOption) (*TransformTextResponse, error) {
	out := new(TransformTextResponse)
	err := c.cc.Invoke(ctx, "/transform.TransformTextService/TransformText", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transformTextServiceClient) TransformAndSplitText(ctx context.Context, opts ...grpc.CallOption) (TransformTextService_TransformAndSplitTextClient, error) {
	stream, err := c.cc.NewStream(ctx, &TransformTextService_ServiceDesc.Streams[0], "/transform.TransformTextService/TransformAndSplitText", opts...)
	if err != nil {
		return nil, err
	}
	x := &transformTextServiceTransformAndSplitTextClient{stream}
	return x, nil
}

type TransformTextService_TransformAndSplitTextClient interface {
	Send(*TransformTextRequest) error
	Recv() (*TransformTextResponse, error)
	grpc.ClientStream
}

type transformTextServiceTransformAndSplitTextClient struct {
	grpc.ClientStream
}

func (x *transformTextServiceTransformAndSplitTextClient) Send(m *TransformTextRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *transformTextServiceTransformAndSplitTextClient) Recv() (*TransformTextResponse, error) {
	m := new(TransformTextResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TransformTextServiceServer is the server API for TransformTextService service.
// All implementations must embed UnimplementedTransformTextServiceServer
// for forward compatibility
type TransformTextServiceServer interface {
	// Unary rpc
	TransformText(context.Context, *TransformTextRequest) (*TransformTextResponse, error)
	// BiDi rpc
	TransformAndSplitText(TransformTextService_TransformAndSplitTextServer) error
	mustEmbedUnimplementedTransformTextServiceServer()
}

// UnimplementedTransformTextServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTransformTextServiceServer struct {
}

func (UnimplementedTransformTextServiceServer) TransformText(context.Context, *TransformTextRequest) (*TransformTextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransformText not implemented")
}
func (UnimplementedTransformTextServiceServer) TransformAndSplitText(TransformTextService_TransformAndSplitTextServer) error {
	return status.Errorf(codes.Unimplemented, "method TransformAndSplitText not implemented")
}
func (UnimplementedTransformTextServiceServer) mustEmbedUnimplementedTransformTextServiceServer() {}

// UnsafeTransformTextServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransformTextServiceServer will
// result in compilation errors.
type UnsafeTransformTextServiceServer interface {
	mustEmbedUnimplementedTransformTextServiceServer()
}

func RegisterTransformTextServiceServer(s grpc.ServiceRegistrar, srv TransformTextServiceServer) {
	s.RegisterService(&TransformTextService_ServiceDesc, srv)
}

func _TransformTextService_TransformText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransformTextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransformTextServiceServer).TransformText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transform.TransformTextService/TransformText",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransformTextServiceServer).TransformText(ctx, req.(*TransformTextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TransformTextService_TransformAndSplitText_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TransformTextServiceServer).TransformAndSplitText(&transformTextServiceTransformAndSplitTextServer{stream})
}

type TransformTextService_TransformAndSplitTextServer interface {
	Send(*TransformTextResponse) error
	Recv() (*TransformTextRequest, error)
	grpc.ServerStream
}

type transformTextServiceTransformAndSplitTextServer struct {
	grpc.ServerStream
}

func (x *transformTextServiceTransformAndSplitTextServer) Send(m *TransformTextResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *transformTextServiceTransformAndSplitTextServer) Recv() (*TransformTextRequest, error) {
	m := new(TransformTextRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TransformTextService_ServiceDesc is the grpc.ServiceDesc for TransformTextService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TransformTextService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "transform.TransformTextService",
	HandlerType: (*TransformTextServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TransformText",
			Handler:    _TransformTextService_TransformText_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "TransformAndSplitText",
			Handler:       _TransformTextService_TransformAndSplitText_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/transform.proto",
}

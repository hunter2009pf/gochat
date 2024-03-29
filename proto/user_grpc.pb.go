// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: user.proto

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

// CommunicationClient is the client API for Communication service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommunicationClient interface {
	ConnectServer(ctx context.Context, in *User, opts ...grpc.CallOption) (Communication_ConnectServerClient, error)
	SendMsg(ctx context.Context, in *SendMsgRequest, opts ...grpc.CallOption) (*SendMsgResponse, error)
}

type communicationClient struct {
	cc grpc.ClientConnInterface
}

func NewCommunicationClient(cc grpc.ClientConnInterface) CommunicationClient {
	return &communicationClient{cc}
}

func (c *communicationClient) ConnectServer(ctx context.Context, in *User, opts ...grpc.CallOption) (Communication_ConnectServerClient, error) {
	stream, err := c.cc.NewStream(ctx, &Communication_ServiceDesc.Streams[0], "/user.Communication/ConnectServer", opts...)
	if err != nil {
		return nil, err
	}
	x := &communicationConnectServerClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Communication_ConnectServerClient interface {
	Recv() (*TextMsg, error)
	grpc.ClientStream
}

type communicationConnectServerClient struct {
	grpc.ClientStream
}

func (x *communicationConnectServerClient) Recv() (*TextMsg, error) {
	m := new(TextMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *communicationClient) SendMsg(ctx context.Context, in *SendMsgRequest, opts ...grpc.CallOption) (*SendMsgResponse, error) {
	out := new(SendMsgResponse)
	err := c.cc.Invoke(ctx, "/user.Communication/SendMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommunicationServer is the server API for Communication service.
// All implementations must embed UnimplementedCommunicationServer
// for forward compatibility
type CommunicationServer interface {
	ConnectServer(*User, Communication_ConnectServerServer) error
	SendMsg(context.Context, *SendMsgRequest) (*SendMsgResponse, error)
	mustEmbedUnimplementedCommunicationServer()
}

// UnimplementedCommunicationServer must be embedded to have forward compatible implementations.
type UnimplementedCommunicationServer struct {
}

func (UnimplementedCommunicationServer) ConnectServer(*User, Communication_ConnectServerServer) error {
	return status.Errorf(codes.Unimplemented, "method ConnectServer not implemented")
}
func (UnimplementedCommunicationServer) SendMsg(context.Context, *SendMsgRequest) (*SendMsgResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsg not implemented")
}
func (UnimplementedCommunicationServer) mustEmbedUnimplementedCommunicationServer() {}

// UnsafeCommunicationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommunicationServer will
// result in compilation errors.
type UnsafeCommunicationServer interface {
	mustEmbedUnimplementedCommunicationServer()
}

func RegisterCommunicationServer(s grpc.ServiceRegistrar, srv CommunicationServer) {
	s.RegisterService(&Communication_ServiceDesc, srv)
}

func _Communication_ConnectServer_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(User)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CommunicationServer).ConnectServer(m, &communicationConnectServerServer{stream})
}

type Communication_ConnectServerServer interface {
	Send(*TextMsg) error
	grpc.ServerStream
}

type communicationConnectServerServer struct {
	grpc.ServerStream
}

func (x *communicationConnectServerServer) Send(m *TextMsg) error {
	return x.ServerStream.SendMsg(m)
}

func _Communication_SendMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMsgRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicationServer).SendMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.Communication/SendMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicationServer).SendMsg(ctx, req.(*SendMsgRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Communication_ServiceDesc is the grpc.ServiceDesc for Communication service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Communication_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.Communication",
	HandlerType: (*CommunicationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMsg",
			Handler:    _Communication_SendMsg_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ConnectServer",
			Handler:       _Communication_ConnectServer_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "user.proto",
}

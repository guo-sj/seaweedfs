// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package messaging_pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// SeaweedMessagingClient is the client API for SeaweedMessaging service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SeaweedMessagingClient interface {
	Subscribe(ctx context.Context, opts ...grpc.CallOption) (SeaweedMessaging_SubscribeClient, error)
	Publish(ctx context.Context, opts ...grpc.CallOption) (SeaweedMessaging_PublishClient, error)
	DeleteTopic(ctx context.Context, in *DeleteTopicRequest, opts ...grpc.CallOption) (*DeleteTopicResponse, error)
	ConfigureTopic(ctx context.Context, in *ConfigureTopicRequest, opts ...grpc.CallOption) (*ConfigureTopicResponse, error)
	GetTopicConfiguration(ctx context.Context, in *GetTopicConfigurationRequest, opts ...grpc.CallOption) (*GetTopicConfigurationResponse, error)
	FindBroker(ctx context.Context, in *FindBrokerRequest, opts ...grpc.CallOption) (*FindBrokerResponse, error)
}

type seaweedMessagingClient struct {
	cc grpc.ClientConnInterface
}

func NewSeaweedMessagingClient(cc grpc.ClientConnInterface) SeaweedMessagingClient {
	return &seaweedMessagingClient{cc}
}

func (c *seaweedMessagingClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (SeaweedMessaging_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &SeaweedMessaging_ServiceDesc.Streams[0], "/messaging_pb.SeaweedMessaging/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &seaweedMessagingSubscribeClient{stream}
	return x, nil
}

type SeaweedMessaging_SubscribeClient interface {
	Send(*SubscriberMessage) error
	Recv() (*BrokerMessage, error)
	grpc.ClientStream
}

type seaweedMessagingSubscribeClient struct {
	grpc.ClientStream
}

func (x *seaweedMessagingSubscribeClient) Send(m *SubscriberMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *seaweedMessagingSubscribeClient) Recv() (*BrokerMessage, error) {
	m := new(BrokerMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *seaweedMessagingClient) Publish(ctx context.Context, opts ...grpc.CallOption) (SeaweedMessaging_PublishClient, error) {
	stream, err := c.cc.NewStream(ctx, &SeaweedMessaging_ServiceDesc.Streams[1], "/messaging_pb.SeaweedMessaging/Publish", opts...)
	if err != nil {
		return nil, err
	}
	x := &seaweedMessagingPublishClient{stream}
	return x, nil
}

type SeaweedMessaging_PublishClient interface {
	Send(*PublishRequest) error
	Recv() (*PublishResponse, error)
	grpc.ClientStream
}

type seaweedMessagingPublishClient struct {
	grpc.ClientStream
}

func (x *seaweedMessagingPublishClient) Send(m *PublishRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *seaweedMessagingPublishClient) Recv() (*PublishResponse, error) {
	m := new(PublishResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *seaweedMessagingClient) DeleteTopic(ctx context.Context, in *DeleteTopicRequest, opts ...grpc.CallOption) (*DeleteTopicResponse, error) {
	out := new(DeleteTopicResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/DeleteTopic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedMessagingClient) ConfigureTopic(ctx context.Context, in *ConfigureTopicRequest, opts ...grpc.CallOption) (*ConfigureTopicResponse, error) {
	out := new(ConfigureTopicResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/ConfigureTopic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedMessagingClient) GetTopicConfiguration(ctx context.Context, in *GetTopicConfigurationRequest, opts ...grpc.CallOption) (*GetTopicConfigurationResponse, error) {
	out := new(GetTopicConfigurationResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/GetTopicConfiguration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *seaweedMessagingClient) FindBroker(ctx context.Context, in *FindBrokerRequest, opts ...grpc.CallOption) (*FindBrokerResponse, error) {
	out := new(FindBrokerResponse)
	err := c.cc.Invoke(ctx, "/messaging_pb.SeaweedMessaging/FindBroker", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SeaweedMessagingServer is the server API for SeaweedMessaging service.
// All implementations must embed UnimplementedSeaweedMessagingServer
// for forward compatibility
type SeaweedMessagingServer interface {
	Subscribe(SeaweedMessaging_SubscribeServer) error
	Publish(SeaweedMessaging_PublishServer) error
	DeleteTopic(context.Context, *DeleteTopicRequest) (*DeleteTopicResponse, error)
	ConfigureTopic(context.Context, *ConfigureTopicRequest) (*ConfigureTopicResponse, error)
	GetTopicConfiguration(context.Context, *GetTopicConfigurationRequest) (*GetTopicConfigurationResponse, error)
	FindBroker(context.Context, *FindBrokerRequest) (*FindBrokerResponse, error)
	mustEmbedUnimplementedSeaweedMessagingServer()
}

// UnimplementedSeaweedMessagingServer must be embedded to have forward compatible implementations.
type UnimplementedSeaweedMessagingServer struct {
}

func (UnimplementedSeaweedMessagingServer) Subscribe(SeaweedMessaging_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedSeaweedMessagingServer) Publish(SeaweedMessaging_PublishServer) error {
	return status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedSeaweedMessagingServer) DeleteTopic(context.Context, *DeleteTopicRequest) (*DeleteTopicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTopic not implemented")
}
func (UnimplementedSeaweedMessagingServer) ConfigureTopic(context.Context, *ConfigureTopicRequest) (*ConfigureTopicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigureTopic not implemented")
}
func (UnimplementedSeaweedMessagingServer) GetTopicConfiguration(context.Context, *GetTopicConfigurationRequest) (*GetTopicConfigurationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopicConfiguration not implemented")
}
func (UnimplementedSeaweedMessagingServer) FindBroker(context.Context, *FindBrokerRequest) (*FindBrokerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindBroker not implemented")
}
func (UnimplementedSeaweedMessagingServer) mustEmbedUnimplementedSeaweedMessagingServer() {}

// UnsafeSeaweedMessagingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SeaweedMessagingServer will
// result in compilation errors.
type UnsafeSeaweedMessagingServer interface {
	mustEmbedUnimplementedSeaweedMessagingServer()
}

func RegisterSeaweedMessagingServer(s grpc.ServiceRegistrar, srv SeaweedMessagingServer) {
	s.RegisterService(&SeaweedMessaging_ServiceDesc, srv)
}

func _SeaweedMessaging_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SeaweedMessagingServer).Subscribe(&seaweedMessagingSubscribeServer{stream})
}

type SeaweedMessaging_SubscribeServer interface {
	Send(*BrokerMessage) error
	Recv() (*SubscriberMessage, error)
	grpc.ServerStream
}

type seaweedMessagingSubscribeServer struct {
	grpc.ServerStream
}

func (x *seaweedMessagingSubscribeServer) Send(m *BrokerMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *seaweedMessagingSubscribeServer) Recv() (*SubscriberMessage, error) {
	m := new(SubscriberMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SeaweedMessaging_Publish_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SeaweedMessagingServer).Publish(&seaweedMessagingPublishServer{stream})
}

type SeaweedMessaging_PublishServer interface {
	Send(*PublishResponse) error
	Recv() (*PublishRequest, error)
	grpc.ServerStream
}

type seaweedMessagingPublishServer struct {
	grpc.ServerStream
}

func (x *seaweedMessagingPublishServer) Send(m *PublishResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *seaweedMessagingPublishServer) Recv() (*PublishRequest, error) {
	m := new(PublishRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SeaweedMessaging_DeleteTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTopicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).DeleteTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/DeleteTopic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).DeleteTopic(ctx, req.(*DeleteTopicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SeaweedMessaging_ConfigureTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureTopicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).ConfigureTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/ConfigureTopic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).ConfigureTopic(ctx, req.(*ConfigureTopicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SeaweedMessaging_GetTopicConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTopicConfigurationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).GetTopicConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/GetTopicConfiguration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).GetTopicConfiguration(ctx, req.(*GetTopicConfigurationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SeaweedMessaging_FindBroker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindBrokerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeaweedMessagingServer).FindBroker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messaging_pb.SeaweedMessaging/FindBroker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeaweedMessagingServer).FindBroker(ctx, req.(*FindBrokerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SeaweedMessaging_ServiceDesc is the grpc.ServiceDesc for SeaweedMessaging service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SeaweedMessaging_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "messaging_pb.SeaweedMessaging",
	HandlerType: (*SeaweedMessagingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteTopic",
			Handler:    _SeaweedMessaging_DeleteTopic_Handler,
		},
		{
			MethodName: "ConfigureTopic",
			Handler:    _SeaweedMessaging_ConfigureTopic_Handler,
		},
		{
			MethodName: "GetTopicConfiguration",
			Handler:    _SeaweedMessaging_GetTopicConfiguration_Handler,
		},
		{
			MethodName: "FindBroker",
			Handler:    _SeaweedMessaging_FindBroker_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _SeaweedMessaging_Subscribe_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Publish",
			Handler:       _SeaweedMessaging_Publish_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "messaging.proto",
}

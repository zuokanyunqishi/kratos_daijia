// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: api/customer/customer.proto

///kratos:clients

package customer

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

const (
	Customer_GetCustomer_FullMethodName = "/api.customer.Customer/GetCustomer"
	Customer_Login_FullMethodName       = "/api.customer.Customer/Login"
)

// CustomerClient is the client API for Customer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CustomerClient interface {
	// rpc CreateCustomer (CreateCustomerRequest) returns (CreateCustomerReply);
	// rpc UpdateCustomer (UpdateCustomerRequest) returns (UpdateCustomerReply);
	// rpc DeleteCustomer (DeleteCustomerRequest) returns (DeleteCustomerReply);
	GetCustomer(ctx context.Context, in *GetCustomerRequest, opts ...grpc.CallOption) (*GetCustomerReply, error)
	// 登录
	Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRes, error)
}

type customerClient struct {
	cc grpc.ClientConnInterface
}

func NewCustomerClient(cc grpc.ClientConnInterface) CustomerClient {
	return &customerClient{cc}
}

func (c *customerClient) GetCustomer(ctx context.Context, in *GetCustomerRequest, opts ...grpc.CallOption) (*GetCustomerReply, error) {
	out := new(GetCustomerReply)
	err := c.cc.Invoke(ctx, Customer_GetCustomer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *customerClient) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRes, error) {
	out := new(LoginRes)
	err := c.cc.Invoke(ctx, Customer_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CustomerServer is the server API for Customer service.
// All implementations must embed UnimplementedCustomerServer
// for forward compatibility
type CustomerServer interface {
	// rpc CreateCustomer (CreateCustomerRequest) returns (CreateCustomerReply);
	// rpc UpdateCustomer (UpdateCustomerRequest) returns (UpdateCustomerReply);
	// rpc DeleteCustomer (DeleteCustomerRequest) returns (DeleteCustomerReply);
	GetCustomer(context.Context, *GetCustomerRequest) (*GetCustomerReply, error)
	// 登录
	Login(context.Context, *LoginReq) (*LoginRes, error)
	mustEmbedUnimplementedCustomerServer()
}

// UnimplementedCustomerServer must be embedded to have forward compatible implementations.
type UnimplementedCustomerServer struct {
}

func (UnimplementedCustomerServer) GetCustomer(context.Context, *GetCustomerRequest) (*GetCustomerReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCustomer not implemented")
}
func (UnimplementedCustomerServer) Login(context.Context, *LoginReq) (*LoginRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedCustomerServer) mustEmbedUnimplementedCustomerServer() {}

// UnsafeCustomerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CustomerServer will
// result in compilation errors.
type UnsafeCustomerServer interface {
	mustEmbedUnimplementedCustomerServer()
}

func RegisterCustomerServer(s grpc.ServiceRegistrar, srv CustomerServer) {
	s.RegisterService(&Customer_ServiceDesc, srv)
}

func _Customer_GetCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomerServer).GetCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Customer_GetCustomer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomerServer).GetCustomer(ctx, req.(*GetCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Customer_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomerServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Customer_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomerServer).Login(ctx, req.(*LoginReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Customer_ServiceDesc is the grpc.ServiceDesc for Customer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Customer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.customer.Customer",
	HandlerType: (*CustomerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCustomer",
			Handler:    _Customer_GetCustomer_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Customer_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/customer/customer.proto",
}

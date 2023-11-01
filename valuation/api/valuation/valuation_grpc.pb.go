// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: valuation/valuation.proto

package valuation

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
	Valuation_GetEstimatePrice_FullMethodName = "/api.valuation.Valuation/GetEstimatePrice"
)

// ValuationClient is the client API for Valuation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ValuationClient interface {
	GetEstimatePrice(ctx context.Context, in *GetEstimatePriceRequest, opts ...grpc.CallOption) (*GetEstimatePriceReply, error)
}

type valuationClient struct {
	cc grpc.ClientConnInterface
}

func NewValuationClient(cc grpc.ClientConnInterface) ValuationClient {
	return &valuationClient{cc}
}

func (c *valuationClient) GetEstimatePrice(ctx context.Context, in *GetEstimatePriceRequest, opts ...grpc.CallOption) (*GetEstimatePriceReply, error) {
	out := new(GetEstimatePriceReply)
	err := c.cc.Invoke(ctx, Valuation_GetEstimatePrice_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ValuationServer is the server API for Valuation service.
// All implementations must embed UnimplementedValuationServer
// for forward compatibility
type ValuationServer interface {
	GetEstimatePrice(context.Context, *GetEstimatePriceRequest) (*GetEstimatePriceReply, error)
	mustEmbedUnimplementedValuationServer()
}

// UnimplementedValuationServer must be embedded to have forward compatible implementations.
type UnimplementedValuationServer struct {
}

func (UnimplementedValuationServer) GetEstimatePrice(context.Context, *GetEstimatePriceRequest) (*GetEstimatePriceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEstimatePrice not implemented")
}
func (UnimplementedValuationServer) mustEmbedUnimplementedValuationServer() {}

// UnsafeValuationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ValuationServer will
// result in compilation errors.
type UnsafeValuationServer interface {
	mustEmbedUnimplementedValuationServer()
}

func RegisterValuationServer(s grpc.ServiceRegistrar, srv ValuationServer) {
	s.RegisterService(&Valuation_ServiceDesc, srv)
}

func _Valuation_GetEstimatePrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEstimatePriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ValuationServer).GetEstimatePrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Valuation_GetEstimatePrice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ValuationServer).GetEstimatePrice(ctx, req.(*GetEstimatePriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Valuation_ServiceDesc is the grpc.ServiceDesc for Valuation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Valuation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.valuation.Valuation",
	HandlerType: (*ValuationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEstimatePrice",
			Handler:    _Valuation_GetEstimatePrice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "valuation/valuation.proto",
}

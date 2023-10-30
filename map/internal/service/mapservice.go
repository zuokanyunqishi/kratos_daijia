package service

import (
	"context"
	"map/internal/biz"

	pb "map/api/mapService"
)

type MapServiceService struct {
	pb.UnimplementedMapServiceServer
	mapDriveBiz *biz.MapServiceBiz
}

func NewMapServiceService(biz *biz.MapServiceBiz) *MapServiceService {
	return &MapServiceService{mapDriveBiz: biz}
}

func (s *MapServiceService) GetDriverInfo(ctx context.Context, req *pb.GetDriverInfoReq) (*pb.GetDriverInfoResp, error) {

	s.mapDriveBiz.GetDriverInfo(ctx, req.Origin, req.Destination)
	return &pb.GetDriverInfoResp{}, nil
}

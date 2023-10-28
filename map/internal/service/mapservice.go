package service

import (
	"context"

	pb "map/api/mapService"
)

type MapServiceService struct {
	pb.UnimplementedMapServiceServer
}

func NewMapServiceService() *MapServiceService {
	return &MapServiceService{}
}

func (s *MapServiceService) GetDrivingInfo(ctx context.Context, req *pb.GetDrivingInfoReq) (*pb.GetDrivingInfoResp, error) {

	return &pb.GetDrivingInfoResp{}, nil
}

package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"map/internal/biz"
	"net/http"

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

	distance, duration, err := s.mapDriveBiz.GetDriverInfo(ctx, req.Origin, req.Destination)

	if err != nil {
		return nil, errors.New(http.StatusBadRequest, err.Error(), "获取路程信息失败")
	}

	return &pb.GetDriverInfoResp{
		Origin:      req.Origin,
		Destination: req.Destination,
		Distance:    distance,
		Duration:    duration,
		Code:        1,
		Message:     "",
	}, nil
}

package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"net/http"
	"strconv"
	"time"
	"valuation/internal/biz"

	pb "valuation/api/valuation"
)

type ValuationService struct {
	pb.UnimplementedValuationServer
	*biz.ValuationBiz
}

func NewValuationService(valuationBiz *biz.ValuationBiz) *ValuationService {
	return &ValuationService{ValuationBiz: valuationBiz}
}

func (s *ValuationService) GetEstimatePrice(ctx context.Context, req *pb.GetEstimatePriceRequest) (*pb.GetEstimatePriceReply, error) {
	drivierInfo, err := s.ValuationBiz.GetDrivingInfo(ctx, req.Origin, req.Destination)

	distance, _ := strconv.ParseInt(drivierInfo.Distance, 10, 64)
	duration, _ := strconv.ParseInt(drivierInfo.Duration, 10, 64)
	price, err := s.GetPrice(ctx, 1, time.Now().Hour(), distance/1000, duration/60)
	if err != nil {
		return nil, errors.New(http.StatusOK, err.Error(), "价格获取失败")
	}
	return &pb.GetEstimatePriceReply{
		Origin:      req.Origin,
		Destination: req.Destination,
		Price:       price,
	}, nil
}

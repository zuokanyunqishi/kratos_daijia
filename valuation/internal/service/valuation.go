package service

import (
	"context"
	"valuation/internal/biz"

	pb "valuation/api/valuation"
)

type ValuationService struct {
	pb.UnimplementedValuationServer
	biz.ValuationBiz
}

func NewValuationService(valuationBiz biz.ValuationBiz) *ValuationService {
	return &ValuationService{ValuationBiz: valuationBiz}
}

func (s *ValuationService) GetEstimatePrice(ctx context.Context, req *pb.GetEstimatePriceRequest) (*pb.GetEstimatePriceReply, error) {
	return &pb.GetEstimatePriceReply{}, nil
}

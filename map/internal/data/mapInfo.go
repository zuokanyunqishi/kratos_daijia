package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type MapInfo struct {
	data *Data
	log  *log.Helper
}

func (m MapInfo) GetDrivingInfo(ctx context.Context, origin, destination string) {
	//TODO implement me
	panic("implement me")
}

func NewMapInfo(data *Data, log *log.Helper) *MapInfo {
	return &MapInfo{data: data, log: log}
}

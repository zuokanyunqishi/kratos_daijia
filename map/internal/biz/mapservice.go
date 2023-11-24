package biz

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-resty/resty/v2"
	"map/internal/conf"
)

type MapDriverInfoRepo interface {
	GetDriverInfo(ctx context.Context, origin, destination string)
}

type MapServiceBiz struct {
	//mdi *MapDriverInfoRepo
	log      *log.Helper
	confAmap *conf.Amap
}

func NewMapServiceBiz(confAmap *conf.Amap, logger log.Logger) *MapServiceBiz {
	return &MapServiceBiz{log: log.NewHelper(logger), confAmap: confAmap}
}

func (b *MapServiceBiz) GetDriverInfo(ctx context.Context, origin, destination string) (string, string, error) {

	url := "https://restapi.amap.com/v3/direction/driving"
	// 高德路径规划服务key
	key := b.confAmap.GetDirection().GetKey()
	httpClient := resty.New()
	httpClient.JSONUnmarshal = sonic.Unmarshal
	httpClient.JSONMarshal = sonic.Marshal
	httpClient.EnableTrace()

	response, err := httpClient.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"origin":      origin,
			"destination": destination,
			"extensions":  "all",
			"output":      "json",
			"key":         key,
		}).Get(url)

	if err != nil {
		return "", "", err
	}

	var directionResp DirectionDrivingResp
	_ = httpClient.JSONUnmarshal(response.Body(), &directionResp)
	if directionResp.Status != "1" {
		return "", "", errors.New(directionResp.Info)
	}
	return directionResp.Route.Paths[0].Distance, directionResp.Route.Paths[0].Duration, nil
}

// DirectionDrivingResp 返回数据
type DirectionDrivingResp struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Count    string `json:"count"`
	Route    Route  `json:"route"`
}

type Route struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	TaxiCost    string `json:"taxi_cost"`
	Paths       []Path `json:"paths"`
}

type Path struct {
	Distance      string `json:"distance"`
	Duration      string `json:"duration"`
	Strategy      string `json:"strategy"`
	Tolls         string `json:"tolls"`
	TollDistance  string `json:"toll_distance"`
	Steps         []Step `json:"steps"`
	Restriction   string `json:"restriction"`
	TrafficLights string `json:"traffic_lights"`
}

type Step struct {
	Instruction     string        `json:"instruction"`
	Orientation     string        `json:"orientation"`
	Distance        string        `json:"distance"`
	Tolls           string        `json:"tolls"`
	TollDistance    string        `json:"toll_distance"`
	TollRoad        []interface{} `json:"toll_road"`
	Duration        string        `json:"duration"`
	Polyline        string        `json:"polyline"`
	Action          interface{}   `json:"action"`
	AssistantAction interface{}   `json:"assistant_action"`
	Tmcs            []Tmc         `json:"tmcs"`
	Cities          []City        `json:"cities"`
	Road            string        `json:"road,omitempty"`
}

type Tmc struct {
	Lcode    []interface{} `json:"lcode"`
	Distance string        `json:"distance"`
	Status   string        `json:"status"`
	Polyline string        `json:"polyline"`
}

type City struct {
	Name      string     `json:"name"`
	Citycode  string     `json:"citycode"`
	Adcode    string     `json:"adcode"`
	Districts []District `json:"districts"`
}

type District struct {
	Name   string `json:"name"`
	Adcode string `json:"adcode"`
}

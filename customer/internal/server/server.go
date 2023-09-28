package server

import (
	//consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/google/wire"
	//consulAPI "github.com/hashicorp/consul/api"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)

//func NewRegistrar(conf *conf.Registry) registry.Registrar {
//	c := consulAPI.DefaultConfig()
//	c.Address = conf.Consul.Address
//	c.Scheme = conf.Consul.Scheme
//	cli, err := consulAPI.NewClient(c)
//	if err != nil {
//		panic(err)
//	}
//	r := consul.New(cli, consul.WithHealthCheck(false))
//	return r
//}

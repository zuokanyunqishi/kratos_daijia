package main

import (
	"flag"
	consulConfig "github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	consulApi "github.com/hashicorp/consul/api"
	"map/internal/boot"
	"os"

	"map/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "mapService"
	// Version is the version of the compiled software.
	Version string = "v1"
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, rr registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(id+"_"+Name),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(rr),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	consulClient, err := consulApi.NewClient(&consulApi.Config{
		Address: bc.Registry.Consul.Address,
	})
	if err != nil {
		panic(err)
	}
	// path  map/amap.yaml
	// amap:
	// direction:
	//        key: xxx
	cs, err := consulConfig.New(consulClient, consulConfig.WithPath("map/"))
	cc := config.New(config.WithSource(cs))

	if err := cc.Load(); err != nil {
		panic(err)
	}
	var cbc conf.Bootstrap
	if err := cc.Scan(&cbc); err != nil {
		panic(err)
	}
	bc.Amap = cbc.Amap

	// trace
	boot.NewTrace(&bc).Run(Name, Version, id)

	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Registry, bc.Amap, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

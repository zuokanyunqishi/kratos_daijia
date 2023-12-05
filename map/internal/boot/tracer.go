package boot

import (
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"map/internal/conf"
)

type Trace struct {
	*conf.Bootstrap
}

func NewTrace(conf *conf.Bootstrap) *Trace {
	return &Trace{conf}
}

func (t *Trace) exporter(types string) (tracesdk.SpanExporter, error) {
	var exporter tracesdk.SpanExporter
	var err error

	switch types {
	case "stdout":
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint())
	case "jaeger":
		exporter, err = jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(t.Trace.Endpoint),
			))
	}
	return exporter, err
}

func (t *Trace) Run(sName, SVersion, InsId string) {

	exporter, err := t.exporter(t.Trace.Exporter)
	if err != nil {
		log.Errorw("trace_exporter_error", err)
		panic(err)
	}
	// 创建 trace provider
	tp := tracesdk.NewTracerProvider(
		// 采样设置
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// 使用 jaeger exporter
		tracesdk.WithBatcher(exporter),
		// 设置服务信息,作为 资源告知
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(sName), // 服务名

			semconv.ServiceVersionKey.String(SVersion), // 服务版本
			semconv.ServiceInstanceIDKey.String(InsId), // 服务实例ID
		),
		),
	)

	otel.SetTracerProvider(tp)
}

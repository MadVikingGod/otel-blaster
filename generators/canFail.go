package generators

import (
	"time"

	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/export/trace"
)

type CanFail struct {
	FailProbibality float64
	FailTime        time.Duration
	Generator       SpanGenerator
}

func (cf *CanFail) Generate(start time.Time, parent apitrace.SpanContext) ([]*sdktrace.SpanData, time.Time) {
	if randGen.Float64() > cf.FailProbibality {
		return cf.Generator.Generate(start, parent)
	}

	ctx := apitrace.SpanContext{
		TraceID:    parent.TraceID,
		SpanID:     spanID(),
		TraceFlags: 1,
	}

	return []*sdktrace.SpanData{
		{
			SpanContext:    ctx,
			ParentSpanID:   parent.SpanID,
			Name:           cf.Generator.Name,
			Attributes:     cf.Generator.Attributes,
			ChildSpanCount: 0,
			StartTime:      start,
			EndTime:        start.Add(cf.FailTime),
			StatusCode:     codes.Error,
		},
	}, start.Add(cf.FailTime)
}

var _ Generator = &CanFail{}

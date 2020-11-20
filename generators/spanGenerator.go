package generators

import (
	"time"

	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/export/trace"
)

type SpanGenerator struct {
	Name       string
	StarPad    time.Duration
	EndPad     time.Duration
	Events     Events
	Attributes []label.KeyValue
	Children   []Generator
}

type Events []Event

type Event struct {
	Name       string
	Attributes []label.KeyValue
	Duration   time.Duration
}

func (e Events) Make(start, end time.Time) []sdktrace.Event {
	events := make([]sdktrace.Event, len(e))
	for i, event := range e {
		at := start.Add(event.Duration)
		if at.After(end) {
			at = end
		}
		events[i] = sdktrace.Event{
			Name:       event.Name,
			Attributes: event.Attributes,
			Time:       at,
		}
	}
	return events
}

func (s *SpanGenerator) Generate(start time.Time, parent apitrace.SpanContext) ([]*sdktrace.SpanData, time.Time) {
	ctx := apitrace.SpanContext{
		TraceID:    parent.TraceID,
		SpanID:     spanID(),
		TraceFlags: 1,
	}
	childStart := start.Add(s.StarPad)
	collection := make([]*sdktrace.SpanData, 0, len(s.Children)*2)
	endTime := childStart
	for _, child := range s.Children {
		spans, end := child.Generate(childStart, ctx)
		collection = append(collection, spans...)
		if end.After(endTime) {
			endTime = end
		}
	}
	endTime = endTime.Add(s.EndPad)

	span := &sdktrace.SpanData{
		SpanContext:    ctx,
		ParentSpanID:   parent.SpanID,
		Name:           s.Name,
		StartTime:      start,
		EndTime:        endTime,
		ChildSpanCount: len(collection),
		Attributes:     s.Attributes,
		MessageEvents:  s.Events.Make(start, endTime),
	}
	return append(collection, span), endTime
}

var _ Generator = &SpanGenerator{}

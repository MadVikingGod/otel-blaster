package generators

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"time"

	apitrace "go.opentelemetry.io/otel/api/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/export/trace"
)

type Generator interface {
	Generate(start time.Time, parent apitrace.SpanContext) ([]*sdktrace.SpanData, time.Time)
}

func Generate(gen Generator, start time.Time) []*sdktrace.SpanData {
	spans, _ := gen.Generate(start, apitrace.SpanContext{
		TraceID:    traceID(),
		TraceFlags: 1,
	})
	return spans
}

var randGen *rand.Rand

func init() {
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	randGen = rand.New(rand.NewSource(rngSeed))
}

func spanID() apitrace.SpanID {
	sid := apitrace.SpanID{}
	randGen.Read(sid[:])
	return sid
}
func traceID() apitrace.ID {
	tid := apitrace.ID{}
	randGen.Read(tid[:])
	return tid
}

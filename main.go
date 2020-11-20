package main

import (
	"context"
	"log"
	"time"

	"github.com/madvikinggod/otel-blaster/generators"

	"go.opentelemetry.io/otel/exporters/stdout"
)

var exampleGenerators = &generators.SpanGenerator{
	Name:    "base",
	StarPad: time.Second,
	EndPad:  5 * time.Second,
	Children: []generators.Generator{
		&generators.SpanGenerator{
			Name:    "child1",
			StarPad: 2 * time.Second,
			Children: []generators.Generator{
				&generators.SpanGenerator{
					Name:   "child3",
					EndPad: 4 * time.Second,
				},
			},
		},
		&generators.SpanGenerator{
			Name:    "child2",
			StarPad: 3 * time.Second,
		},
	},
}

func main() {

	exporter, err := stdout.NewExporter([]stdout.Option{
		stdout.WithPrettyPrint(),
	}...)
	if err != nil {
		log.Fatalf("failed to initialize stdout export pipeline: %v", err)
	}

	spans := generators.Generate(exampleGenerators, time.Now())

	exporter.ExportSpans(context.Background(), spans)
}

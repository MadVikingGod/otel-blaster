package generators

import (
	"testing"
	"time"

	"go.opentelemetry.io/otel/label"
)

func TestGenerate(t *testing.T) {
	type fields struct {
		Name       string
		StarPad    time.Duration
		EndPad     time.Duration
		Events     Events
		Attributes []label.KeyValue
		Children   []Generator
	}
	tests := []struct {
		name         string
		fields       fields
		want         time.Duration
		wantChildren int
	}{
		{
			name: "base case",
			fields: fields{
				Name:    "base",
				StarPad: time.Second,
				EndPad:  5 * time.Second,
			},
			want:         6 * time.Second,
			wantChildren: 0,
		},
		{
			name: "parallel children",
			fields: fields{
				Name:    "base",
				StarPad: time.Second,
				EndPad:  5 * time.Second,
				Children: []Generator{
					&SpanGenerator{
						Name:    "child1",
						StarPad: 2 * time.Second,
					},
					&SpanGenerator{
						Name:    "child2",
						StarPad: 3 * time.Second,
					},
				},
			},
			want:         9 * time.Second,
			wantChildren: 2,
		},
		{
			name: "sub children",
			fields: fields{
				Name:    "base",
				StarPad: time.Second,
				EndPad:  5 * time.Second,
				Children: []Generator{
					&SpanGenerator{
						Name:    "child1",
						StarPad: 2 * time.Second,
						Children: []Generator{
							&SpanGenerator{
								Name:   "child3",
								EndPad: 4 * time.Second,
							},
						},
					},
					&SpanGenerator{
						Name:    "child2",
						StarPad: 3 * time.Second,
					},
				},
			},
			want:         12 * time.Second,
			wantChildren: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SpanGenerator{
				Name:       tt.fields.Name,
				StarPad:    tt.fields.StarPad,
				EndPad:     tt.fields.EndPad,
				Events:     tt.fields.Events,
				Attributes: tt.fields.Attributes,
				Children:   tt.fields.Children,
			}
			now := time.Now()
			got := Generate(s, now)

			// The last one returned should be the root span
			rootSpan := got[len(got)-1]
			if rootSpan.Name != s.Name {
				t.Errorf("SpanGenerator.Generate() name got = %v, want %v", got[len(got)-1].Name, s.Name)
			}
			if rootSpan.EndTime != now.Add(tt.want) {
				t.Errorf("SpanGenerator.Generate() endTime got = %v, want %v", rootSpan.EndTime, now.Add(tt.want))
			}
			if rootSpan.ChildSpanCount != tt.wantChildren {
				t.Errorf("SpanGenerator.Generate() Child Count got = %v, want %v", rootSpan.ChildSpanCount, tt.wantChildren)
			}
		})
	}
}

func BenchmarkGenerate(b *testing.B) {

	gen := &SpanGenerator{
		Name:    "base",
		StarPad: time.Second,
		EndPad:  5 * time.Second,
		Children: []Generator{
			&SpanGenerator{
				Name:    "child1",
				StarPad: 2 * time.Second,
				Children: []Generator{
					&SpanGenerator{
						Name:   "child3",
						EndPad: 4 * time.Second,
					},
				},
			},
			&SpanGenerator{
				Name:    "child2",
				StarPad: 3 * time.Second,
			},
		},
	}
	now := time.Now()

	for i := 0; i < b.N; i++ {
		Generate(gen, now)
	}
}

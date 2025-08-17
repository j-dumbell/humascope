package humascope

import (
	"context"
	"runtime/pprof"

	"github.com/danielgtaylor/huma/v2"
	"github.com/grafana/pyroscope-go"
)

// NewPyroscopeMW creates a new middleware which instruments the handler with profiling labels.
// The labels should be provided as-per pprof.Labels.
func NewPyroscopeMW(labels ...string) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		pyroscope.TagWrapper(ctx.Context(), pprof.Labels(labels...), func(_ context.Context) {
			next(ctx)
		})
	}
}

// Register wraps huma.Register, instrumenting the handler with 'method' and 'path'
// profiling labels.
func Register[I, O any](api huma.API, op huma.Operation, handler func(context.Context, *I) (*O, error)) {
	pyroscopeMW := NewPyroscopeMW("method", op.Method, "path", op.Path)
	op.Middlewares = append([]func(ctx huma.Context, next func(huma.Context)){pyroscopeMW}, op.Middlewares...)
	huma.Register(api, op, handler)
}

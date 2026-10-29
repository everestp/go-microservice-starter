package telemetry

import "context"

// InitTracerProvider is a small compatibility hook for the bootstrap scaffold.
func InitTracerProvider(_ context.Context, _ string) (struct{}, error) {
	return struct{}{}, nil
}

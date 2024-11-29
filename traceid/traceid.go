package traceid

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/things-go/gin-contrib/utilities/sequence"
)

// Key to use when setting the trace id.
type ctxTraceIdKey struct{}

// Config defines the config for TraceId middleware
type Config struct {
	traceIdHeader string
	nextTraceID   func() string
}

// Option TraceId option
type Option func(*Config)

// WithTraceIdHeader optional request id header (default "X-Trace-Id")
func WithTraceIdHeader(s string) Option {
	return func(c *Config) {
		c.traceIdHeader = s
	}
}

// WithNextTraceId optional next trace id function (default NewSequence function use utilities/sequence)
func WithNextTraceId(f func() string) Option {
	return func(c *Config) {
		c.nextTraceID = f
	}
}

// TraceId is a middleware that injects a trace id into the context of each
// request. if it is empty, set to write head
//   - traceIdHeader is the name of the HTTP Header which contains the trace id.
//     Exported so that it can be changed by developers. (default "X-Trace-Id")
//   - nextTraceID generates the next trace id.(default NewSequence function use utilities/sequence)
func TraceId(opts ...Option) gin.HandlerFunc {
	cc := &Config{
		traceIdHeader: "X-Trace-Id",
		nextTraceID:   sequence.New().NewSequence,
	}
	for _, opt := range opts {
		opt(cc)
	}

	return func(c *gin.Context) {
		traceId := c.Request.Header.Get(cc.traceIdHeader)
		if traceId == "" {
			traceId = cc.nextTraceID()
		}
		// set response header
		c.Header(cc.traceIdHeader, traceId)
		// set request context
		c.Request = c.Request.WithContext(WithTraceId(c.Request.Context(), traceId))
		c.Next()
	}
}

func WithTraceId(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, ctxTraceIdKey{}, traceId)
}

// FromTraceId returns a trace id from the given context if one is present.
// Returns the empty string if a trace id cannot be found.
func FromTraceId(ctx context.Context) string {
	traceId, _ := ctx.Value(ctxTraceIdKey{}).(string)
	return traceId
}

func InjectNewFromTraceId(ctx, newCtx context.Context) context.Context {
	return WithTraceId(newCtx, FromTraceId(ctx))
}

// GetTraceId get trace id from gin.Context.
func GetTraceId(c *gin.Context) string {
	return FromTraceId(c.Request.Context())
}

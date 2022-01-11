package tracer

import (
	"context"
	"net/http"

	"github.com/reyhanfahlevi/pkg/go/tracer/nr"
)

// Config config
type Config struct {
	Appname          string
	IsEnableJaeger   bool
	IsEnableNewRelic bool
	NewRelic         NewRelicConfig
}

// NewRelicConfig is a newrelic's config
type NewRelicConfig struct {
	SecretKey string
	LogLevel  string
}

// Transaction bg transaction - at http handlers
type Transaction struct {
	ctx context.Context

	// md contains the metadata for a transaction
	md http.Header
}

// Finish finishing a span / transaction
func (s *Transaction) Finish(err ...*error) {
	if len(err) > 0 && *err[0] != nil {
		_ = nr.Error(s.ctx, *err[0]) // notice error
	}
	nr.EndTransaction(s.ctx)
}

type ISegment interface {
	End()
}

// Span bg at funcs
type Span struct {
	ctx       context.Context
	nrSegment ISegment

	// md contains the metadata for a transaction
	md http.Header
}

func (s *Span) Finish() {
	if s.nrSegment != nil {
		s.nrSegment.End()
	}
}

type DBConInfo struct {
	Name string
	Host string
	Port string
}

type sqlExtraArgs struct {
	DBConInfo
	Query string
}

// Options is options for span
type Options struct {
	SpanType   int
	ExtraArgs  interface{}
	ExtraParam map[string]interface{}
}

const (
	// SpanTypeSQL is span as SQL span
	SpanTypeSQL int = 101

	// SpanTypeRedis is span as Redis call span
	SpanTypeRedis int = 102

	// SpanTypeHTTPCall is span as HTTP call span
	SpanTypeHTTPCall int = 103

	// SpanTypeNSQPublish is a span as Message Broker NSQ Publish
	SpanTypeNSQPublish int = 104

	// SpanTypeElasticsearch is span as Elasticsarch call span
	SpanTypeElasticsearch int = 105
)

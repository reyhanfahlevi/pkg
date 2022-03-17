package nr

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/ennobelprakoso/pkg/go/log"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrlogrus"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

var (
	app *newrelic.Application

	once sync.Once
)

type Config struct {
	AppName   string
	SecretKey string
	LogLevel  string
}

// Init to initialize new relic in go app
func Init(cfg Config) (err error) {

	once.Do(func() {
		app, err = makeApplication(cfg)
		if err != nil {
			log.Errorf("Could not initialize newrelic tracer: %s", err.Error())
			return
		}
	})

	return err
}

func makeApplication(cfg Config) (*newrelic.Application, error) {
	if cfg.AppName == "" {
		return nil, errors.New("new relic app name is nil")
	}

	if cfg.SecretKey == "" {
		return nil, errors.New("new relic secret key is nil")
	}

	logLevel := logrus.WarnLevel
	if cfg.LogLevel != "" {
		logLevel = getLevel(cfg.LogLevel)
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigEnabled(true),
		newrelic.ConfigAppName(cfg.AppName),
		newrelic.ConfigLicense(cfg.SecretKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(config *newrelic.Config) {
			log := logrus.New()
			log.SetLevel(logLevel)
			config.Logger = nrlogrus.Transform(log)
		},
	)

	if err != nil {
		return nil, err
	}

	return app, nil
}

// StartTransactionWithName to create a new Transaction with Name
func StartTransactionWithName(ctx context.Context, name string) context.Context {
	if app == nil {
		return ctx
	}

	txn := app.StartTransaction(name)
	ctx = newrelic.NewContext(ctx, txn)

	return ctx
}

// StartGinTransactionWithName start gin request trx
func StartGinTransactionWithName(c *gin.Context, name string) *gin.Context {
	if app == nil {
		return c
	}

	txn := app.StartTransaction(name)
	txn.SetName(name)
	txn.SetWebRequestHTTP(c.Request)
	txn.SetWebResponse(c.Writer)

	c.Set("newRelicTransaction", txn)

	return c
}

// AddAttribute to add attribute to span
func AddAttribute(ctx context.Context, key string, value interface{}) context.Context {
	if app == nil {
		return ctx
	}

	txn := newrelic.FromContext(ctx)
	txn.AddAttribute(key, value)
	return newrelic.NewContext(ctx, txn)
}

func GetMetadataFromContext(ctx context.Context) http.Header {
	hdr := http.Header{}
	if app == nil {
		return hdr
	}

	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return hdr
	}
	txn.InsertDistributedTraceHeaders(hdr)
	return hdr
}

// EndTransaction to end NR transaction from context
func EndTransaction(ctx context.Context) {
	if app == nil {
		return
	}

	txn := newrelic.FromContext(ctx)
	txn.End()
}

// Error to notice error to NR
func Error(ctx context.Context, err error) error {
	if err == nil {
		return err
	}

	if app == nil {
		return err
	}

	txn := newrelic.FromContext(ctx)
	txn.NoticeError(err)

	// return parent error
	return err
}

// StartSegment start segment
func StartSegment(ctx context.Context, name string) *newrelic.Segment {
	if app == nil {
		return nil
	}

	txn := newrelic.FromContext(ctx)
	return txn.StartSegment(name)
}

// StartNSQSegment starts a new relic segment for NSQ produce topic
func StartNSQSegment(ctx context.Context, topicName string) *newrelic.MessageProducerSegment {
	seg := &newrelic.MessageProducerSegment{
		Library:         "NSQ",
		DestinationType: newrelic.MessageTopic,
		DestinationName: topicName,
	}
	txn := newrelic.FromContext(ctx)
	seg.StartTime = txn.StartSegmentNow()
	return seg
}

// StartPostgresSegment to start a postgres segment
func StartPostgresSegment(ctx context.Context, query, collection, operation string, params map[string]interface{}) *newrelic.DatastoreSegment {
	return StartPostgresSegmentWithDBName(ctx, "", "", "", query, collection, operation, params)
}

// StartPostgresSegmentWithDBName to start a postgres segment with DB name & the DB host
func StartPostgresSegmentWithDBName(
	ctx context.Context, dbName, dbHost, dbPort, query, collection, operation string, params map[string]interface{},
) *newrelic.DatastoreSegment {
	if app == nil {
		return &newrelic.DatastoreSegment{}
	}

	datastore := newrelic.DatastoreSegment{
		Product:            newrelic.DatastorePostgres,
		DatabaseName:       dbName,
		Host:               dbHost,
		PortPathOrID:       dbPort,
		Collection:         collection,
		Operation:          operation,
		ParameterizedQuery: query,
		QueryParameters:    params,
	}

	txn := newrelic.FromContext(ctx)
	datastore.StartTime = txn.StartSegmentNow()
	return &datastore
}

func getLevel(level string) logrus.Level {
	allLevel := map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"error": logrus.WarnLevel,
	}

	if value, ok := allLevel[level]; ok {
		return value
	}
	return logrus.ErrorLevel
}

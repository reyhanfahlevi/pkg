package tracer

import (
	"strings"

	"github.com/reyhanfahlevi/pkg/go/tracer/nr"
)

// Init to initialize tracing package
func Init(cfg Config) error {
	nrCfg := nr.Config{
		AppName:   cfg.Appname,
		SecretKey: cfg.NewRelic.SecretKey,
		LogLevel:  cfg.NewRelic.LogLevel,
	}
	err := nr.Init(nrCfg)
	if err != nil {
		return err
	}

	return nil
}

// getOperationFromSQLQuery to get DDL / DML operation
// example: query `SELECT $1 FROM table_name`, this func will return SELECT
// newrelic pkg will send both operation name and raw query to their data, so we can explore queries that slow and so on...
func getOperationFromSQLQuery(query string) string {
	// replace /t with space
	query = strings.Replace(query, "\t", " ", -1)

	// replace /n with space
	query = strings.Replace(query, "\n", " ", -1)

	// uppercase all
	query = strings.ToUpper(query)

	// trim space
	query = strings.TrimSpace(query)

	if strings.Contains(query, "INSERT ") {
		return "INSERT"
	}
	if strings.Contains(query, "UPDATE ") {
		return "UPDATE"
	}
	if strings.Contains(query, "DELETE ") {
		return "DELETE"
	}
	if strings.Contains(query, "SELECT ") {
		return "SELECT"
	}

	return strings.Split(query, " ")[0]
}

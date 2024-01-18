package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jsmit257/huautla/internal/metrics"
	"github.com/jsmit257/huautla/types"

	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

type (
	query interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...any) *sql.Row
	}

	Conn struct {
		query
		generateUUID uuidgen
		logger       *log.Entry
		sql          map[string]map[string]string
	}

	uuidgen func() uuid.UUID

	getMockDB func() *sql.DB

	deferred func(start time.Time, err error, l *log.Entry)
)

const connformat = ""

var mtrcs = metrics.DataMetrics.MustCurryWith(prometheus.Labels{"pkg": "data"})

func New(logger *log.Entry) (types.DB, error) {

	var err error

	result := &Conn{
		generateUUID: uuid.New,
		logger:       logger,
		sql:          readSQL("pgsql.yaml"),
	}

	result.query, err = sql.Open(
		"postgres",
		fmt.Sprintf(connformat))

	return result, err
}

func readSQL(filename string) map[string]map[string]string {
	var err error

	result := make(map[string]map[string]string)

	// open the file,
	// parse as yaml or panic
	if err != nil {
		panic(err)
	}

	return result
}

func initVendorFuncs(method string, l *log.Entry, err error, id types.UUID, cid types.CID) (deferred, time.Time, *log.Entry) {
	start := time.Now()
	l = l.WithFields(log.Fields{
		"method": method,
		"cid":    cid,
		"id":     id,
	})

	l.Info("starting work")

	return func(start time.Time, err error, l *log.Entry) {
			duration := time.Since(start)

			l.
				WithField("duration", duration).
				WithError(err).
				Infof("finished work")

			// TODO: metrics
		},
		start,
		l
}

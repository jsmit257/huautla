package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/jsmit257/huautla/types"

	// "github.com/prometheus/client_golang/prometheus"

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
		sql          map[string]string
	}

	uuidgen func() uuid.UUID

	getMockDB func() *sql.DB
)

const connformat = ""

// var mtrcs = metrics.DataMetrics.MustCurryWith(prometheus.Labels{"pkg": "mysql"})

func New(logger *log.Entry) (types.DB, error) {

	var err error

	result := &Conn{
		generateUUID: uuid.New,
		logger:       logger,
		sql:          readSQL(""),
	}

	result.query, err = sql.Open(
		"postgres",
		fmt.Sprintf(connformat))

	return result, err
}

func readSQL(filename string) map[string]string {
	var err error

	result := make(map[string]string)

	// open the file,
	// parse as yaml or panic
	if err != nil {
		panic(err)
	}

	return result
}

func mockUUIDGen() uuid.UUID {
	return uuid.Must(uuid.FromBytes([]byte("0123456789abcdef")))
}

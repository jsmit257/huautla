package data

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jsmit257/huautla/internal/metrics"
	"github.com/jsmit257/huautla/types"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
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

func New(cnxInfo string, log *log.Entry) (types.DB, error) {

	var err error
	var sqls map[string]map[string]string
	var query *sql.DB

	if sqls, err = readSQL("./pgsql.yaml"); err != nil {
		return nil, err
	} else if query, err = sql.Open("postgres", cnxInfo); err != nil {
		return nil, err
	}

	return &Conn{
		query:        query,
		generateUUID: uuid.New,
		logger:       log,
		sql:          sqls,
	}, nil
}

func (db *Conn) deleteByUUID(ctx context.Context, id types.UUID, cid types.CID, method, table string, l *log.Entry) error {
	var err error

	deferred, start, l := initAccessFuncs(method, l, id, cid)
	defer deferred(start, err, l)

	var result sql.Result

	result, err = db.ExecContext(ctx, db.sql[table]["delete"], id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		// this won't be reported in the WithError log in `defer ...`, b/c it's operator error
		return fmt.Errorf("%s could not be deleted: '%s'", table, id)
	}

	return err
}

func readSQL(filename string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	if yamlFile, err := os.ReadFile(filename); err != nil {
		wd, _ := os.Getwd()
		err = fmt.Errorf("pwd: '%s', err: %v", wd, err)
		return result, err
	} else if err = yaml.Unmarshal(yamlFile, &result); err != nil {
		return result, err
	}

	return result, nil
}

func initAccessFuncs(method string, l *log.Entry, id types.UUID, cid types.CID) (deferred, time.Time, *log.Entry) {
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

func isPrimaryKeyViolation(err error) bool {
	return false
}

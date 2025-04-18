package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"

	"github.com/google/uuid"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	log "github.com/sirupsen/logrus"

	pq "github.com/lib/pq"
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
		// sql          map[string]map[string]string
	}

	uuidgen func() uuid.UUID

	getMockDB func(*sql.DB, sqlmock.Sqlmock, error) *sql.DB

	deferred func(*error, *log.Entry)
)

func New(cnxInfo string, log *log.Entry) (types.DB, error) {
	var err error
	var query *sql.DB

	if query, err = sql.Open("postgres", cnxInfo); err != nil {
		return nil, err
	} else if err = query.Ping(); err != nil {
		return nil, err
	}

	return &Conn{
		query:        query,
		generateUUID: uuid.New,
		logger:       log,
	}, nil
}

func (db *Conn) deleteByUUID(ctx context.Context, id types.UUID, cid types.CID, method, table string, l *log.Entry) error {
	var err error

	deferred, l := initAccessFuncs(method, l, id, cid)
	defer deferred(&err, l)

	var result sql.Result

	result, err = db.ExecContext(ctx, psqls[table]["delete"], id)
	if err != nil {
		return pqerr(err)
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		// this won't be reported in the WithError log in `defer ...`, b/c it's operator error
		return fmt.Errorf("%s could not be deleted: '%s'", table, id)
	}

	return err
}

func initAccessFuncs(fn string, l *log.Entry, id any, cid types.CID) (deferred, *log.Entry) {
	start := time.Now()
	l = l.WithFields(log.Fields{
		"function": fn,
		"cid":      cid,
	})
	if id != nil {
		l = l.WithField("key", id)
	}

	l.Info("starting work")

	return func(err *error, l *log.Entry) {
		duration := time.Since(start)

		if err != nil {
			l.WithError(*err)
		}
		l.WithField("duration", duration).Infof("finished work")

		// TODO: metrics
	}, l
}

type PKeyError error
type FKeyError error
type UniqueError error
type ConstraintError error
type NullError error

func pqerr(err error) error {
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return err
	}
	switch pqErr.Code {
	case "23505": // unique_key
		// FIXME: can't tell the difference between primary key and other
		//  unique constraints; searching for `_pkey` seems too clumsy to
		//  be the right solution; is there another? until then, the result
		//  is always false
		return UniqueError(fmt.Errorf("unique key violation: %s",
			pqErr.Detail))
	case "23503":
		return FKeyError(fmt.Errorf("foreign key violation: %s, %s.%s",
			pqErr.Detail,
			pqErr.Table,
			pqErr.Column))
	case "23514":
		return ConstraintError(fmt.Errorf("constraint violation: %s, %s, %s.%s",
			pqErr.Detail,
			pqErr.Constraint,
			pqErr.Table,
			pqErr.Column))
	case "23502":
		return NullError(fmt.Errorf("field not nullable: %s, %s.%s",
			pqErr.Detail,
			pqErr.Table,
			pqErr.Column))
	}
	return pqErr
}
func isPrimaryKeyViolation(err error) bool {
	pqErr, ok := err.(*pq.Error)

	// see pqerr above for why this returns what it does
	return false && ok && pqErr.Code == "23505"
}

func isForeignKeyViolation(err error) bool {
	pqErr, ok := err.(*pq.Error)

	return ok && pqErr.Code == "23503" // FIXME: should be right
}

package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jsmit257/huautla/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	_ets = []eventtype{
		{UUID: "0", Name: "eventtype 0", Severity: "severity", Stage: types.Stage{UUID: "0", Name: "stage 0"}},
		{UUID: "1", Name: "eventtype 1", Severity: "severity", Stage: types.Stage{UUID: "1", Name: "stage 1"}},
		{UUID: "2", Name: "eventtype 2", Severity: "severity", Stage: types.Stage{UUID: "1", Name: "stage 1"}},
	}
	etFields = row{"id", "name", "severity", "stage_uuid", "stage_name"}
	etValues = [][]driver.Value{
		{_ets[0].UUID, _ets[0].Name, _ets[0].Severity, _ets[0].Stage.UUID, _ets[0].Stage.Name},
		{_ets[1].UUID, _ets[1].Name, _ets[1].Severity, _ets[1].Stage.UUID, _ets[1].Stage.Name},
		{_ets[2].UUID, _ets[2].Name, _ets[2].Severity, _ets[2].Stage.UUID, _ets[2].Stage.Name},
	}
)

func Test_SelectAllEventTypes(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "Test_SelectAllEventTypes")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.EventType
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				etFields.mock(mock, etValues...)
				return db
			},
			result: []types.EventType{
				types.EventType(_ets[0]),
				types.EventType(_ets[1]),
				types.EventType(_ets[2]),
			},
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, eventFields.fail())
				return db
			},
			err: eventFields.err(),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectAllEventTypes(context.Background(), "Test_SelectAllEventTypes")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectEventType(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectStrain")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.EventType
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				etFields.mock(mock, etValues[0])
				return db
			},
			id:     "0",
			result: types.EventType(_ets[0]),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, etFields.fail())
				return db
			},
			err: etFields.err(),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectEventType(context.Background(), tc.id, "Test_EventType")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertEventType(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "InsertEventType")

	tcs := map[string]struct {
		db     getMockDB
		result types.EventType
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			result: types.EventType{UUID: "30313233-3435-3637-3839-616263646566", Name: "eventtype 0", Stage: types.Stage{}},
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			result: types.EventType{UUID: "30313233-3435-3637-3839-616263646566", Name: "eventtype 0", Stage: types.Stage{}},
			err:    fmt.Errorf("eventtype was not added"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			result: types.EventType{UUID: "30313233-3435-3637-3839-616263646566", Name: "eventtype 0", Stage: types.Stage{}},
			err:    fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			result: types.EventType{UUID: "30313233-3435-3637-3839-616263646566", Name: "eventtype 0", Stage: types.Stage{}},
			err:    fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).InsertEventType(
				context.Background(),
				types.EventType{UUID: "0", Name: "eventtype 0", Stage: types.Stage{}},
				"Test_InsertEventType")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateEventType(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateStrain")

	tcs := map[string]struct {
		db  getMockDB
		id  types.UUID
		err error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("eventtype was not updated: '0'"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).UpdateEventType(
				context.Background(),
				tc.id,
				types.EventType{Name: "strain " + string(tc.id)},
				"Test_UpdateStrains")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_DeleteEventType(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "DeleteStrain")

	tcs := map[string]struct {
		db  getMockDB
		id  types.UUID
		err error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("eventtype could not be deleted: '0'"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).DeleteEventType(
				context.Background(),
				tc.id,
				"Test_DeleteStrain")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_EventTypeReport(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectStrain")

	tcs := map[string]struct {
		db                       getMockDB
		result, actual, expected types.Entity
		err                      error
	}{
		"happy_lifecycle_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					etFields.set(etValues[0]),
					lcFields.set(lcValues),
					eventFields.set(),
					attrFields.set(),
					ingFields.set(),
					ingFields.set(),
					napFields.set(),
					noteFields.set(),
					photoFields.set(),
					genFields.set())

				return db
			},
			result: func(e types.Entity) types.Entity {
				e["lifecycles"] = []types.Entity{mustEntity(_lc)}

				return e
			}(mustEntity(_ets[0])),
		},
		"lifecycle_path_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, etFields.set(etValues[0]), lcFields.fail())

				return db
			},
			err: lcFields.err(),
		},
		"happy_generation_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					etFields.set(etValues[0]),
					lcFields.set(),
					genFields.set(genValues),
					eventFields.set(),
					srcFields.set(),
					ingFields.set(),
					ingFields.set(),
					napFields.set(),
					noteFields.set(),
					strainFields.set())

				return db
			},
			result: func(e types.Entity) types.Entity {
				e["generations"] = []types.Entity{mustEntity(_gen)}

				return e
			}(mustEntity(_ets[0])),
		},
		"generation_path_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					etFields.set(etValues[0]),
					lcFields.set(),
					genFields.fail())

				return db
			},
			err: genFields.err(),
		},
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					etFields.set(etValues[0]),
					lcFields.set(),
					genFields.set())

				return db
			},
			result: func(e types.Entity) types.Entity {
				return e
			}(mustEntity(_ets[0])),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).EventTypeReport(context.Background(), "tc.id", "Test_EventType")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

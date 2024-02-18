package data

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jsmit257/huautla/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_SelectLifecycle(t *testing.T) {
	t.Parallel()

	fieldnames := []string{
		"name",
		"location",
		"graincost",
		"bulkcost",
		"yield",
		"count",
		"gross",
		"mtime",
		"ctime",
		"strain_uuid",
		"strain_name",
		"strain_vendor_uuid",
		"strain_vendor_name",
		"grain_substrate_uuid",
		"grain_substrate_name",
		"grain_substrate_type",
		"grain_vendor_uuid",
		"grain_vendor_name",
		"bulk_substrate_uuid",
		"bulk_substrate_name",
		"bulk_substrate_type",
		"bulk_vendor_uuid",
		"bulk_vendor_name",
	}

	e0, e1, e2 := types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

	whenwillthenbenow := time.Now() // time.Soon()

	l := log.WithField("test", "SelectLifecycle")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Lifecycle
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(fieldnames).
						AddRow(
							"name",
							"location",
							0,
							0,
							0,
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"0",
							"strain 0",
							"x",
							"vendor x",
							"gs",
							"gs",
							types.GrainType,
							"1",
							"vendor 1",
							"bs",
							"bs",
							types.BulkType,
							"2",
							"vendor 2"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "value"}).
						AddRow("0", "name 0", "value 0").
						AddRow("1", "name 1", "value 1").
						AddRow("2", "name 2", "value 2"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name"}).
						AddRow("0", "ingredient 0").
						AddRow("1", "ingredient 1").
						AddRow("2", "ingredient 2"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name"}).
						AddRow("0", "ingredient 0").
						AddRow("1", "ingredient 1").
						AddRow("2", "ingredient 2"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "eventtype_name", "eventtype_severity", "stage_uuid", "stage_name"}).
						AddRow(e0.UUID, e0.Temperature, e0.Humidity, e0.MTime, e0.CTime, e0.EventType.UUID, e0.EventType.Name, e0.EventType.Severity, e0.EventType.Stage.UUID, e0.EventType.Stage.Name).
						AddRow(e1.UUID, e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name).
						AddRow(e2.UUID, e2.Temperature, e2.Humidity, e2.MTime, e2.CTime, e2.EventType.UUID, e2.EventType.Name, e2.EventType.Severity, e2.EventType.Stage.UUID, e2.EventType.Stage.Name))

				return db
			},
			id: "0",
			result: types.Lifecycle{
				UUID:      "0",
				Name:      "name",
				Location:  "location",
				GrainCost: 0,
				BulkCost:  0,
				Yield:     0,
				Count:     0,
				Gross:     0,
				MTime:     whenwillthenbenow,
				CTime:     whenwillthenbenow,
				Strain: types.Strain{
					UUID: "0",
					Name: "strain 0",
					Vendor: types.Vendor{
						UUID: "x",
						Name: "vendor x",
					},
					Attributes: []types.StrainAttribute{
						{UUID: "0", Name: "name 0", Value: "value 0"},
						{UUID: "1", Name: "name 1", Value: "value 1"},
						{UUID: "2", Name: "name 2", Value: "value 2"},
					},
				},
				GrainSubstrate: types.Substrate{
					UUID: "gs",
					Name: "gs",
					Type: types.GrainType,
					Vendor: types.Vendor{
						UUID: "1",
						Name: "vendor 1",
					},
					Ingredients: []types.Ingredient{
						{UUID: "0", Name: "ingredient 0"},
						{UUID: "1", Name: "ingredient 1"},
						{UUID: "2", Name: "ingredient 2"},
					},
				},
				BulkSubstrate: types.Substrate{
					UUID: "bs",
					Name: "bs",
					Type: types.BulkType,
					Vendor: types.Vendor{
						UUID: "2",
						Name: "vendor 2",
					},
					Ingredients: []types.Ingredient{
						{UUID: "0", Name: "ingredient 0"},
						{UUID: "1", Name: "ingredient 1"},
						{UUID: "2", Name: "ingredient 2"},
					},
				},
				Events: []types.Event{e0, e1, e2},
			},
		},
		"all_attrs_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(fieldnames).
						AddRow(
							"name",
							"location",
							0,
							0,
							0,
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"0",
							"strain 0",
							"x",
							"vendor x",
							"gs",
							"gs",
							types.GrainType,
							"1",
							"vendor 1",
							"bs",
							"bs",
							types.BulkType,
							"2",
							"vendor 2"))
				mock.ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))

				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"get_grain_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(fieldnames).
						AddRow(
							"name",
							"location",
							0,
							0,
							0,
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"0",
							"strain 0",
							"x",
							"vendor x",
							"gs",
							"gs",
							types.GrainType,
							"1",
							"vendor 1",
							"bs",
							"bs",
							types.BulkType,
							"2",
							"vendor 2"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "value"}).
						AddRow("0", "name 0", "value 0").
						AddRow("1", "name 1", "value 1").
						AddRow("2", "name 2", "value 2"))
				mock.ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))

				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"get_bulk_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(fieldnames).
						AddRow(
							"name",
							"location",
							0,
							0,
							0,
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"0",
							"strain 0",
							"x",
							"vendor x",
							"gs",
							"gs",
							types.GrainType,
							"1",
							"vendor 1",
							"bs",
							"bs",
							types.BulkType,
							"2",
							"vendor 2"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "value"}).
						AddRow("0", "name 0", "value 0").
						AddRow("1", "name 1", "value 1").
						AddRow("2", "name 2", "value 2"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "value"}).
						AddRow("0", "name 0", "value 0").
						AddRow("1", "name 1", "value 1").
						AddRow("2", "name 2", "value 2"))
				mock.ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))

				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectLifecycle(context.Background(), tc.id, "Test_SelectLifecycle")

			require.Equal(t, tc.err, err)
			// require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertLifecycle(t *testing.T) {
	t.Parallel()

	fieldnames := []string{
		"name",
		"location",
		"graincost",
		"bulkcost",
		"yield",
		"count",
		"gross",
		"mtime",
		"ctime",
		"strain_uuid",
		"strain_name",
		"strain_vendor_uuid",
		"strain_vendor_name",
		"grain_substrate_uuid",
		"grain_substrate_name",
		"grain_substrate_type",
		"grain_vendor_uuid",
		"grain_vendor_name",
		"bulk_substrate_uuid",
		"bulk_substrate_name",
		"bulk_substrate_type",
		"bulk_vendor_uuid",
		"bulk_vendor_name",
	}

	e0, e1, e2 := types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

	whenwillthenbenow := time.Now() // time.Soon()

	l := log.WithField("test", "InsertLifecycle")

	tcs := map[string]struct {
		db  getMockDB
		err error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(fieldnames).
						AddRow(
							"name",
							"location",
							0,
							0,
							0,
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"0",
							"strain 0",
							"x",
							"vendor x",
							"gs",
							"gs",
							types.GrainType,
							"1",
							"vendor 1",
							"bs",
							"bs",
							types.BulkType,
							"2",
							"vendor 2"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "value"}).
						AddRow("0", "name 0", "value 0").
						AddRow("1", "name 1", "value 1").
						AddRow("2", "name 2", "value 2"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name"}).
						AddRow("0", "ingredient 0").
						AddRow("1", "ingredient 1").
						AddRow("2", "ingredient 2"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name"}).
						AddRow("0", "ingredient 0").
						AddRow("1", "ingredient 1").
						AddRow("2", "ingredient 2"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "eventtype_name", "eventtype_severity", "stage_uuid", "stage_name"}).
						AddRow(e0.UUID, e0.Temperature, e0.Humidity, e0.MTime, e0.CTime, e0.EventType.UUID, e0.EventType.Name, e0.EventType.Severity, e0.EventType.Stage.UUID, e0.EventType.Stage.Name).
						AddRow(e1.UUID, e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name).
						AddRow(e2.UUID, e2.Temperature, e2.Humidity, e2.MTime, e2.CTime, e2.EventType.UUID, e2.EventType.Name, e2.EventType.Severity, e2.EventType.Stage.UUID, e2.EventType.Stage.Name))

				return db
			},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("lifecycle was not added: 0"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lc, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).InsertLifecycle(
				context.Background(),
				types.Lifecycle{},
				"Test_InsertLifecycle")

			require.Equal(t, tc.err, err)
			require.Equal(t, lc.UUID, types.UUID(mockUUIDGen().String()))
		})
	}
}

func Test_UpdateLifecycle(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateLifecycle")

	tcs := map[string]struct {
		db  getMockDB
		lc  types.Lifecycle
		err error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			lc: types.Lifecycle{},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			lc:  types.Lifecycle{},
			err: fmt.Errorf("lifecycle was not updated"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			lc:  types.Lifecycle{},
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			lc:  types.Lifecycle{},
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// there's no good way to test the returned lifecycle, to start, the
			// timestamps are non-deterministic; system tests will vet the rest
			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).UpdateLifecycle(context.Background(), tc.lc, "Test_UpdateLifecycle")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_DeleteLifecycle(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "DeleteLifecycle")

	tcs := map[string]struct {
		db  getMockDB
		id  types.UUID
		err error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("lifecycle could not be deleted: '0'"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
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
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).DeleteLifecycle(
				context.Background(),
				tc.id,
				"Test_DeleteLifecycle")

			require.Equal(t, tc.err, err)
		})
	}
}

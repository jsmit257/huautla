package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jsmit257/huautla/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	e0, e1, e2 = types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

	lcfieldnames = []string{
		"uuid",
		"location",
		"straincost",
		"graincost",
		"bulkcost",
		"yield",
		"count",
		"gross",
		"mtime",
		"ctime",
		"strain_uuid",
		"strain_species",
		"strain_name",
		"generation_uuid",
		"strain_ctime",
		"strain_vendor_uuid",
		"strain_vendor_name",
		"strain_vendor_website",
		"grain_substrate_uuid",
		"grain_substrate_name",
		"grain_substrate_type",
		"grain_vendor_uuid",
		"grain_vendor_name",
		"grain_vendor_website",
		"bulk_substrate_uuid",
		"bulk_substrate_name",
		"bulk_substrate_type",
		"bulk_vendor_uuid",
		"bulk_vendor_name",
		"bulk_vendor_website",
	}

	lctestrow = []driver.Value{
		"30313233-3435-3637-3839-616263646566",
		"location",
		0,
		0,
		0,
		0,
		0,
		0,
		whenwillthenbenow,
		whenwillthenbenow,
		"0",
		"X.species",
		"strain 0",
		"nil",
		whenwillthenbenow,
		"x",
		"vendor x",
		"website",
		"gs",
		"gs",
		types.GrainType,
		"1",
		"vendor 1",
		"website",
		"bs",
		"bs",
		types.BulkType,
		"2",
		"vendor 2",
		"website",
	}

	lchappyresults = []types.Lifecycle{
		{
			UUID:       "0",
			Location:   "location",
			StrainCost: 0,
			GrainCost:  0,
			BulkCost:   0,
			Yield:      0,
			Count:      0,
			Gross:      0,
			MTime:      whenwillthenbenow,
			CTime:      whenwillthenbenow,
			Strain: types.Strain{
				UUID:    "0",
				Species: "X.species",
				Name:    "strain 0",
				CTime:   whenwillthenbenow,
				Vendor: types.Vendor{
					UUID:    "x",
					Name:    "vendor x",
					Website: "website",
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
					UUID:    "1",
					Name:    "vendor 1",
					Website: "website",
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
					UUID:    "2",
					Name:    "vendor 2",
					Website: "website",
				},
				Ingredients: []types.Ingredient{
					{UUID: "0", Name: "ingredient 0"},
					{UUID: "1", Name: "ingredient 1"},
					{UUID: "2", Name: "ingredient 2"},
				},
			},
			Events: []types.Event{e0, e1, e2},
		},
	}
)

func Test_SelectLifecycleIndex(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectLifecycleIndex")

	tcs := map[string]struct {
		db     getMockDB
		result []types.Lifecycle
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{
							"uuid",
							"location",
							"mtime",
							"ctime",
							"strain_uuid",
							"strain_species",
							"strain_name",
							"strain_ctime",
							"vendor_uuid",
							"vendor_name",
							"vendor_website",
							"event_uuid",
							"temp",
							"humidity",
							"event_mtime",
							"event_ctime",
							"et_uuid",
							"et_name",
							"et_sev",
							"stage_uuid",
							"stage_name"}).
						AddRow(
							"0",
							"happy_path",
							whenwillthenbenow,
							whenwillthenbenow,
							"strain 0",
							"strain 0",
							"strain 0",
							whenwillthenbenow,
							"vendor 0",
							"vendor 0",
							"vendor 0",
							"event 0",
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"type 0",
							"type 0",
							"type 0",
							"stage 0",
							"stage 0",
						).
						AddRow(
							"1",
							"happy_path 2",
							whenwillthenbenow,
							whenwillthenbenow,
							"strain 0",
							"strain 0",
							"strain 0",
							whenwillthenbenow,
							"vendor 0",
							"vendor 0",
							"vendor 0",
							"event 0",
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"type 0",
							"type 0",
							"type 0",
							"stage 0",
							"stage 0",
						).
						AddRow(
							"1",
							"happy_path 2",
							whenwillthenbenow,
							whenwillthenbenow,
							"strain 0",
							"strain 0",
							"strain 0",
							whenwillthenbenow,
							"vendor 0",
							"vendor 0",
							"vendor 0",
							"event 1",
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"type 0",
							"type 0",
							"type 0",
							"stage 0",
							"stage 0",
						))
				return db
			},
			result: []types.Lifecycle{
				{
					UUID:     "0",
					Location: "happy_path",
					MTime:    whenwillthenbenow,
					CTime:    whenwillthenbenow,
					Strain: types.Strain{
						UUID:    "strain 0",
						Name:    "strain 0",
						Species: "strain 0",
						CTime:   whenwillthenbenow,
						Vendor: types.Vendor{
							UUID:    "vendor 0",
							Name:    "vendor 0",
							Website: "vendor 0",
						},
					},
					Events: []types.Event{{
						UUID:  "event 0",
						MTime: whenwillthenbenow,
						CTime: whenwillthenbenow,
						EventType: types.EventType{
							UUID:     "type 0",
							Name:     "type 0",
							Severity: "type 0",
							Stage: types.Stage{
								UUID: "stage 0",
								Name: "stage 0",
							},
						},
					}},
				},
				{
					UUID:     "1",
					Location: "happy_path 2",
					MTime:    whenwillthenbenow,
					CTime:    whenwillthenbenow,
					Strain: types.Strain{
						UUID:    "strain 0",
						Name:    "strain 0",
						Species: "strain 0",
						CTime:   whenwillthenbenow,
						Vendor: types.Vendor{
							UUID:    "vendor 0",
							Name:    "vendor 0",
							Website: "vendor 0",
						},
					},
					Events: []types.Event{
						{
							UUID:  "event 0",
							MTime: whenwillthenbenow,
							CTime: whenwillthenbenow,
							EventType: types.EventType{
								UUID:     "type 0",
								Name:     "type 0",
								Severity: "type 0",
								Stage: types.Stage{
									UUID: "stage 0",
									Name: "stage 0",
								},
							},
						},
						{
							UUID:  "event 1",
							MTime: whenwillthenbenow,
							CTime: whenwillthenbenow,
							EventType: types.EventType{
								UUID:     "type 0",
								Name:     "type 0",
								Severity: "type 0",
								Stage: types.Stage{
									UUID: "stage 0",
									Name: "stage 0",
								},
							},
						},
					},
				},
			},
		},
		"db_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		// "row_error": {
		// 	db: func() *sql.DB {
		// 		db, mock, _ := sqlmock.New()
		// 		mock.
		// 			ExpectQuery("").
		// 			WillReturnRows(sqlmock.
		// 				NewRows([]string{"uuid", "location", "ctime"}).
		// 				AddRow("0", "row_error", whenwillthenbenow).
		// 				RowError(1, fmt.Errorf("some error")))
		// 		return db
		// 	},
		// 	err:    fmt.Errorf("some error"),
		// },
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {

			result, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectLifecycleIndex(context.Background(), "Test_SelectLifecycleIndex")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectLifecycle(t *testing.T) {
	t.Parallel()

	e0, e1, e2 := types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

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
						NewRows(lcfieldnames).
						AddRow(lctestrow...))
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
						NewRows([]string{"id", "temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "event_severity", "eventtype_name", "stage_uuid", "stage_name", "has_photos", "has_notes"}).
						AddRow(e0.UUID, e0.Temperature, e0.Humidity, e0.MTime, e0.CTime, e0.EventType.UUID, e0.EventType.Name, e0.EventType.Severity, e0.EventType.Stage.UUID, e0.EventType.Stage.Name, 0, 0).
						AddRow(e1.UUID, e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name, 0, 0).
						AddRow(e2.UUID, e2.Temperature, e2.Humidity, e2.MTime, e2.CTime, e2.EventType.UUID, e2.EventType.Name, e2.EventType.Severity, e2.EventType.Stage.UUID, e2.EventType.Stage.Name, 0, 0))

				return db
			},
			id:     "0",
			result: lchappyresults[0],
		},
		"get_events_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(lcfieldnames).
						AddRow(lctestrow...))
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
					WillReturnError(fmt.Errorf("some error"))

				return db
			},
			id: "0",
			result: types.Lifecycle{
				UUID:       "0",
				Location:   "location",
				StrainCost: 0,
				GrainCost:  0,
				BulkCost:   0,
				Yield:      0,
				Count:      0,
				Gross:      0,
				MTime:      whenwillthenbenow,
				CTime:      whenwillthenbenow,
				Strain: types.Strain{
					UUID:    "0",
					Species: "X.species",
					Name:    "strain 0",
					CTime:   whenwillthenbenow,
					Vendor: types.Vendor{
						UUID:    "x",
						Name:    "vendor x",
						Website: "website",
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
						UUID:    "1",
						Name:    "vendor 1",
						Website: "website",
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
						UUID:    "2",
						Name:    "vendor 2",
						Website: "website",
					},
					Ingredients: []types.Ingredient{
						{UUID: "0", Name: "ingredient 0"},
						{UUID: "1", Name: "ingredient 1"},
						{UUID: "2", Name: "ingredient 2"},
					},
				},
				Events: []types.Event{e0, e1, e2},
			},
			err: fmt.Errorf("some error"),
		},
		"all_attrs_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(lcfieldnames).
						AddRow(lctestrow...))
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
						NewRows(lcfieldnames).
						AddRow(lctestrow...))
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
						NewRows(lcfieldnames).
						AddRow(lctestrow...))
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
					WillReturnError(fmt.Errorf("some error"))

				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"no_rows": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(lcfieldnames))
				return db
			},
			id:  "0",
			err: sql.ErrNoRows,
		},

		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
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

			_, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectLifecycle(context.Background(), tc.id, "Test_SelectLifecycle")

			if !assert.Equal(t, tc.err, err) {
				panic(err)
			}
			// require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertLifecycle(t *testing.T) {
	t.Parallel()

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
						NewRows(lcfieldnames).
						AddRow(lctestrow...))
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
						NewRows([]string{"id", "temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "event_severity", "eventtype_name", "stage_uuid", "stage_name", "has_photos", "has_notes"}).
						AddRow(e0.UUID, e0.Temperature, e0.Humidity, e0.MTime, e0.CTime, e0.EventType.UUID, e0.EventType.Name, e0.EventType.Severity, e0.EventType.Stage.UUID, e0.EventType.Stage.Name, 0, 0).
						AddRow(e1.UUID, e1.Temperature, e1.Humidity, e1.MTime, e1.CTime, e1.EventType.UUID, e1.EventType.Name, e1.EventType.Severity, e1.EventType.Stage.UUID, e1.EventType.Stage.Name, 0, 0).
						AddRow(e2.UUID, e2.Temperature, e2.Humidity, e2.MTime, e2.CTime, e2.EventType.UUID, e2.EventType.Name, e2.EventType.Severity, e2.EventType.Stage.UUID, e2.EventType.Stage.Name, 0, 0))

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
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
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

			// there's no good way to test the returned lifecycle, to start, the
			// timestamps are non-deterministic; system tests will vet the rest
			_, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).UpdateLifecycle(context.Background(), types.Lifecycle{}, "Test_UpdateLifecycle")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_UpdateLifecycleMTime(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateLifecycleMTime")

	now := time.Now()

	tcs := map[string]struct {
		db       getMockDB
		modified time.Time
		err      error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			modified: now,
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("mtime was not updated"),
		},
		"row_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		"db_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec("").
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
			}).UpdateLifecycleMTime(context.Background(), &types.Lifecycle{}, time.Now(), "Test_UpdateLifecycle")

			require.Equal(t, tc.err, err)
			// require.Equal(t, tc.modified, lc.MTime)
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

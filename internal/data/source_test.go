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
	"github.com/stretchr/testify/require"
)

var (
	_src = types.Source{
		UUID:      "sourceuuid",
		Type:      "sourcetype",
		Lifecycle: nil,
		Strain:    types.Strain(_strain),
	}
	srcFields = row{"uuid", "type", "progenitor_uuid", "lifecycle_uuid", "strain_uuid", "strain_name", "&strain_species", "strain_ctime", "strain_dtime", "strain_vendor_id", "strain_vendor_name", "strain_vendor_website"}
	srcValues = [][]driver.Value{{_src.UUID, _src.Type, "pgid", nil, _src.Strain.UUID, _src.Strain.Name, _src.Strain.Species, _strain.CTime, _strain.DTime, _strain.Vendor.UUID, _strain.Vendor.Name, _strain.Vendor.Website}}
)

func Test_GetSources(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "GetSources")

	tcs := map[string]struct {
		db     getMockDB
		result []types.Source
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				srcFields.mock(mock, [][]driver.Value{
					{"uuid", "type", "progenitor_uuid", "lifecycle_uuid", "strain_uuid", "strain_name", "strain_species", wwtbn, nil, "strain_vendor_id", "strain_vendor_name", "strain_vendor_website"},
					{"uuid", "type", "progenitor_uuid", nil, "strain_uuid", "strain_name", "strain_species", wwtbn, nil, "strain_vendor_id", "strain_vendor_name", "strain_vendor_website"},
				}...)
				lcFields.mock(mock, []driver.Value{"uuid", "location", 0, 0, 0, 0, 0, 0, wwtbn, wwtbn, "0", "X.species", "strain 0", nil, wwtbn, nil, "x", "vendor x", "website", "gs", "gs", types.GrainType, "1", "vendor 1", "website", "bs", "bs", types.BulkType, "2", "vendor 2", "website"})
				eventFields.mock(mock, eventValues...)
				attrFields.mock(mock, attrValues...)
				ingFields.mock(mock, ingValues...)
				ingFields.mock(mock, ingValues...)

				return db
			},
			result: []types.Source{
				{
					UUID:      "uuid",
					Type:      "type",
					Lifecycle: &types.Lifecycle{UUID: "lifecycle_uuid"},
					Strain:    types.Strain(_strain),
				},
				{
					UUID:   "uuid",
					Type:   "type",
					Strain: types.Strain(_strain),
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
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			g := types.Generation{}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).GetSources(context.Background(), &g, "Test_GetSources")

			require.Equal(t, tc.err, err)
			require.Equal(t, mustObject(tc.result), mustObject(g.Sources))
		})
	}
}

func Test_AddStrainSource(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "AddStrainSource")

	tcs := map[string]struct {
		db     getMockDB
		s      types.Source
		result []types.Source
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))

				strainFields.mock(mock, strainValues)
				attrFields.mock(mock, attrValues...)

				return db
			},
			s: types.Source{Type: "Clone"},
			result: func(strain strain) []types.Source {
				strain.Attributes = []types.StrainAttribute{
					types.StrainAttribute(_attrs[0]),
					types.StrainAttribute(_attrs[1]),
					types.StrainAttribute(_attrs[2]),
				}
				return []types.Source{
					{
						UUID:   types.UUID(mockUUIDGen().String()),
						Type:   "Clone",
						Strain: types.Strain(strain),
					},
				}
			}(_strain),
		},
		"select_strain_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			s:   types.Source{Type: "Clone"},
			err: fmt.Errorf("couldn't fetch strain"),
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("source was not added"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			g := types.Generation{}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).AddStrainSource(
				context.Background(),
				&g,
				tc.s,
				"Test_AddStrainSource")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, g.Sources)
		})
	}
}

func Test_AddEventSource(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "AddEventSource")

	whenwillthenbenow := time.Now() // time.Soon()

	tcs := map[string]struct {
		db     getMockDB
		e      types.Event
		result []types.Source
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"uuid"}).
						AddRow("uuid"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "eventtype_severity", "eventtype_name", "stage_uuid", "stage_name"}).
						AddRow(
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"e0.EventType.UUID",
							"Clone",
							"Info",
							"e0.EventType.Stage.UUID",
							"e0.EventType.Stage.Name"))
				// mock.ExpectQuery("").
				// 	WillReturnRows(sqlmock.
				// 		NewRows([]string{"name", "severity", "stage_uuid", "stage_name"}).
				// 		AddRow("Clone", "Info", "0", "stage 0"))
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"uuid", "species", "name", "ctime", "mtime", "vendor_uuid", "vendor_name", "vendor_website", "generation_uuid"}).
						AddRow("uuid", "X.species", "strain 0", whenwillthenbenow, nil, "0", "vendor 0", "website", nil))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "value"}).
						AddRow("0", "name 0", "value 0").
						AddRow("1", "name 1", "value 1").
						AddRow("2", "name 2", "value 2"))
				return db
			},
			e: types.Event{EventType: types.EventType{UUID: ""}},
			result: []types.Source{
				{
					UUID:      types.UUID(mockUUIDGen().String()),
					Type:      "Clone",
					Lifecycle: nil,
					Strain: types.Strain{
						UUID:    "uuid",
						Name:    "strain 0",
						Species: "X.species",
						CTime:   whenwillthenbenow,
						Vendor: types.Vendor{
							UUID:    "0",
							Name:    "vendor 0",
							Website: "website",
						},
						Attributes: []types.StrainAttribute{
							{UUID: "0", Name: "name 0", Value: "value 0"},
							{UUID: "1", Name: "name 1", Value: "value 1"},
							{UUID: "2", Name: "name 2", Value: "value 2"},
						},
					},
				},
			},
		},
		"strain_id_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("couldn't get strain for AddEventSource (%#v)", types.Event{}),
		},
		"select_eventtype_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"uuid"}).
						AddRow("uuid"))
				mock.ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("couldn't get eventtype for AddEventSource"),
		},
		"select_strain_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"uuid"}).
						AddRow("uuid"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "eventtype_severity", "eventtype_name", "stage_uuid", "stage_name"}).
						AddRow(
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"e0.EventType.UUID",
							"Clone",
							"Info",
							"e0.EventType.Stage.UUID",
							"e0.EventType.Stage.Name"))
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			// e: types.Event{EventType: types.EventType{UUID: ""}},
			err: fmt.Errorf("couldn't fetch strain"),
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"uuid"}).
						AddRow("uuid"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "eventtype_severity", "eventtype_name", "stage_uuid", "stage_name"}).
						AddRow(
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"e0.EventType.UUID",
							"Clone",
							"Info",
							"e0.EventType.Stage.UUID",
							"e0.EventType.Stage.Name"))
				mock.
					ExpectExec("").
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("source was not added"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"uuid"}).
						AddRow("uuid"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "eventtype_severity", "eventtype_name", "stage_uuid", "stage_name"}).
						AddRow(
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"e0.EventType.UUID",
							"Clone",
							"Info",
							"e0.EventType.Stage.UUID",
							"e0.EventType.Stage.Name"))
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
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"uuid"}).
						AddRow("uuid"))
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"temperature", "humidity", "mtime", "ctime", "eventtype_uuid", "eventtype_severity", "eventtype_name", "stage_uuid", "stage_name"}).
						AddRow(
							0,
							0,
							whenwillthenbenow,
							whenwillthenbenow,
							"e0.EventType.UUID",
							"Clone",
							"Info",
							"e0.EventType.Stage.UUID",
							"e0.EventType.Stage.Name"))
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

			g := types.Generation{}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).AddEventSource(
				context.Background(),
				&g,
				tc.e,
				"Test_AddEventSource")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, g.Sources)
		})
	}
}

func Test_ChangeSource(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "ChangeSource")

	tcs := map[string]struct {
		db getMockDB
		s  types.Source
		sources,
		result []types.Source
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
			s: types.Source{UUID: "1", Type: "Spore"},
			sources: []types.Source{
				{UUID: "0"},
				{UUID: "1", Type: "Clone"},
				{UUID: "2"},
			},
			result: []types.Source{
				{UUID: "1", Type: "Spore"},
				{UUID: "0"},
				{UUID: "2"},
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
			err: fmt.Errorf("source was not changed"),
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

			g := types.Generation{Sources: tc.sources}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).ChangeSource(
				context.Background(),
				&g,
				tc.s,
				"Test_ChangeSource")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, g.Sources)
		})
	}
}

func Test_RemoveSource(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "RemoveSource")

	tcs := map[string]struct {
		db getMockDB
		id types.UUID
		sources,
		result []types.Source
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
			sources: []types.Source{
				{UUID: "0"},
				{UUID: "1"},
				{UUID: "2"},
			},
			id: "1",
			result: []types.Source{
				{UUID: "0"},
				{UUID: "2"},
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
			id:  "0",
			err: fmt.Errorf("source could not be deleted: '0'"),
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

			g := types.Generation{Sources: tc.sources}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).RemoveSource(
				context.Background(),
				&g,
				tc.id,
				"Test_RemoveSource")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, g.Sources)
		})
	}
}

package data

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

func Test_SelectGenerationIndex(t *testing.T) {
	// t.Skip()
	t.Parallel()

	whenwillthenbenow := time.Now()
	fields := [25]string{
		"uuid",
		"plating_id",
		"plating_name",
		"plating_type",
		"plating_vendor_id",
		"plating_vendor_name",
		"plating_vendor_website",
		"liquid_id",
		"liquid_name",
		"liquid_type",
		"liquid_vendor_id",
		"liquid_vendor_name",
		"liquid_vendor_website",
		"source_uuid",
		"type",
		"observable_uuid",
		"strain_uuid",
		"strain_name",
		"strain_species",
		"strain_ctime",
		"strain_vendor_id",
		"strain_vendor_name",
		"strain_vendor_website",
		"generation_mtime",
		"generation_ctime",
	}

	l := log.WithField("test", "SelectGenerationIndex")

	tcs := map[string]struct {
		db     getMockDB
		result []types.Generation
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(fields[:]).
						AddRow(
							"happy_path",
							"plating_id",
							"plating_name",
							"plating_type",
							"plating_vendor_id",
							"plating_vendor_name",
							"plating_vendor_website",
							"liquid_id",
							"liquid_name",
							"liquid_type",
							"liquid_vendor_id",
							"liquid_vendor_name",
							"liquid_vendor_website",
							"source_uuid 0",
							"spore",
							"lifecycle_uuid",
							"strain_uuid",
							"strain_name",
							"strain_species",
							whenwillthenbenow,
							"strain_vendor_id",
							"strain_vendor_name",
							"strain_vendor_website",
							whenwillthenbenow,
							whenwillthenbenow).
						AddRow(
							"happy_path",
							"plating_id",
							"plating_name",
							"plating_type",
							"plating_vendor_id",
							"plating_vendor_name",
							"plating_vendor_website",
							"liquid_id",
							"liquid_name",
							"liquid_type",
							"liquid_vendor_id",
							"liquid_vendor_name",
							"liquid_vendor_website",
							"source_uuid 1",
							"spore",
							"lifecycle_uuid",
							"strain_uuid",
							"strain_name",
							"strain_species",
							whenwillthenbenow,
							"strain_vendor_id",
							"strain_vendor_name",
							"strain_vendor_website",
							whenwillthenbenow,
							whenwillthenbenow).
						AddRow(
							"happy_path 2",
							"plating_id",
							"plating_name",
							"plating_type",
							"plating_vendor_id",
							"plating_vendor_name",
							"plating_vendor_website",
							"liquid_id",
							"liquid_name",
							"liquid_type",
							"liquid_vendor_id",
							"liquid_vendor_name",
							"liquid_vendor_website",
							"source_uuid 0",
							"spore",
							"lifecycle_uuid",
							"strain_uuid",
							"strain_name",
							"strain_species",
							whenwillthenbenow,
							"strain_vendor_id",
							"strain_vendor_name",
							"strain_vendor_website",
							whenwillthenbenow,
							whenwillthenbenow))
				return db
			},
			result: []types.Generation{
				{
					UUID:  "happy_path",
					MTime: whenwillthenbenow,
					CTime: whenwillthenbenow,
					PlatingSubstrate: types.Substrate{
						UUID: "plating_id",
						Name: "plating_name",
						Type: "plating_type",
						Vendor: types.Vendor{
							UUID:    "plating_vendor_id",
							Name:    "plating_vendor_name",
							Website: "plating_vendor_website",
						},
					},
					LiquidSubstrate: types.Substrate{
						UUID: "liquid_id",
						Name: "liquid_name",
						Type: "liquid_type",
						Vendor: types.Vendor{
							UUID:    "liquid_vendor_id",
							Name:    "liquid_vendor_name",
							Website: "liquid_vendor_website",
						},
					},
					Sources: []types.Source{
						{
							UUID:      "source_uuid 1",
							Type:      "spore",
							Lifecycle: &types.Lifecycle{UUID: "lifecycle_uuid"},
							Strain: types.Strain{
								UUID:    "strain_uuid",
								Name:    "strain_name",
								Species: "strain_species",
								CTime:   whenwillthenbenow,
								Vendor: types.Vendor{
									UUID:    "strain_vendor_id",
									Name:    "strain_vendor_name",
									Website: "strain_vendor_website",
								},
							},
						},
						{
							UUID:      "source_uuid 0",
							Type:      "spore",
							Lifecycle: &types.Lifecycle{UUID: "lifecycle_uuid"},
							Strain: types.Strain{
								UUID:    "strain_uuid",
								Name:    "strain_name",
								Species: "strain_species",
								CTime:   whenwillthenbenow,
								Vendor: types.Vendor{
									UUID:    "strain_vendor_id",
									Name:    "strain_vendor_name",
									Website: "strain_vendor_website",
								},
							},
						},
					},
				},
				{
					UUID:  "happy_path 2",
					MTime: whenwillthenbenow,
					CTime: whenwillthenbenow,
					PlatingSubstrate: types.Substrate{
						UUID: "plating_id",
						Name: "plating_name",
						Type: "plating_type",
						Vendor: types.Vendor{
							UUID:    "plating_vendor_id",
							Name:    "plating_vendor_name",
							Website: "plating_vendor_website",
						},
					},
					LiquidSubstrate: types.Substrate{
						UUID: "liquid_id",
						Name: "liquid_name",
						Type: "liquid_type",
						Vendor: types.Vendor{
							UUID:    "liquid_vendor_id",
							Name:    "liquid_vendor_name",
							Website: "liquid_vendor_website",
						},
					},
					Sources: []types.Source{
						{
							UUID:      "source_uuid 0",
							Type:      "spore",
							Lifecycle: &types.Lifecycle{UUID: "lifecycle_uuid"},
							Strain: types.Strain{
								UUID:    "strain_uuid",
								Name:    "strain_name",
								Species: "strain_species",
								CTime:   whenwillthenbenow,
								Vendor: types.Vendor{
									UUID:    "strain_vendor_id",
									Name:    "strain_vendor_name",
									Website: "strain_vendor_website",
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
	}

	for k, v := range tcs {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        v.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", k),
			}).SelectGenerationIndex(context.Background(), "Test_SelectGenerationIndex")

			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_SelectGeneration(t *testing.T) {
	t.Parallel()

	fieldnames := []string{
		"platingsubstrate_uuid",
		"platingsubstrate_name",
		"platingsubstrate_type",
		"platingsubstrate_vendor_uuid",
		"platingsubstrate_vendor_name",
		"platingsubstrate_vendor_website",
		"liquidsubstrate_uuid",
		"liquidsubstrate_name",
		"liquidsubstrate_type",
		"liquidsubstrate_vendor_uuid",
		"liquidsubstrate_vendor_name",
		"liquidsubstrate_vendor_website",
		"mtime",
		"ctime",
	}

	e0, e1, e2 := types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

	whenwillthenbenow := time.Now() // time.Soon()

	l := log.WithField("test", "SelectGenerationIndex")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Generation
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(fieldnames).
						AddRow(
							"platingsubstrate_uuid",
							"platingsubstrate_name",
							"platingsubstrate_type",
							"platingsubstrate_vendor_uuid",
							"platingsubstrate_vendor_name",
							"platingsubstrate_vendor_website",
							"liquidsubstrate_uuid",
							"liquidsubstrate_name",
							"liquidsubstrate_type",
							"liquidsubstrate_vendor_uuid",
							"liquidsubstrate_vendor_name",
							"liquidsubstrate_vendor_website",
							whenwillthenbenow,
							whenwillthenbenow,
						))
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
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{
							"uuid",
							"type",
							"pgid",
							"lcid",
							"strain_uuid",
							"strain_name",
							"species",
							"strain_ctime",
							"strain_vendor_uuid",
							"strain_vendor_name",
							"strain_vendor_website",
						}).
						AddRow(
							"uuid",
							"type",
							"pgid",
							nil,
							"strain_uuid",
							"strain_name",
							"species",
							whenwillthenbenow,
							"strain_vendor_uuid",
							"strain_vendor_name",
							"strain_vendor_website",
						))
				return db
			},
			id: "0",
			result: types.Generation{
				Sources: []types.Source{
					{
						UUID: "uuid",
						Type: "type",
						Strain: types.Strain{
							UUID:    "strain_uuid",
							Name:    "strain_name",
							Species: "species",
							CTime:   whenwillthenbenow,
							Vendor: types.Vendor{
								UUID:    "strain_vendor_uuid",
								Name:    "strain_vendor_name",
								Website: "strain_vendor_website",
							},
						},
					},
				},
				MTime: whenwillthenbenow,
				CTime: whenwillthenbenow,
				PlatingSubstrate: types.Substrate{
					UUID: "platingsubstrate_uuid",
					Name: "platingsubstrate_name",
					Type: "platingsubstrate_type",
					Vendor: types.Vendor{
						UUID:    "platingsubstrate_vendor_uuid",
						Name:    "platingsubstrate_vendor_name",
						Website: "platingsubstrate_vendor_website",
					},
				},
				LiquidSubstrate: types.Substrate{
					UUID: "liquidsubstrate_uuid",
					Name: "liquidsubstrate_name",
					Type: "liquidsubstrate_type",
					Vendor: types.Vendor{
						UUID:    "liquidsubstrate_vendor_uuid",
						Name:    "liquidsubstrate_vendor_name",
						Website: "liquidsubstrate_vendor_website",
					},
				},
				Events: []types.Event{e0, e1, e2},
			},
		},
		"plating_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(fieldnames).
						AddRow(
							"platingsubstrate_uuid",
							"platingsubstrate_name",
							"platingsubstrate_type",
							"platingsubstrate_vendor_uuid",
							"platingsubstrate_vendor_name",
							"platingsubstrate_vendor_website",
							"liquidsubstrate_uuid",
							"liquidsubstrate_name",
							"liquidsubstrate_type",
							"liquidsubstrate_vendor_uuid",
							"liquidsubstrate_vendor_name",
							"liquidsubstrate_vendor_website",
							whenwillthenbenow,
							whenwillthenbenow,
						))
				mock.ExpectQuery("").
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"liquid_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(fieldnames).
						AddRow(
							"platingsubstrate_uuid",
							"platingsubstrate_name",
							"platingsubstrate_type",
							"platingsubstrate_vendor_uuid",
							"platingsubstrate_vendor_name",
							"platingsubstrate_vendor_website",
							"liquidsubstrate_uuid",
							"liquidsubstrate_name",
							"liquidsubstrate_type",
							"liquidsubstrate_vendor_uuid",
							"liquidsubstrate_vendor_name",
							"liquidsubstrate_vendor_website",
							whenwillthenbenow,
							whenwillthenbenow,
						))
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
		"events_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(fieldnames).
						AddRow(
							"platingsubstrate_uuid",
							"platingsubstrate_name",
							"platingsubstrate_type",
							"platingsubstrate_vendor_uuid",
							"platingsubstrate_vendor_name",
							"platingsubstrate_vendor_website",
							"liquidsubstrate_uuid",
							"liquidsubstrate_name",
							"liquidsubstrate_type",
							"liquidsubstrate_vendor_uuid",
							"liquidsubstrate_vendor_name",
							"liquidsubstrate_vendor_website",
							whenwillthenbenow,
							whenwillthenbenow,
						))
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
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"source_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows(fieldnames).
						AddRow(
							"platingsubstrate_uuid",
							"platingsubstrate_name",
							"platingsubstrate_type",
							"platingsubstrate_vendor_uuid",
							"platingsubstrate_vendor_name",
							"platingsubstrate_vendor_website",
							"liquidsubstrate_uuid",
							"liquidsubstrate_name",
							"liquidsubstrate_type",
							"liquidsubstrate_vendor_uuid",
							"liquidsubstrate_vendor_name",
							"liquidsubstrate_vendor_website",
							whenwillthenbenow,
							whenwillthenbenow,
						))
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
			}).SelectGeneration(context.Background(), tc.id, "Test_SelectGenerationIndex")

			require.Equal(t, tc.err, err)
			// require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertGeneration(t *testing.T) {
	t.Parallel()

	fieldnames := []string{
		"platingsubstrate_uuid",
		"platingsubstrate_name",
		"platingsubstrate_type",
		"platingsubstrate_vendor_uuid",
		"platingsubstrate_vendor_name",
		"platingsubstrate_vendor_website",
		"liquidsubstrate_uuid",
		"liquidsubstrate_name",
		"liquidsubstrate_type",
		"liquidsubstrate_vendor_uuid",
		"liquidsubstrate_vendor_name",
		"liquidsubstrate_vendor_website",
		"mtime",
		"ctime",
	}

	e0, e1, e2 := types.Event{UUID: "0"},
		types.Event{UUID: "1"},
		types.Event{UUID: "2"}

	whenwillthenbenow := time.Now() // time.Soon()

	l := log.WithField("test", "InsertGeneration")

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
							"platingsubstrate_uuid",
							"platingsubstrate_name",
							"platingsubstrate_type",
							"platingsubstrate_vendor_uuid",
							"platingsubstrate_vendor_name",
							"platingsubstrate_vendor_website",
							"liquidsubstrate_uuid",
							"liquidsubstrate_name",
							"liquidsubstrate_type",
							"liquidsubstrate_vendor_uuid",
							"liquidsubstrate_vendor_name",
							"liquidsubstrate_vendor_website",
							whenwillthenbenow,
							whenwillthenbenow,
						))
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
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{
							"uuid",
							"type",
							"pgid",
							"lcid",
							"strain_uuid",
							"strain_name",
							"species",
							"strain_ctime",
							"strain_vendor_uuid",
							"strain_vendor_name",
							"strain_vendor_website",
						}).
						AddRow(
							"uuid",
							"type",
							"pgid",
							nil,
							"strain_uuid",
							"strain_name",
							"species",
							whenwillthenbenow,
							"strain_vendor_uuid",
							"strain_vendor_name",
							"strain_vendor_website",
						))
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
			err: fmt.Errorf("generation was not added: 0"),
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

			g, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).InsertGeneration(
				context.Background(),
				types.Generation{},
				"Test_InsertGeneration")

			require.Equal(t, tc.err, err)
			require.Equal(t, g.UUID, types.UUID(mockUUIDGen().String()))
		})
	}
}

func Test_UpdateGeneration(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateGeneration")

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
			err: fmt.Errorf("generation was not updated"),
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

			_, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).UpdateGeneration(context.Background(), types.Generation{}, "Test_UpdateGeneration")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_UpdateGenerationMTime(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateGenerationMTime")

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
			}).UpdateGenerationMTime(context.Background(), &types.Generation{}, time.Now(), "Test_UpdateGeneration")

			require.Equal(t, tc.err, err)
			// require.Equal(t, tc.modified, lc.MTime)
		})
	}
}

func Test_DeleteGeneration(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "DeleteGeneration")

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
			err: fmt.Errorf("generation could not be deleted: '0'"),
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
			}).DeleteGeneration(
				context.Background(),
				tc.id,
				"Test_DeleteGeneration")

			require.Equal(t, tc.err, err)
		})
	}
}

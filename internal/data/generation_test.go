package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

var (
	_gen = generation{
		UUID: "uuid",
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
		MTime: wwtbn,
		CTime: wwtbn,
	}
	genFields = row{
		"uuid",
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
		"dtime",
	}
	genValues = []driver.Value{
		_gen.UUID,
		_gen.PlatingSubstrate.UUID,
		_gen.PlatingSubstrate.Name,
		_gen.PlatingSubstrate.Type,
		_gen.PlatingSubstrate.Vendor.UUID,
		_gen.PlatingSubstrate.Vendor.Name,
		_gen.PlatingSubstrate.Vendor.Website,
		_gen.LiquidSubstrate.UUID,
		_gen.LiquidSubstrate.Name,
		_gen.LiquidSubstrate.Type,
		_gen.LiquidSubstrate.Vendor.UUID,
		_gen.LiquidSubstrate.Vendor.Name,
		_gen.LiquidSubstrate.Vendor.Website,
		_gen.MTime,
		_gen.CTime,
		nil,
	}
)

func Test_SelectGenerationIndex(t *testing.T) {
	t.Parallel()

	fields := [26]string{
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
		"generation_dtime",
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
				mock.ExpectQuery("").WillReturnRows(sqlmock.
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
						wwtbn,
						"strain_vendor_id",
						"strain_vendor_name",
						"strain_vendor_website",
						wwtbn,
						wwtbn,
						nil).
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
						wwtbn,
						"strain_vendor_id",
						"strain_vendor_name",
						"strain_vendor_website",
						wwtbn,
						wwtbn,
						nil).
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
						wwtbn,
						"strain_vendor_id",
						"strain_vendor_name",
						"strain_vendor_website",
						wwtbn,
						wwtbn,
						nil))
				return db
			},
			result: []types.Generation{
				{
					UUID:  "happy_path",
					MTime: wwtbn,
					CTime: wwtbn,
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
							Strain:    types.Strain(_strain),
						},
						{
							UUID:      "source_uuid 0",
							Type:      "spore",
							Lifecycle: &types.Lifecycle{UUID: "lifecycle_uuid"},
							Strain:    types.Strain(_strain),
						},
					},
				},
				{
					UUID:  "happy_path 2",
					MTime: wwtbn,
					CTime: wwtbn,
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
							Strain:    types.Strain(_strain),
						},
					},
				},
			},
		},
		"db_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
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
			require.Equal(t, mustObject(v.result), mustObject(result))
		})
	}
}

func Test_SelectGenerationsByStrain(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectGenerationsByStrain")

	tcs := map[string]struct {
		db     getMockDB
		result []types.Generation
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				genFields.mock(mock, genValues)
				ingFields.mock(mock)
				ingFields.mock(mock)
				eventFields.mock(mock)
				srcFields.mock(mock)

				return db
			},
			result: []types.Generation{
				{
					UUID:  "uuid",
					MTime: wwtbn,
					CTime: wwtbn,
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
					Events: []types.Event{},
				},
			},
		},
		"db_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
	}

	for k, v := range tcs {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			p, _ := types.NewReportAttrs(map[string][]string{"strain-id": {"0"}})

			result, err := (&Conn{
				query:        v.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", k),
			}).selectGenerations(context.Background(), p, "Test_SelectGenerationsByStrain")

			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_SelectGeneration(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectGeneration")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Generation
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				genFields.mock(mock, genValues)
				eventFields.mock(mock, eventValues...)
				srcFields.mock(mock, srcValues...)

				return db
			},
			id: "0",
			result: types.Generation{
				Sources: []types.Source{
					{
						UUID:   "uuid",
						Type:   "type",
						Strain: types.Strain(_strain),
					},
				},
				MTime: wwtbn,
				CTime: wwtbn,
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
				Events: []types.Event{
					types.Event(_events[0]),
					types.Event(_events[1]),
					types.Event(_events[2]),
				},
			},
		},
		"events_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				genFields.mock(mock, genValues)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"source_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				genFields.mock(mock, genValues)
				eventFields.mock(mock, eventValues...)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"no_rows": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				genFields.mock(mock)
				return db
			},
			id:  "0",
			err: sql.ErrNoRows,
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
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
			}).SelectGeneration(context.Background(), tc.id, "Test_SelectGeneration")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result.UUID, types.UUID(""))
		})
	}
}

func Test_InsertGeneration(t *testing.T) {
	t.Parallel()

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
				genFields.mock(mock, append([]driver.Value{mockUUIDGen().String()}, genValues[1:]...))
				eventFields.mock(mock, eventValues...)
				srcFields.mock(mock, srcValues...)

				return db
			},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("generation was not added: 0"),
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

			g, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).InsertGeneration(
				context.Background(),
				types.Generation{},
				"Test_InsertGeneration")

			require.Equal(t, tc.err, err)
			require.Equal(t, types.UUID(mockUUIDGen().String()), g.UUID)
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
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("generation was not updated"),
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
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			modified: now,
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("mtime was not updated"),
		},
		"row_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		"db_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
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
		err error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			// id: "0",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			err: fmt.Errorf("generation could not be deleted: 'tc.id'"),
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

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).DeleteGeneration(
				context.Background(),
				"tc.id",
				"Test_DeleteGeneration")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_GenerationReport(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "GenerationReport")

	tcs := map[string]struct {
		db     getMockDB
		result types.Entity
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				newBuilder(mock,
					genFields.set(genValues),
					eventFields.set(eventValues...),
					srcFields.set(srcValues...),
					ingFields.set(ingValues...),
					ingFields.set(ingValues...),
					napFields.set(),
					noteFields.set(noteValues...),
					strainFields.set())

				return db
			},
		},
		"progeny_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				genFields.mock(mock, genValues)
				eventFields.mock(mock, eventValues...)
				// also used different values, like above
				srcFields.mock(mock, srcValues...)
				ingFields.mock(mock, ingValues...)
				ingFields.mock(mock, ingValues...)
				napFields.mock(mock)
				noteFields.mock(mock, noteValues...)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		"notes_error": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				genFields.mock(mock, genValues)
				eventFields.mock(mock, eventValues...)
				srcFields.mock(mock, srcValues...)
				ingFields.mock(mock, ingValues...)
				ingFields.mock(mock, ingValues...)
				napFields.mock(mock)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		"plating_ingredient_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				genFields.mock(mock, genValues)
				eventFields.mock(mock, eventValues...)
				srcFields.mock(mock, srcValues...)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		"liquid_ingredient_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				genFields.mock(mock, genValues)
				eventFields.mock(mock, eventValues...)
				srcFields.mock(mock, srcValues...)
				ingFields.mock(mock, ingValues...)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		"notes_and_photos_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				genFields.mock(mock, genValues)
				eventFields.mock(mock, eventValues...)
				srcFields.mock(mock, srcValues...)
				ingFields.mock(mock, ingValues...)
				ingFields.mock(mock, ingValues...)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		// "no_id": {
		// 	db: func() *sql.DB {
		// 		db, _, _ := sqlmock.New()
		// 		return db
		// 	},
		// 	err: fmt.Errorf("failed to find param values in the following fields: [generation-id]"),
		// },
		"events_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				genFields.mock(mock, genValues)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))

				return db
			},
			err: fmt.Errorf("some error"),
		},
		"source_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				genFields.mock(mock, genValues)
				eventFields.mock(mock, eventValues...)

				mock.ExpectQuery("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			err: fmt.Errorf("some error"),
		},
		"no_rows": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				genFields.mock(mock)
				return db
			},
			err: sql.ErrNoRows,
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
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

			_, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).GenerationReport(context.Background(), "tc.id", "Test_GenerationReport")

			require.Equal(t, tc.err, err)
			// require.Equal(t, tc.result.UUID, types.UUID(""))
		})
	}
}

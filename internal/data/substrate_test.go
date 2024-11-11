package data

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/jsmit257/huautla/types"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/require"
)

var (
	_subs = []substrate{
		{UUID: "0", Name: "substrate 0", Type: types.PlatingType, Vendor: types.Vendor{UUID: "0", Name: "vendor 0", Website: "website 0"}, Ingredients: []types.Ingredient{}},
		{UUID: "1", Name: "substrate 1", Type: types.GrainType, Vendor: types.Vendor{UUID: "1", Name: "vendor 1", Website: "website 1"}, Ingredients: []types.Ingredient{}},
		{UUID: "2", Name: "substrate 2", Type: types.GrainType, Vendor: types.Vendor{UUID: "1", Name: "vendor 1", Website: "website 1"}, Ingredients: []types.Ingredient{}},
	}
	subFields = row{"id", "name", "type", "vendor_uuid", "vendor_name", "vendor_website"}
	subValues = [][]driver.Value{
		{_subs[0].UUID, _subs[0].Name, _subs[0].Type, _subs[0].Vendor.UUID, _subs[0].Vendor.Name, _subs[0].Vendor.Website},
		{_subs[1].UUID, _subs[1].Name, _subs[1].Type, _subs[1].Vendor.UUID, _subs[1].Vendor.Name, _subs[1].Vendor.Website},
		{_subs[2].UUID, _subs[2].Name, _subs[2].Type, _subs[2].Vendor.UUID, _subs[2].Vendor.Name, _subs[2].Vendor.Website},
	}
)

func Test_SelectAllSubstrates(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectAllSubstrates")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.Substrate
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					subFields.set(subValues...),
					ingFields.set(),
					ingFields.set(),
					ingFields.set())

				return db
			},
			result: []types.Substrate{
				types.Substrate(_subs[0]),
				types.Substrate(_subs[1]),
				types.Substrate(_subs[2]),
			},
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
			}).SelectAllSubstrates(context.Background(), "Test_SelectAllSubstrates")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectSubstrate(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectSubstrate")

	tcs := map[string]struct {
		db     getMockDB
		noid   bool
		result types.Substrate
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					subFields.set(subValues[0]),
					ingFields.set(ingValues...))

				return db
			},
			result: types.Substrate{
				UUID:   _subs[0].UUID,
				Name:   _subs[0].Name,
				Type:   _subs[0].Type,
				Vendor: _subs[0].Vendor,
				Ingredients: []types.Ingredient{
					types.Ingredient(_ingredients[0]),
					types.Ingredient(_ingredients[1]),
					types.Ingredient(_ingredients[2]),
				}},
		},
		"ingredients_fail": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {

				newBuilder(mock,
					subFields.set(subValues[0]),
					ingFields.fail())

				return db
			},
			err: ingFields.err(),
		},
		"subs_empty": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, subFields.set())
				return db
			},
			err: sql.ErrNoRows,
		},
		"missing_id": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, subFields.set())
				return db
			},
			noid: true,
			err:  fmt.Errorf("failed to find param values in the following fields: [substrate-id]"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
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

			var id types.UUID
			if !tc.noid {
				id = "0"
			}

			result, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectSubstrate(context.Background(), id, "Test_SelectSubstrate")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertSubstrate(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "InsertSubstrate")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		tp     types.SubstrateType
		result types.Substrate
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}},
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}},
			err:    fmt.Errorf("substrate was not added"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}, Ingredients: nil},
			err:    fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			id:     "0",
			tp:     types.GrainType,
			result: types.Substrate{UUID: "30313233-3435-3637-3839-616263646566", Name: "substrate 0", Type: types.GrainType, Vendor: types.Vendor{}, Ingredients: nil},
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
			}).InsertSubstrate(
				context.Background(),
				types.Substrate{Name: "substrate " + string(tc.id), Type: tc.tp, Vendor: types.Vendor{}, Ingredients: nil},
				"Test_InsertSubstrate")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateSubstrate(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateSubstrate")

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
			err: fmt.Errorf("substrate was not updated: '0'"),
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
			}).UpdateSubstrate(
				context.Background(),
				tc.id,
				types.Substrate{Name: "substrate " + string(tc.id)},
				"Test_UpdateSubstrates")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_DeleteSubstrate(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "DeleteSubstrate")

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
			err: fmt.Errorf("substrate could not be deleted: '0'"),
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
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).DeleteSubstrate(
				context.Background(),
				tc.id,
				"Test_DeleteSubstrate")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_SubstrateReport(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SubstrateReport")

	tcs := map[string]struct {
		db     getMockDB
		noid   bool
		result types.Entity
		err    error
	}{
		"happy_generation_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					subFields.set(subValues[0]),
					ingFields.set(ingValues...),
					genFields.set(genValues),
					eventFields.set(eventValues...),
					srcFields.set(srcValues[0]),
					ingFields.set(ingValues...),
					ingFields.set(ingValues...),
					napFields.set(),
					noteFields.set(noteValues...),
					strainFields.set())

				return db
			},

			result: func(s types.Entity) types.Entity {

				g := mustEntity(_gen)

				g["events"] = events
				g["plating_substrate"].(map[string]interface{})["ingredients"] = ingredients
				g["liquid_substrate"].(map[string]interface{})["ingredients"] = ingredients
				g["notes"] = notes
				g["sources"] = []interface{}{
					map[string]interface{}{
						"id":     _src.UUID,
						"strain": mustEntity(_src.Strain),
						"type":   _src.Type,
					},
				}

				s["generations"] = []types.Entity{g}
				s["ingredients"] = ingredients

				return s
			}(mustEntity(_subs[0])),
		},
		"generation_path_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					subFields.set(subValues[0]),
					ingFields.set(ingValues...),
					genFields.fail())

				return db
			},
			err: genFields.err(),
		},
		"happy_lifecycle_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					subFields.set(subValues[1]),
					ingFields.set(),
					lcFields.set(lcValues),
					eventFields.set(),
					attrFields.set(),
					ingFields.set(),
					ingFields.set(),
					napFields.set(),
					noteFields.set(),
					genFields.set())

				return db
			},
			result: func(s types.Entity) types.Entity {
				s["lifecycles"] = []types.Entity{mustEntity(_lc)}
				return s
			}(mustEntity(_subs[1])),
		},
		"lifecycle_path_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					subFields.set(subValues[1]),
					ingFields.set(),
					lcFields.fail())

				return db
			},
			err: lcFields.err(),
		},
		"happy_path_no_children": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock,
					subFields.set(subValues[0]),
					ingFields.set(ingValues...),
					genFields.set())

				return db
			},
			result: func(s types.Entity) types.Entity {
				s["ingredients"] = ingredients
				return s
			}(mustEntity(_subs[0])),
		},
		"no_rows": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, subFields.set())
				return db
			},
			err: sql.ErrNoRows,
		},
		"missing_id": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, subFields.set())
				return db
			},
			noid: true,
			err:  fmt.Errorf("failed to find param values in the following fields: [substrate-id]"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				newBuilder(mock, subFields.fail())
				return db
			},
			err: subFields.err(),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var id types.UUID
			if !tc.noid {
				id = "0"
			}

			result, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SubstrateReport(context.Background(), id, "Test_SubstrateReport")

			js, _ := json.Marshal(result)

			require.Equal(t, tc.err, err)
			require.Equal(t, mustObject(tc.result), mustObject(result), string(js))
		})
	}
}

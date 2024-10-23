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
	_ven = vendor{
		UUID:    "vendoruuid",
		Name:    "vendorname",
		Website: "vendorwebsite",
	}
	venFields = row{"uuid", "name", "website"}
	venValue  = []driver.Value{_ven.UUID, _ven.Name, _ven.Website}
)

func Test_SelectAllVendors(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectAllVendors")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.Vendor
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				venFields.mock(mock, [][]driver.Value{
					{"0", "vendor 0", "website"},
					{"1", "vendor 1", "website"},
					{"2", "vendor 2", "website"},
				}...)

				return db
			},
			result: []types.Vendor{
				{UUID: "0", Name: "vendor 0", Website: "website"},
				{UUID: "1", Name: "vendor 1", Website: "website"},
				{UUID: "2", Name: "vendor 2", Website: "website"},
			},
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

			result, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectAllVendors(context.Background(), "Test_SelectAllVendors")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_SelectVendor(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "SelectVendor")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result types.Vendor
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				venFields.mock(mock, venValue)
				return db
			},
			id:     "vendoruuid",
			result: types.Vendor(_ven),
		},
		"no_result": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				venFields.mock(mock)
				return db
			},
			result: types.Vendor{},
			err:    sql.ErrNoRows,
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

			result, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).SelectVendor(context.Background(), tc.id, "Test_SelectVendors")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_InsertVendor(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "InsertVendor")

	tcs := map[string]struct {
		db     getMockDB
		result types.Vendor
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			result: types.Vendor{UUID: "30313233-3435-3637-3839-616263646566", Name: "vendor 0"},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			result: types.Vendor{UUID: "30313233-3435-3637-3839-616263646566", Name: "vendor 0"},
			err:    fmt.Errorf("vendor was not added"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			result: types.Vendor{UUID: "30313233-3435-3637-3839-616263646566", Name: "vendor 0"},
			err:    fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			result: types.Vendor{UUID: "30313233-3435-3637-3839-616263646566", Name: "vendor 0"},
			err:    fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).InsertVendor(
				context.Background(),
				types.Vendor{UUID: "0", Name: "vendor 0"},
				"Test_InsertVendors")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

func Test_UpdateVendor(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "UpdateVendor")

	tcs := map[string]struct {
		db  getMockDB
		id  types.UUID
		v   types.Vendor
		err error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("vendor was not updated: '0'"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
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
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).UpdateVendor(
				context.Background(),
				tc.id,
				types.Vendor{},
				// types.Vendor{Name: "vendor " + string(tc.id)},
				"Test_UpdateVendors")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_DeleteVendor(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "DeleteVendor")

	tcs := map[string]struct {
		db  getMockDB
		id  types.UUID
		err error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("vendor could not be deleted: '0'"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			id:  "0",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
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
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).DeleteVendor(
				context.Background(),
				tc.id,
				"Test_DeleteVendors")

			require.Equal(t, tc.err, err)
		})
	}
}

func Test_VendorReport(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "VendorReport")

	tcs := map[string]struct {
		db     getMockDB
		result types.Entity
		err    error
	}{
		"happy_strain_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				newBuilder(mock,
					venFields.set(venValue),
					subFields.set(),
					strainFields.set(strainValues),
					attrFields.set(),
					// add a generation
					genFields.set(genValues),
					eventFields.set(),
					srcFields.set(),
					napFields.set(),
					ingFields.set(),
					noteFields.set(),
					strainFields.set(),
					// add a lifecycle
					lcFields.set(lcValues),
					eventFields.set(),
					attrFields.set(),
					ingFields.set(),
					ingFields.set(),
					napFields.set(),
					noteFields.set(),
					photoFields.set())

				return db
			},
			result: func(v types.Entity) types.Entity {
				str := mustEntity(_strain)
				str["lifecycles"] = []types.Entity{mustEntity(_lc)}
				str["generations"] = []types.Entity{mustEntity(_gen)}

				v["strains"] = []types.Entity{str}

				return v
			}(mustEntity(_ven)),
		},
		"strain_path_fail": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				newBuilder(mock,
					venFields.set(venValue),
					subFields.set(),
					strainFields.set(strainValues),
					// add a generation
					genFields.fail())

				return db
			},
			err: genFields.err(),
		},
		"happy_substrate_path_for_lifecycles": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				newBuilder(mock,
					venFields.set(venValue),
					subFields.set(subValues[1]),
					ingFields.set(),
					// create a lifecycle from substrate
					lcFields.set(lcValues),
					eventFields.set(),
					attrFields.set(),
					ingFields.set(),
					ingFields.set(),
					napFields.set(),
					noteFields.set(),
					genFields.set(),
					// no strains for this path
					strainFields.set())

				return db
			},
			result: func(v types.Entity) types.Entity {
				sub := mustEntity(_subs[1])
				sub["lifecycles"] = []types.Entity{mustObject(_lc)}

				v["substrates"] = []types.Entity{sub}

				return v
			}(mustEntity(_ven)),
		},
		"substrate_path_for_lifecycles_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				newBuilder(mock,
					venFields.set(venValue),
					subFields.set(subValues[1]),
					ingFields.set(),
					// create a lifecycle from substrate
					lcFields.set(lcValues),
					eventFields.fail())

				return db
			},
			err: eventFields.err(),
		},
		"happy_substrate_path_for_generations": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				newBuilder(mock,
					venFields.set(venValue),
					subFields.set(subValues[0]),
					ingFields.set(),
					// create a generation from substrate
					genFields.set(genValues),
					eventFields.set(),
					srcFields.set(srcValues[0]),
					ingFields.set(),
					ingFields.set(),
					napFields.set(),
					noteFields.set(),
					strainFields.set(),
					// no strains for this path
					strainFields.set())

				return db
			},
			result: func(v types.Entity) types.Entity {
				gen := mustEntity(_gen)
				gen["sources"] = []interface{}{mustObject(_src)}

				sub := mustEntity(_subs[0])
				sub["generations"] = []types.Entity{gen}

				v["substrates"] = []types.Entity{sub}

				return v
			}(mustEntity(_ven)),
		},
		"substrate_path_for_generations_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()

				newBuilder(mock,
					venFields.set(venValue),
					subFields.set(subValues[0]),
					ingFields.set(),
					// create a generation from substrate
					genFields.set(genValues),
					eventFields.fail(),
				)

				return db
			},
			err: eventFields.err(),
		},
		"happy_path_no_children": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				newBuilder(mock, venFields.set(venValue), subFields.set(), strainFields.set())
				return db
			},
			result: mustEntity(_ven),
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

			result, err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).VendorReport(context.Background(), "tc.id", "Test_VendorReport")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}

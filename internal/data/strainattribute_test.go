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
	_attrs = []attr{
		{UUID: "attruuid 0", Name: "attrname 0", Value: "attrvalue 0"},
		{UUID: "attruuid 1", Name: "attrname 1", Value: "attrvalue 1"},
		{UUID: "attruuid 2", Name: "attrname 2", Value: "attrvalue 2"},
	}
	attrFields = row{"uuid", "name", "value"}
	attrValues = [][]driver.Value{
		{_attrs[0].UUID, _attrs[0].Name, _attrs[0].Value},
		{_attrs[1].UUID, _attrs[1].Name, _attrs[1].Value},
		{_attrs[2].UUID, _attrs[2].Name, _attrs[2].Value},
	}
)

func Test_KnownAttributeNames(t *testing.T) {
	//ctx context.Context, cid types.CID) ([]string, error)
	t.Parallel()

	l := log.WithField("test", "KnownAttributeNames")

	tcs := map[string]struct {
		db     getMockDB
		result []string
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"name"}).
						AddRow("name 0").
						AddRow("name 1").
						AddRow("name 2"))
				return db
			},
			result: []string{"name 0", "name 1", "name 2"},
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

			s, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).KnownAttributeNames(context.Background(), "Test_KnownAttributeNames")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, s)
		})
	}
}

func Test_GetAllAttributes(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "GetAllAttributes")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.StrainAttribute
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectQuery("").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "value"}).
						AddRow("0", "name 0", "value 0").
						AddRow("1", "name 1", "value 1").
						AddRow("2", "name 2", "value 2"))
				return db
			},
			result: []types.StrainAttribute{
				{UUID: "0", Name: "name 0", Value: "value 0"},
				{UUID: "1", Name: "name 1", Value: "value 1"},
				{UUID: "2", Name: "name 2", Value: "value 2"},
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

			s := &types.Strain{UUID: tc.id}

			err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).GetAllAttributes(context.Background(), s, "Test_GetAllAttributes")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, s.Attributes)
		})
	}
}

func Test_AddAttribute(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "InsertStrain")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		a      types.StrainAttribute
		result []types.StrainAttribute
		err    error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
			result: []types.StrainAttribute{
				{UUID: "30313233-3435-3637-3839-616263646566"},
			},
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("attribute was not added"),
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

			s := &types.Strain{}

			a, err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).AddAttribute(
				context.Background(),
				s,
				types.StrainAttribute{},
				"Test_InsertStrains")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, s.Attributes)
			require.NotEmpty(t, a)
		})
	}
}

func Test_ChangeAttribute(t *testing.T) {
	//ctx context.Context, s *Strain, n, v string, cid CID) error
	t.Parallel()

	l := log.WithField("test", "RemoveAttribute")

	tcs := map[string]struct {
		db    getMockDB
		attrs []types.StrainAttribute
		id    types.UUID
		n, v  string
		err   error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			attrs: []types.StrainAttribute{
				{UUID: "0", Name: "Mojo", Value: "Lost"},
				{UUID: "1", Name: "Yield", Value: "Some"},
			},
			id: "1",
			n:  "Yield",
			v:  "Lots!!",
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			n:   "Yield",
			v:   "Lots!!",
			err: fmt.Errorf("attribute was not changed"),
		},
		"query_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnError(fmt.Errorf("some error"))
				return db
			},
			n:   "Yield",
			v:   "Lots!!",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
				return db
			},
			n:   "Yield",
			v:   "Lots!!",
			err: fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := &types.Strain{UUID: "0", Attributes: tc.attrs}

			err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).ChangeAttribute(
				context.Background(),
				s,
				types.StrainAttribute{UUID: tc.id, Name: tc.n, Value: tc.v},
				"Test_RemoveAttribute")

			require.Equal(t, tc.err, err)
		})
	}

}

func Test_RemoveAttribute(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "RemoveAttribute")

	tcs := map[string]struct {
		db    getMockDB
		attrs []types.StrainAttribute
		id    types.UUID
		err   error
	}{
		"happy_path": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			attrs: []types.StrainAttribute{
				{UUID: "1"},
				{UUID: "0"},
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func(db *sql.DB, mock sqlmock.Sqlmock, err error) *sql.DB {
				mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("attribute was not removed"),
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

			s := &types.Strain{Attributes: tc.attrs}

			err := (&Conn{
				query:        tc.db(sqlmock.New()),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).RemoveAttribute(
				context.Background(),
				s,
				tc.id,
				"Test_RemoveAttribute")

			require.Equal(t, tc.err, err)
		})
	}
}

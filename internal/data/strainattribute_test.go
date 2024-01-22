package data

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jsmit257/huautla/types"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_KnownAttributeNames(t *testing.T) {
	//ctx context.Context, cid types.CID) ([]string, error)
	t.Parallel()

	querypat, l := sqls["get-unique-names"],
		log.WithField("test", "KnownAttributeNames")

	tcs := map[string]struct {
		db     getMockDB
		result []string
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery(querypat).
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
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery(querypat).
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			result: []string{},
			err:    fmt.Errorf("some error"),
		},
		// "query_result_nil": {}, // FIXME: how to mock?
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s, err := (&Conn{
				query:        tc.db(),
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

	querypat, l := sqls["all-attributes"],
		log.WithField("test", "GetAllAttributes")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.StrainAttribute
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery(querypat).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name", "value"}).
						AddRow("0", "name 0", "value 0").
						AddRow("1", "name 1", "value 1").
						AddRow("2", "name 2", "value 2"))
				return db
			},
			result: []types.StrainAttribute{
				{"0", "name 0", "value 0"},
				{"1", "name 1", "value 1"},
				{"2", "name 2", "value 2"},
			},
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectQuery(querypat).
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			result: []types.StrainAttribute{},
			err:    fmt.Errorf("some error"),
		},
		// "query_result_nil": {}, // FIXME: how to mock?
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := &types.Strain{UUID: tc.id}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).GetAllAttributes(context.Background(), s, "Test_GetAllAttributes")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, s.Attributes)
		})
	}
}

func Test_AddAttribute(t *testing.T) {
	//ctx context.Context, s *Strain, sa StrainAttribute, cid CID) error
	t.Parallel()

	var querypat = sqls["insert"]

	l := log.WithField("test", "InsertStrain")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		n, v   string
		result []types.StrainAttribute
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			id: "0",
			result: []types.StrainAttribute{
				types.StrainAttribute{UUID: "30313233-3435-3637-3839-616263646566"},
			},
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("attribute was not added"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
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
					ExpectExec(querypat).
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

			s := &types.Strain{}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).AddAttribute(
				context.Background(),
				s,
				tc.n,
				tc.v,
				"Test_InsertStrains")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, s.Attributes)
		})
	}
}

func Test_ChangeAttribute(t *testing.T) {
	//ctx context.Context, s *Strain, n, v string, cid CID) error
	t.Parallel()

	var querypat = sqls["delete"]

	l := log.WithField("test", "RemoveAttribute")

	tcs := map[string]struct {
		db    getMockDB
		attrs []types.StrainAttribute
		id    types.UUID
		n, v  string
		err   error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			attrs: []types.StrainAttribute{
				{"0", "Mojo", "Lost"},
				{"1", "Yield", "Some"},
			},
			id: "0",
			n:  "Yield",
			v:  "Lots!!",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			n:   "Yield",
			v:   "Lots!!",
			err: fmt.Errorf("attribute was not changed"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnError(fmt.Errorf("some error"))
				return db
			},
			n:   "Yield",
			v:   "Lots!!",
			err: fmt.Errorf("some error"),
		},
		"result_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("some error")))
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

			s := &types.Strain{UUID: tc.id, Attributes: tc.attrs}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).ChangeAttribute(
				context.Background(),
				s,
				tc.id,
				tc.n,
				tc.v,
				"Test_RemoveAttribute")

			require.Equal(t, tc.err, err)
		})
	}

}

func Test_RemoveAttribute(t *testing.T) {
	//ctx context.Context, s *Strain, id UUID, cid CID) error
	t.Parallel()

	var querypat = sqls["delete"]

	l := log.WithField("test", "RemoveAttribute")

	tcs := map[string]struct {
		db    getMockDB
		attrs []types.StrainAttribute
		id    types.UUID
		err   error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 1))
				return db
			},
			attrs: []types.StrainAttribute{
				{UUID: "1"},
				{UUID: "0"},
			},
			id: "0",
		},
		"no_rows_affected": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
					WillReturnResult(sqlmock.NewResult(0, 0))
				return db
			},
			id:  "0",
			err: fmt.Errorf("attribute was not removed"),
		},
		"query_fails": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.
					ExpectExec(querypat).
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
					ExpectExec(querypat).
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

			s := &types.Strain{Attributes: tc.attrs}

			err := (&Conn{
				query:        tc.db(),
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

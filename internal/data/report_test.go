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

func Test_newRpt(t *testing.T) {
	t.Parallel()

	l := log.WithField("test", "newRpt")

	tcs := map[string]struct {
		db     getMockDB
		data   any
		parent *rpttree
		err    error
	}{
		"happy_path_no_children": {
			db: func() *sql.DB {
				db, _, _ := sqlmock.New()
				return db
			},
			data: types.Vendor{UUID: "0"},
		},
		"has_cycle": {
			db: func() *sql.DB {
				db, _, _ := sqlmock.New()
				return db
			},
			parent: &rpttree{id: "vendor#0"},
			data:   types.Vendor{UUID: "0"},
		},
		"no_cycle": {
			db: func() *sql.DB {
				db, _, _ := sqlmock.New()
				return db
			},
			parent: &rpttree{id: "strain#0"},
			data:   types.Vendor{UUID: "0"},
		},
		"type_error": {
			db: func() *sql.DB {
				db, _, _ := sqlmock.New()
				return db
			},
			data: struct{}{},
			err:  fmt.Errorf("couldn't determine entity type: '{}' 'struct {}'"),
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
			}).newRpt(context.Background(), tc.data, "Test_newRpt", tc.parent)
			require.Equal(t, tc.err, err)
		})
	}
}

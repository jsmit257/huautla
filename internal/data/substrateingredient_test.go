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

func Test_GetAllIngredients(t *testing.T) {
	t.Parallel()

	querypat, l := sqls["all-ingredients"],
		log.WithField("test", "GetAllIngredients")

	tcs := map[string]struct {
		db     getMockDB
		id     types.UUID
		result []types.Ingredient
		err    error
	}{
		"happy_path": {
			db: func() *sql.DB {
				db, mock, _ := sqlmock.New()
				mock.ExpectQuery(querypat).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "name"}).
						AddRow("0", "ingredient 0").
						AddRow("1", "ingredient 1").
						AddRow("2", "ingredient 2"))
				return db
			},
			result: []types.Ingredient{
				types.Ingredient{"0", "ingredient 0"},
				types.Ingredient{"1", "ingredient 1"},
				types.Ingredient{"2", "ingredient 2"},
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
			result: []types.Ingredient{},
			err:    fmt.Errorf("some error"),
		},
	}

	for name, tc := range tcs {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := &types.Substrate{UUID: tc.id}

			err := (&Conn{
				query:        tc.db(),
				generateUUID: mockUUIDGen,
				logger:       l.WithField("name", name),
			}).GetAllIngredients(context.Background(), s, "Test_GetAllIngredients")

			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, s.Ingredients)
		})
	}
}

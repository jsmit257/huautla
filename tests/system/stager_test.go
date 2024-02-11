package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var stages = []types.Stage{
	{UUID: "0", Name: "Gestation"},
	{UUID: "1", Name: "Colonization"},
	{UUID: "2", Name: "Majority"},
	{UUID: "3", Name: "Vacation"},
}

func Test_SelectAllStages(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		result []types.Stage
		err    error
	}{
		"happy_path": {
			result: stages,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectAllStages(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result[0:len(v.result)])
		})
	}
}
func Test_SelectStage(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		result types.Stage
		err    error
	}{
		"happy_path": {
			id:     stages[2].UUID,
			result: stages[2],
		},
		"no_row_returned": {
			id:     "8",
			result: types.Stage{UUID: "8", Name: ""},
			err:    fmt.Errorf("sql: no rows in result set"),
		},
		"query_fails": {
			id:     invalidUUID,
			result: types.Stage{UUID: "0", Name: ""},
			err:    fmt.Errorf("sql: no rows in result set"), // XXX
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectStage(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}
func Test_InsertStage(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s   types.Stage
		err error
	}{
		"happy_path": {
			s: types.Stage{Name: "bogus!"},
		},
		"duplicate_name_violation": {
			s:   stages[0],
			err: fmt.Errorf(""),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertStage(context.Background(), v.s, types.CID(k))
			require.Equal(t, v.err, err)
			require.NotEmpty(t, result.UUID)
		})
	}
}
func Test_UpdateStage(t *testing.T) {
	set := map[string]struct {
		id  types.UUID
		s   types.Stage
		err error
	}{
		"happy_path": {
			id: "update me!",
			s:  types.Stage{Name: "Renamed"},
		},
		"no_rows_affected": {
			id:  "missing",
			err: fmt.Errorf("stage was not updated: 'missing'"),
		},
		"duplicate_name_violation": {
			id: "update me!",
			s:  stages[0],
		},
		"query_fails": {
			id: invalidUUID,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.UpdateStage(context.Background(), v.id, v.s, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
func Test_DeleteStage(t *testing.T) {
	set := map[string]struct {
		id  types.UUID
		err error
	}{
		"happy_path": {
			id: "delete me!",
		},
		"no_rows_affected": {
			id:  "missing",
			err: fmt.Errorf("stage could not be deleted: 'missing'"),
		},
		"referential_violation": {
			id:  stages[0].UUID,
			err: fmt.Errorf("referential constraint"),
		},
		"query_fails": {
			id:  invalidUUID,
			err: fmt.Errorf(""),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteStage(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}

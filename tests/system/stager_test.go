package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var stages = map[string]types.Stage{}

func init() {
	for _, id := range []types.UUID{"0", "1", "2", "3", "4"} {
		if s, err := db.SelectStage(context.Background(), id, "stager_init"); err != nil {
			panic(err)
		} else {
			stages[s.Name] = s
		}
	}
}

func Test_SelectAllStages(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		result []types.Stage
		err    error
	}{
		"happy_path": {
			result: append(
				[]types.Stage{stages["Gestation"]},
				stages["Colonization"],
				stages["Majority"],
				stages["Vacation"],
				stages["Any"],
			),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectAllStages(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Subset(t, result, v.result, "wtf yo! %#v", v.result)
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
			id:     stages["Majority"].UUID,
			result: stages["Majority"],
		},
		"no_row_returned": {
			id:     "8",
			result: types.Stage{UUID: "8", Name: ""},
			err:    fmt.Errorf("sql: no rows in result set"),
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
			s:   stages["Gestation"],
			err: fmt.Errorf(uniqueKeyViolation, "stages_name_key"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertStage(context.Background(), v.s, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.NotEmpty(t, result.UUID)
		})
	}
}
func Test_UpdateStage(t *testing.T) {
	t.Skip()
	t.Parallel()

	set := map[string]struct {
		id  types.UUID
		s   types.Stage
		err error
	}{
		"happy_path": {
			id: "update me!",
			s:  types.Stage{Name: "renamed"},
		},
		"no_rows_affected": {
			id:  "missing",
			err: fmt.Errorf("stage was not updated: 'missing'"),
		},
		"duplicate_name_violation": {
			id:  "update me!",
			s:   stages["Gestation"],
			err: fmt.Errorf(uniqueKeyViolation, "stages_name_key"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.UpdateStage(context.Background(), v.id, v.s, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}
func Test_DeleteStage(t *testing.T) {
	t.Parallel()

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
			id:  stages["Colonization"].UUID,
			err: fmt.Errorf("foreign key violation: Key (uuid)=(1) is still referenced from table \"event_types\"., event_types."),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteStage(context.Background(), v.id, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

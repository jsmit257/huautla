package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var substrates = map[types.SubstrateType][]types.Substrate{}

func init() {
	for _, id := range []types.UUID{"0", "1", "2", "3", "update generation"} {
		if s, err := db.SelectSubstrate(context.Background(), id, "substrate_init"); err != nil {
			panic(err)
		} else {
			substrates[types.SubstrateType(s.Type)] = append(substrates[types.SubstrateType(s.Type)], s)
		}
	}
}

func Test_SelectAllSubstrates(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		result []types.Substrate
		err    error
	}{
		"happy_path": {
			result: append(
				[]types.Substrate{substrates[types.BulkType][0]},
				substrates[types.GrainType][0],
				substrates[types.LiquidType][0],
				substrates[types.PlatingType][0],
			),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectAllSubstrates(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Subset(t, result, v.result)
		})
	}
}

func Test_SelectSubstrate(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		result types.Substrate
		err    error
	}{
		"happy_path": {
			id:     substrates[types.GrainType][0].UUID,
			result: substrates[types.GrainType][0],
		},
		"no_rows_returned": {
			id:     "missing",
			result: types.Substrate{UUID: "missing"},
			err:    sql.ErrNoRows,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectSubstrate(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_InsertSubstrate(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s   types.Substrate
		err error
	}{
		"happy_path": {
			s: types.Substrate{Name: "Honey Solution", Type: types.BulkType, Vendor: vendors["localhost"]},
		},
		"unique_key_violation": {
			s:   substrates[types.BulkType][0],
			err: fmt.Errorf(uniqueKeyViolation, "substrates_name_vendor_uuid_key"),
		},
		"check_constraint_violation": {
			s:   types.Substrate{Name: "Maltodextrin", Type: "Stardust", Vendor: vendors["localhost"]},
			err: fmt.Errorf(checkConstraintViolation, "substrates", "substrates_type_check"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertSubstrate(context.Background(), v.s, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.NotEmpty(t, result.UUID)
		})
	}
}

func Test_UpdateSubstrate(t *testing.T) {
	t.Parallel()

	substrate, err := db.SelectSubstrate(context.Background(), "update me!", "Test_UpdateSubstrate")
	require.Nil(t, err)

	set := map[string]struct {
		s   types.Substrate
		err error
	}{
		"happy_path": {
			s: func() types.Substrate {
				result := substrate
				result.Name = "Updated"
				return result
			}(),
		},
		"unique_key_violation_name": {
			s: func() types.Substrate {
				result := substrate
				result.Name = "Millet"
				return result
			}(),
			err: fmt.Errorf(uniqueKeyViolation, "substrates_name_vendor_uuid_key"),
		},
		// "unique_key_violation_vendor": { // XXX: can't currently update vendpr
		// 	s: func() types.Substrate {
		// 		result := substrate
		// 		result.Vendor = types.Vendor{UUID: "missing"}
		// 		return result
		// 	}(),
		// 	err: fmt.Errorf("duplicate key violation"),
		// },
		// "check_constraint_violation": { // XXX: can't currently update type
		// 	s: func() types.Substrate {
		// 		result := substrate
		// 		result.Type = "Stardust"
		// 		return result
		// 	}(),
		// 	err: fmt.Errorf("check constraint"),
		// },
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.UpdateSubstrate(context.Background(), substrate.UUID, v.s, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

func Test_DeleteSubstrate(t *testing.T) {
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
			err: fmt.Errorf("substrate could not be deleted: 'missing'"),
		},
		"referential_violation": {
			id: substrates[types.GrainType][0].UUID,
			err: fmt.Errorf(foreignKeyViolation1toMany,
				"substrates",
				"substrate_ingredients_substrate_uuid_fkey",
				"substrate_ingredients"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteSubstrate(context.Background(), v.id, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

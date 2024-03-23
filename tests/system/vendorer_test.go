package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"

	"github.com/stretchr/testify/require"
)

var vendors = map[string]types.Vendor{}

func init() {
	for _, id := range []types.UUID{"localhost"} {
		if v, err := db.SelectVendor(context.Background(), id, "substrate_init"); err != nil {
			panic(err)
		} else {
			vendors[string(v.UUID)] = v
		}
	}
}

func Test_SelectAllVendors(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		result []types.Vendor
		err    error
	}{
		"happy_path": {
			result: []types.Vendor{vendors["localhost"]},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectAllVendors(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result[:len(v.result)])
		})
	}
}
func Test_SelectVendor(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		result types.Vendor
		err    error
	}{
		// "happy_path": { // already been done
		// 	id:     "localhost",
		// 	result: vendors["localhost"],
		// },
		"no_row_returned": {
			id:     "8",
			result: types.Vendor{UUID: "8", Name: ""},
			err:    fmt.Errorf("sql: no rows in result set"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectVendor(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}
func Test_InsertVendor(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		v      types.Vendor
		result types.Vendor
		err    error
	}{
		"happy_path": {
			v: types.Vendor{Name: "inserted vendor"},
		},
		"duplicate_name_violation": {
			v:   vendors["localhost"],
			err: fmt.Errorf(`pq: duplicate key value violates unique constraint "vendors_name_key"`),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertVendor(context.Background(), v.v, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.NotEmpty(t, result.UUID)
		})
	}
}
func Test_UpdateVendor(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id  types.UUID
		v   types.Vendor
		err error
	}{
		"happy_path": {
			id: "update me!",
			v:  types.Vendor{Name: "localhost"},
		},
		"duplicate_name_violation": {
			id:  "update me!",
			v:   vendors["localhost"],
			err: fmt.Errorf(uniqueKeyViolation, "vendors_name_key"),
		},
		"no_rows_affected": {
			id:  "missing",
			err: fmt.Errorf("vendor was not updated: 'missing'"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.UpdateVendor(context.Background(), v.id, v.v, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}
func Test_DeleteVendor(t *testing.T) {
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
			err: fmt.Errorf("vendor could not be deleted: 'missing'"),
		},
		"referential_violation": {
			id: "localhost",
			err: fmt.Errorf(foreignKeyViolation1toMany,
				"vendors",
				"substrates_vendor_uuid_fkey",
				"substrates"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteVendor(context.Background(), v.id, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

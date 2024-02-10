package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var vendor0 = types.Vendor{UUID: "0", Name: "127.0.0.1"}

func Test_SelectAllVendors(t *testing.T) {
	set := map[string]struct {
		result []types.Vendor
		err    error
	}{
		"happy_path": {
			result: []types.Vendor{vendor0},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectAllVendors(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}
func Test_SelectVendor(t *testing.T) {
	set := map[string]struct {
		id     types.UUID
		result types.Vendor
		err    error
	}{
		"happy_path": {
			id:     "0",
			result: vendor0,
		},
		"no_row_returned": {
			id:     "8",
			result: types.Vendor{UUID: "8", Name: ""},
			err:    fmt.Errorf("sql: no rows in result set"),
		},
		"query_fails": {
			id:     "8888888888888888888888888888888888888888888888888888888888888888",
			result: types.Vendor{UUID: "0", Name: ""},
			err:    noRows, // XXX
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
	set := map[string]struct {
		v      types.Vendor
		result types.Vendor
		err    error
	}{
		"happy_path": {
			v: types.Vendor{Name: "inserted vendor"},
		},
		"duplicate_name_violation": {
			v: vendor0,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertVendor(context.Background(), v.v, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}
func Test_UpdateVendor(t *testing.T) {
	set := map[string]struct {
		id  types.UUID
		v   types.Vendor
		err error
	}{
		"happy_path": {
			id: "0",
			v:  types.Vendor{Name: "localhost"},
		},
		"duplicate_name_violation": {
			id: "1",
			v:  vendor0,
		},
		"no_rows_affected": {
			id:  "foobar",
			err: fmt.Errorf("vendor table was not updated 'foobar'"),
		},
		"query_fails": {
			id: "12121212121212121212121212121212121212121",
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.UpdateVendor(context.Background(), v.id, v.v, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
func Test_DeleteVendor(t *testing.T) {
	set := map[string]struct {
		id  types.UUID
		err error
	}{
		"happy_path": {
			id: "-1",
		},
		"no_rows_affected": {
			id:  "foobar",
			err: fmt.Errorf("vendor table was not deleted 'foobar'"),
		},
		"query_fails": {
			id:  "01234567890123456789012345678901234567891",
			err: fmt.Errorf("some error"),
		},
		"referential_violation": {
			id:  "0",
			err: fmt.Errorf("referential constraint"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteVendor(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}

package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

func Test_SelectAllIngredients(t *testing.T) {
	set := map[string]struct {
		result []types.Ingredient
		err    error
	}{
		"happy_path": {
			result: []types.Ingredient{
				{UUID: "0", Name: "Vermiculite"},
				{UUID: "1", Name: "Maltodexterin"},
				{UUID: "2", Name: "Rye"},
				{UUID: "3", Name: "White Millet"},
				{UUID: "4", Name: "Popcorn"},
				{UUID: "5", Name: "Manure"},
				{UUID: "6", Name: "Coir"},
				{UUID: "7", Name: "Honey"},
				{UUID: "8", Name: "Agar"},
				{UUID: "9", Name: "Rice Flour"},
				{UUID: "10", Name: "White Milo"},
				{UUID: "11", Name: "Red Milo"},
				{UUID: "12", Name: "Red Millet"},
				{UUID: "13", Name: "Gypsum"},
				{UUID: "14", Name: "Calcium phosphate"},
				{UUID: "15", Name: "Diammonium phosphate"},
			},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectAllIngredients(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}
func Test_SelectIngredient(t *testing.T) {
	set := map[string]struct {
		id     types.UUID
		result types.Ingredient
		err    error
	}{
		"happy_path": {
			id:     "13",
			result: types.Ingredient{UUID: "13", Name: "Gypsum"},
		},
		"no_rows_returned": {
			id:     "foobar",
			result: types.Ingredient{UUID: "foobar", Name: ""},
			err:    fmt.Errorf("sql: no rows in result set"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectIngredient(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}
func Test_InsertIngredient(t *testing.T) {
	set := map[string]struct {
		i      types.Ingredient
		result []types.Ingredient
		err    error
	}{
		"happy_path": {
			i: types.Ingredient{Name: "bogus!"},
		},
		"no_rows_affected": {}, // ???
		"duplicate_name_violation": {
			i:      types.Ingredient{Name: "Coir"},
			result: []types.Ingredient{},
			err:    fmt.Errorf(""),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertIngredient(context.Background(), v.i, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}
func Test_UpdateIngredient(t *testing.T) {
	set := map[string]struct {
		id  types.UUID
		i   types.Ingredient
		err error
	}{
		"happy_path": {
			id: "0",
			i:  types.Ingredient{Name: "renamed"},
		},
		"no_rows_affected": {
			id:  "foobar",
			err: fmt.Errorf("ingredient was not updated: '0'"),
		},
		"duplicate_name_violation": {
			id: "0",
			i:  types.Ingredient{Name: "Honey"},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.UpdateIngredient(context.Background(), v.id, v.i, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
func Test_DeleteIngredient(t *testing.T) {
	set := map[string]struct {
		id  types.UUID
		err error
	}{
		"happy_path": {
			id: "-1",
		},
		"no_rows_affected": {
			id:  "foobar",
			err: fmt.Errorf("ingredient was not deleted: '0'"),
		},
		"query_fails": {
			id:  "01234567890123456789012345678901234567891",
			err: fmt.Errorf("ingredient was not deleted: '0'"),
		},
		"referential_violation": {
			id:  "2",
			err: fmt.Errorf("referential constraint"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteIngredient(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}

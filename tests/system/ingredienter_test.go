package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var ingredients = []types.Ingredient{
	{UUID: "0", Name: "Vermiculite"},
	{UUID: "1", Name: "Maltodextrin"},
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
}

func Test_SelectAllIngredients(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		result []types.Ingredient
		err    error
	}{
		"happy_path": {result: ingredients},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectAllIngredients(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Subset(t, result, v.result)
		})
	}
}

func Test_SelectIngredient(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		result types.Ingredient
		err    error
	}{
		"happy_path": {
			id:     ingredients[13].UUID,
			result: ingredients[13],
		},
		"no_rows_returned": {
			id:     "missing",
			result: types.Ingredient{UUID: "missing", Name: ""},
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
	t.Parallel()

	set := map[string]struct {
		i   types.Ingredient
		err error
	}{
		"happy_path": {
			i: types.Ingredient{Name: "bogus!"},
		},
		"duplicate_name_violation": {
			i:   types.Ingredient{Name: "Coir"},
			err: fmt.Errorf(uniqueKeyViolation, "ingredients_name_key"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertIngredient(context.Background(), v.i, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.NotEmpty(t, result.UUID)
		})
	}
}
func Test_UpdateIngredient(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id  types.UUID
		i   types.Ingredient
		err error
	}{
		"happy_path": {
			id: "update me!",
			i:  types.Ingredient{Name: "renamed"},
		},
		"no_rows_affected": {
			id:  "missing",
			err: fmt.Errorf("ingredient was not updated: 'missing'"),
		},
		"duplicate_name_violation": {
			id:  "update me!",
			i:   types.Ingredient{Name: "Honey"},
			err: fmt.Errorf(uniqueKeyViolation, "ingredients_name_key"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.UpdateIngredient(context.Background(), v.id, v.i, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}
func Test_DeleteIngredient(t *testing.T) {
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
			err: fmt.Errorf("ingredient could not be deleted: 'missing'"),
		},
		"referential_violation": {
			id: ingredients[12].UUID,
			err: fmt.Errorf(foreignKeyViolation1toMany,
				"ingredients",
				"substrate_ingredients_ingredient_uuid_fkey",
				"substrate_ingredients"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteIngredient(context.Background(), v.id, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

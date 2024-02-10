package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

func Test_GetAllIngredients(t *testing.T) {
	set := map[string]struct {
		s      types.Substrate
		result []types.Ingredient
		err    error
	}{
		"happy_path": {
			s: types.Substrate{UUID: "1"},
			result: []types.Ingredient{
				{UUID: "3", Name: "White Millet"},
				{UUID: "12", Name: "Red Millet"},
			},
		},
		"no_rows_found": {
			s:   types.Substrate{UUID: "missing"},
			err: fmt.Errorf("sql: no rows in result set"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.GetAllIngredients(context.Background(), &v.s, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, v.s.Ingredients)
		})
	}
}
func Test_AddIngredient(t *testing.T) {
	set := map[string]struct {
		s      types.Substrate
		i      types.Ingredient
		result []types.Ingredient
		err    error
	}{
		"happy_path": {
			s:      types.Substrate{UUID: "0"},
			i:      types.Ingredient{UUID: "9"},
			result: []types.Ingredient{{UUID: "9"}},
		},
		"duplicate_key_violation": {
			s:   types.Substrate{UUID: "0"},
			i:   types.Ingredient{UUID: "2"},
			err: fmt.Errorf("duplicate key violation"),
		},
		"no_rows_affected_ingredient": {
			s:   types.Substrate{UUID: "0"},
			i:   types.Ingredient{UUID: "-2"},
			err: fmt.Errorf("substrateingredient was not added"),
		},
		"no_rows_affected_substrate": {
			s:   types.Substrate{UUID: "-0"},
			i:   types.Ingredient{UUID: "2"},
			err: fmt.Errorf("substrateingredient was not added"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.AddIngredient(context.Background(), &v.s, v.i, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, v.s.Ingredients)
		})
	}
}
func Test_ChangeIngredient(t *testing.T) {
	set := map[string]struct {
		s          types.Substrate
		oldI, newI types.Ingredient
		result     []types.Ingredient
		err        error
	}{
		"happy_path": {
			s: types.Substrate{UUID: "1", Ingredients: []types.Ingredient{
				{UUID: "3"},
				{UUID: "12"},
			}},
			oldI: types.Ingredient{UUID: "3"},
			newI: types.Ingredient{UUID: "4"},
			result: []types.Ingredient{
				{UUID: "4"},
				{UUID: "12"},
			},
		},
		"no_rows_affected": {
			s: types.Substrate{UUID: "1", Ingredients: []types.Ingredient{
				{UUID: "3"},
				{UUID: "12"},
			}},
			oldI: types.Ingredient{UUID: "4"},
			newI: types.Ingredient{UUID: "4"},
			result: []types.Ingredient{
				{UUID: "3"},
				{UUID: "12"},
			},
			err: fmt.Errorf("substrateingredient was not changed"),
		},
		"unique_key_violation": {
			s: types.Substrate{UUID: "1", Ingredients: []types.Ingredient{
				{UUID: "3"},
				{UUID: "12"},
			}},
			oldI: types.Ingredient{UUID: "3"},
			newI: types.Ingredient{UUID: "12"},
			result: []types.Ingredient{
				{UUID: "3"},
				{UUID: "12"},
			},
			err: fmt.Errorf("unique key violation"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.ChangeIngredient(context.Background(), &v.s, v.oldI, v.newI, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, v.s.Ingredients)
		})
	}
}
func Test_RemoveIngredient(t *testing.T) {
	set := map[string]struct {
		s      types.Substrate
		i      types.Ingredient
		result []types.Ingredient
		err    error
	}{
		"happy_path": {
			s: types.Substrate{UUID: "0", Ingredients: []types.Ingredient{
				{UUID: "1"},
				{UUID: "2"},
				{UUID: "3"},
			}},
			i: types.Ingredient{UUID: "2"},
			result: []types.Ingredient{
				{UUID: "1"},
				{UUID: "3"},
			},
		},
		"no_rows_affected": {
			s:   types.Substrate{UUID: "0"},
			i:   types.Ingredient{UUID: "12"},
			err: fmt.Errorf("substrateingredient was not removed"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.RemoveIngredient(context.Background(), &v.s, v.i, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, v.s.Ingredients)
		})
	}
}

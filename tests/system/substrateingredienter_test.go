package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

func Test_GetAllIngredients(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		s      types.Substrate
		result []types.Ingredient
		err    error
	}{
		"happy_path": {
			s:      substrates[1],
			result: []types.Ingredient{ingredients[3], ingredients[12]},
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
	t.Parallel()

	substrate, err := db.SelectSubstrate(context.Background(), "add ingredient", "Test_AddIngredient")
	require.Nil(t, err)

	set := map[string]struct {
		s      types.Substrate
		i      types.Ingredient
		result []types.Ingredient
		err    error
	}{
		"happy_path": {
			s:      substrate,
			i:      ingredients[9],
			result: []types.Ingredient{ingredients[2], ingredients[9]},
		},
		"duplicate_key_violation": {
			s:   substrate,
			i:   ingredients[2],
			err: fmt.Errorf("duplicate key violation"),
		},
		"no_rows_affected_ingredient": {
			s:   substrate,
			i:   types.Ingredient{UUID: "missing"},
			err: fmt.Errorf("substrateingredient was not added"),
		},
		"no_rows_affected_substrate": {
			s:   types.Substrate{UUID: "missing"},
			i:   types.Ingredient{UUID: "3"},
			err: fmt.Errorf("substrateingredient was not added"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.AddIngredient(context.Background(), &v.s, v.i, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, v.s.Ingredients)
		})
	}
}
func Test_ChangeIngredient(t *testing.T) {
	t.Parallel()

	substrate, err := db.SelectSubstrate(context.Background(), "add ingredient", "Test_ChangeIngredient")
	require.Nil(t, err)

	finalstate := []types.Ingredient{ingredients[4], ingredients[12]}

	set := map[string]struct {
		oldI, newI types.Ingredient
		err        error
	}{
		"happy_path": { // order matters, so does synchronous execution
			oldI: ingredients[3],
			newI: ingredients[4],
		},
		"no_rows_affected_old_ingredient": {
			oldI: types.Ingredient{UUID: "missing"},
			newI: ingredients[4],
			err:  fmt.Errorf("substrateingredient was not changed"),
		},
		"unique_key_violation_ingredient": {
			oldI: ingredients[4],
			newI: ingredients[12],
			err:  fmt.Errorf("unique key violation"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			// t.Parallel()
			err := db.ChangeIngredient(context.Background(), &substrate, v.oldI, v.newI, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, finalstate, substrate.Ingredients)
		})
	}
}
func Test_RemoveIngredient(t *testing.T) {
	t.Parallel()

	substrate, err := db.SelectSubstrate(context.Background(), "remove ingredient", "Test_RemoveIngredient")
	require.Nil(t, err)

	result := []types.Ingredient{ingredients[0], ingredients[2]}

	set := map[string]struct {
		i   types.Ingredient
		err error
	}{
		"happy_path": { // happy path has to run first
			i: substrate.Ingredients[1],
		},
		"no_rows_affected_ingredient": {
			i:   types.Ingredient{UUID: "missing"},
			err: fmt.Errorf("substrateingredient was not removed"),
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.RemoveIngredient(context.Background(), &substrate, v.i, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, result, substrate.Ingredients)
		})
	}
}

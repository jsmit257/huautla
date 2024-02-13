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
			s: types.Substrate{UUID: "missing"},
		},
	}
	for k, v := range set {
		t.Run(k, func(t *testing.T) {
			err := db.GetAllIngredients(context.Background(), &v.s, types.CID(k))
			require.Equal(t, v.err, err)
			require.ElementsMatch(t, v.s.Ingredients, v.result)
		})
	}
}
func Test_AddIngredient(t *testing.T) {
	t.Parallel()

	substrate, err := db.SelectSubstrate(context.Background(), "add ingredient", "Test_AddIngredient")
	require.Nil(t, err)

	set := map[string]struct {
		i      types.Ingredient
		result []types.Ingredient
		err    error
	}{
		"happy_path": {
			i:      ingredients[9],
			result: append(substrate.Ingredients, ingredients[9]),
		},
		"duplicate_key_violation": {
			i:      ingredients[2],
			result: substrate.Ingredients[:],
			err:    fmt.Errorf(uniqueKeyViolation, "substrate_ingredients_substrate_uuid_ingredient_uuid_key"),
		},
		"no_rows_affected_ingredient": {
			i:      types.Ingredient{UUID: "missing"},
			result: substrate.Ingredients[:],
			err:    fmt.Errorf("substrateingredient was not added"),
		},
	}
	for k, v := range set {
		k, v, substrate := k, v, substrate
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.AddIngredient(context.Background(), &substrate, v.i, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.ElementsMatch(t, v.result, substrate.Ingredients)
		})
	}
}
func Test_ChangeIngredient(t *testing.T) {
	t.Parallel()

	substrate, err := db.SelectSubstrate(context.Background(), "change ingredient", "Test_ChangeIngredient")
	require.Nil(t, err)

	set := map[string]struct {
		oldI, newI types.Ingredient
		result     []types.Ingredient
		err        error
	}{
		"no_rows_affected_old_ingredient": {
			oldI:   types.Ingredient{UUID: "missing"},
			newI:   ingredients[4],
			result: substrate.Ingredients,
			err:    fmt.Errorf("substrateingredient was not changed"),
		},
		"unique_key_violation_ingredient": {
			oldI:   ingredients[3],
			newI:   ingredients[12],
			result: substrate.Ingredients,
			err:    fmt.Errorf(uniqueKeyViolation, "substrate_ingredients_substrate_uuid_ingredient_uuid_key"),
		},
		"no_rows_affected_new_ingredient": {
			oldI:   ingredients[3],
			newI:   types.Ingredient{UUID: "missing"},
			result: substrate.Ingredients,
			err: fmt.Errorf(
				foreignKeyViolation1to1,
				"substrate_ingredients",
				"substrate_ingredients_ingredient_uuid_fkey"),
		},
		// "happy_path": { // this case borks everything, but only sometimes, and only sometimes in the same way
		// 	oldI:   ingredients[3],
		// 	newI:   ingredients[4],
		// 	result: []types.Ingredient{ingredients[4], ingredients[12]},
		// },
	}
	for k, v := range set {
		k, v, substrate := k, v, substrate
		t.Run(k, func(t *testing.T) {
			// t.Parallel() // don't get the shenanigans here but whatever, for now
			err := db.ChangeIngredient(context.Background(), &substrate, v.oldI, v.newI, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.ElementsMatch(t, v.result, substrate.Ingredients)
		})
	}
	// moved the happy path here for now, something is really wonky here, and why does changing strains affect this branch?
	err = db.ChangeIngredient(context.Background(), &substrate, ingredients[3], ingredients[4], "Test_ChangeIngredient")
	require.Nil(t, err)
	require.ElementsMatch(t, []types.Ingredient{ingredients[4], ingredients[12]}, substrate.Ingredients)
}
func Test_RemoveIngredient(t *testing.T) {
	t.Parallel()

	substrate, err := db.SelectSubstrate(context.Background(), "remove ingredient", "Test_RemoveIngredient")
	require.Nil(t, err)

	set := map[string]struct {
		i      types.Ingredient
		result []types.Ingredient
		err    error
	}{
		"happy_path": { // happy path has to run first
			i:      substrate.Ingredients[1],
			result: []types.Ingredient{ingredients[12], ingredients[14]},
		},
		"no_rows_affected_ingredient": {
			i:      types.Ingredient{UUID: "missing"},
			result: substrate.Ingredients,
			err:    fmt.Errorf("substrateingredient was not removed"),
		},
	}
	for k, v := range set {
		k, v, substrate := k, v, substrate
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.RemoveIngredient(context.Background(), &substrate, v.i, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.ElementsMatch(t, v.result, substrate.Ingredients)
		})
	}
}

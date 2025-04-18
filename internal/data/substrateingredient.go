package data

import (
	"context"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) GetAllIngredients(ctx context.Context, s *types.Substrate, cid types.CID) error {
	var err error
	deferred, l := initAccessFuncs("GetAllIngredients", db.logger, s.UUID, cid)
	defer deferred(&err, l)

	rows, err := db.query.QueryContext(ctx, psqls["substrate-ingredient"]["all"], s.UUID)
	if err != nil {
		return err
	}

	ing := make([]types.Ingredient, 0, 100)
	for rows.Next() {
		row := types.Ingredient{}
		if err = rows.Scan(
			&row.UUID,
			&row.Name,
		); err != nil {
			return err
		}
		ing = append(ing, row)
	}

	s.Ingredients = ing

	return nil
}

func (db *Conn) AddIngredient(ctx context.Context, s *types.Substrate, i types.Ingredient, cid types.CID) error {
	var err error
	deferred, l := initAccessFuncs("AddIngredient", db.logger, "nil", cid)
	defer deferred(&err, l)

	var rows int64
	result, err := db.query.ExecContext(ctx, psqls["substrate-ingredient"]["add"], db.generateUUID(), s.UUID, i.UUID)
	if err != nil {
		if isPrimaryKeyViolation(err) {
			return db.AddIngredient(ctx, s, i, cid) // FIXME: infinite loop?
		}
		return err
	} else if rows, err = result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		err = fmt.Errorf("substrateingredient was not added")
	} else {
		s.Ingredients = append(s.Ingredients, i)
	}

	return err
}

func (db *Conn) ChangeIngredient(ctx context.Context, s *types.Substrate, oldI, newI types.Ingredient, cid types.CID) error {
	var err error
	deferred, l := initAccessFuncs("ChangeIngredient", db.logger, "nil", cid)
	defer deferred(&err, l)

	var rows int64
	result, err := db.query.ExecContext(ctx, psqls["substrate-ingredient"]["change"], newI.UUID, s.UUID, oldI.UUID)
	if err != nil {
		return err
	} else if rows, err = result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("substrateingredient was not changed")
	}
	i, j := 0, len(s.Ingredients)
	for i < j && s.Ingredients[i].UUID != oldI.UUID {
		i++
	}
	s.Ingredients[i] = newI

	return err
}

func (db *Conn) RemoveIngredient(ctx context.Context, s *types.Substrate, i types.Ingredient, cid types.CID) error {
	var err error
	deferred, l := initAccessFuncs("RemoveIngredient", db.logger, s.UUID, cid)
	defer deferred(&err, l)

	var rows int64
	result, err := db.query.ExecContext(ctx, psqls["substrate-ingredient"]["remove"], s.UUID, i.UUID)
	if err != nil {
		return err
	} else if rows, err = result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		err = fmt.Errorf("substrateingredient was not removed")
		return err
	}

	ndx, j := 0, len(s.Ingredients)
	for ndx < j && s.Ingredients[ndx].UUID != i.UUID {
		ndx++
	}
	s.Ingredients = append(s.Ingredients[:ndx], s.Ingredients[ndx+1:]...)

	return err
}

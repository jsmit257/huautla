package data

import (
	"context"
	"database/sql"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) GetAllIngredients(ctx context.Context, s *types.Substrate, cid types.CID) error {
	var err error

	deferred, start, l := initVendorFuncs("GetAllIngredients", db.logger, err, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	s.Ingredients = make([]types.Ingredient, 0, 100)

	rows, err = db.query.QueryContext(ctx, db.sql["substrate-ingredient"]["all-ingredients"], s.UUID)
	if err != nil {
		return err
	}

	for rows.Next() {
		row := types.Ingredient{}
		rows.Scan(
			&row.UUID,
			&row.Name)
		s.Ingredients = append(s.Ingredients, row)
	}

	return err
}

func (db *Conn) AddIngredient(ctx context.Context, s types.Substrate, i types.Ingredient, cid types.CID) error {
	return nil
}

func (db *Conn) ChangeIngredient(ctx context.Context, s types.Substrate, oldI, newI types.Ingredient, cid types.CID) error {
	return nil
}

func (db *Conn) RemoveIngredient(ctx context.Context, sid types.UUID, i types.Ingredient, cid types.CID) error {
	return nil
}

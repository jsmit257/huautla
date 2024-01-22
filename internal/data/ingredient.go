package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllIngredients(ctx context.Context, cid types.CID) ([]types.Ingredient, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectAllIngredients", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.Ingredient, 0, 100)

	rows, err = db.query.QueryContext(ctx, db.sql["ingredient"]["select-all"])
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := types.Ingredient{}
		rows.Scan(&row.UUID, &row.Name)
		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectIngredient(ctx context.Context, id types.UUID, cid types.CID) (types.Ingredient, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectIngredient", db.logger, id, cid)
	defer deferred(start, err, l)

	result := types.Ingredient{UUID: id}
	err = db.
		QueryRowContext(ctx, db.sql["ingredient"]["select"], id).
		Scan(&result.Name)

	return result, err
}

func (db *Conn) InsertIngredient(ctx context.Context, i types.Ingredient, cid types.CID) (types.Ingredient, error) {
	var err error

	i.UUID = types.UUID(db.generateUUID().String())

	deferred, start, l := initVendorFuncs("InsertIngredient", db.logger, i.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["ingredient"]["insert"], i.UUID, i.Name)
	if err != nil {
		// FIXME: choose what to do based on the tupe of error
		duplicatePrimaryKeyErr := false
		if duplicatePrimaryKeyErr {
			return db.InsertIngredient(ctx, i, cid) // FIXME: infinite loop?
		}
		return i, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return i, err
	} else if rows != 1 {
		return i, fmt.Errorf("ingredient was not added")
	}

	return i, err
}

func (db *Conn) UpdateIngredient(ctx context.Context, id types.UUID, i types.Ingredient, cid types.CID) error {
	var err error

	deferred, start, l := initVendorFuncs("UpdateIngredient", db.logger, id, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["ingredient"]["update"], i.Name, id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("ingredient was not updated: '%s'", id)
	}
	return nil
}

func (db *Conn) DeleteIngredient(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteIngredient", "ingredient", db.logger)
}

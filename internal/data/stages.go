package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllStages(ctx context.Context, cid types.CID) ([]types.Stage, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectAllStages", db.logger, err, types.UUID("nil"), cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.Stage, 0, 100)

	rows, err = db.query.QueryContext(ctx, db.sql["stage"]["select-all"])
	if err != nil {
		return nil, err
	} else if rows == nil {
		return nil, fmt.Errorf("no result returned from SelectAllStages")
	}

	for rows.Next() {
		row := types.Stage{}
		rows.Scan(&row.UUID, &row.Name)
		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectStage(ctx context.Context, id types.UUID, cid types.CID) (types.Stage, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectStage", db.logger, err, id, cid)
	defer deferred(start, err, l)

	result := types.Stage{UUID: id}
	err = db.
		QueryRowContext(ctx, db.sql["stage"]["select"], id).
		Scan(&result.Name)

	return result, err
}

func (db *Conn) InsertStage(ctx context.Context, v types.Stage, cid types.CID) (types.Stage, error) {
	var err error

	v.UUID = types.UUID(db.generateUUID().String())

	deferred, start, l := initVendorFuncs("InsertStage", db.logger, err, v.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["stage"]["insert"], v.UUID, v.Name)
	if err != nil {
		// FIXME: choose what to do based on the tupe of error
		duplicatePrimaryKeyErr := false
		if duplicatePrimaryKeyErr {
			return db.InsertStage(ctx, v, cid) // FIXME: infinite loop?
		}
		return v, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return v, err
	} else if rows != 1 {
		return v, fmt.Errorf("stage was not added")
	}

	return v, err
}

func (db *Conn) UpdateStage(ctx context.Context, id types.UUID, v types.Stage, cid types.CID) error {
	var err error

	deferred, start, l := initVendorFuncs("UpdateStage", db.logger, err, id, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["stage"]["update"], v.Name, id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("stage was not updated: '%s'", id)
	}
	return nil
}

func (db *Conn) DeleteStage(ctx context.Context, id types.UUID, cid types.CID) error {
	var err error

	deferred, start, l := initVendorFuncs("DeleteStage", db.logger, err, id, cid)
	defer deferred(start, err, l)

	var result sql.Result

	l.Info("starting work")

	result, err = db.ExecContext(ctx, db.sql["stage"]["delete"], id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		// this won't be reported in the WithError log in `defer ...`, b/c it's operator error
		return fmt.Errorf("stage could not be deleted: '%s'", id)
	}

	return err
}

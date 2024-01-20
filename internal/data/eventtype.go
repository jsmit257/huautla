package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) SelectAllEventTypes(ctx context.Context, cid types.CID) ([]types.EventType, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectAllEventTypes", db.logger, err, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	result := make([]types.EventType, 0, 100)

	rows, err = db.query.QueryContext(ctx, db.sql["eventtype"]["select-all"])
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row := types.EventType{}
		rows.Scan(
			&row.UUID,
			&row.Name,
			&row.Stage.UUID,
			&row.Stage.Name)
		result = append(result, row)
	}

	return result, err
}

func (db *Conn) SelectEventType(ctx context.Context, id types.UUID, cid types.CID) (types.EventType, error) {
	var err error

	deferred, start, l := initVendorFuncs("SelectEventType", db.logger, err, id, cid)
	defer deferred(start, err, l)

	result := types.EventType{UUID: id}
	err = db.
		QueryRowContext(ctx, db.sql["eventtype"]["select"], id).
		Scan(
			&result.Name,
			&result.Stage.UUID,
			&result.Stage.Name)

	return result, err
}

func (db *Conn) InsertEventType(ctx context.Context, e types.EventType, cid types.CID) (types.EventType, error) {
	var err error

	e.UUID = types.UUID(db.generateUUID().String())

	deferred, start, l := initVendorFuncs("InsertEventType", db.logger, err, e.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["eventtype"]["insert"], e.UUID, e.Name, e.Stage.UUID)
	if err != nil {
		// FIXME: choose what to do based on the tupe of error
		duplicatePrimaryKeyErr := false
		if duplicatePrimaryKeyErr {
			return db.InsertEventType(ctx, e, cid) // FIXME: infinite loop?
		}
		return e, err
	} else if rows, err := result.RowsAffected(); err != nil {
		return e, err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return e, fmt.Errorf("eventtype was not added")
	}

	return e, err
}

func (db *Conn) UpdateEventType(ctx context.Context, id types.UUID, s types.EventType, cid types.CID) error {
	var err error

	deferred, start, l := initVendorFuncs("UpdateEventType", db.logger, err, id, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["eventtype"]["update"], s.Name, id)
	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("eventtype was not updated: '%s'", id)
	}
	return nil
}

func (db *Conn) DeleteEventType(ctx context.Context, id types.UUID, cid types.CID) error {
	return db.deleteByUUID(ctx, id, cid, "DeleteEventType", "eventtype", db.logger)
}

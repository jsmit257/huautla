package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) KnownAttributeNames(ctx context.Context, cid types.CID) ([]string, error) {
	var err error
	var result = []string{}
	var s string

	deferred, start, l := initVendorFuncs("KnownAttributeNames", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	rows, err = db.query.QueryContext(ctx, db.sql["strainattribute"]["all-attributes"])
	if err != nil {
		return result, err
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&s)
		result = append(result, s)
	}

	return result, err
}

func (db *Conn) GetAllAttributes(ctx context.Context, s *types.Strain, cid types.CID) error {
	var err error

	deferred, start, l := initVendorFuncs("GetAllAttributes", db.logger, "nil", cid)
	defer deferred(start, err, l)

	var rows *sql.Rows

	s.Attributes = make([]types.StrainAttribute, 0, 100)

	rows, err = db.query.QueryContext(ctx, db.sql["strainattribute"]["all"])
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		row := types.StrainAttribute{}
		rows.Scan(
			&row.UUID,
			&row.Name,
			&row.Value)
		s.Attributes = append(s.Attributes, row)
	}

	return err
}

func (db *Conn) AddAttribute(ctx context.Context, s *types.Strain, n, v string, cid types.CID) error {
	var err error

	id := types.UUID(db.generateUUID().String())

	deferred, start, l := initVendorFuncs("AddAttribute", db.logger, s.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["strainattribute"]["add"], id, n, v, s.UUID)

	if err != nil {
		if isUniqueViolation(err) {
			return db.AddAttribute(ctx, s, n, v, cid) // FIXME: infinite loop?
		}
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("attribute was not added")
	}

	s.Attributes = append(s.Attributes, types.StrainAttribute{id, n, v})

	return err
}

func (db *Conn) ChangeAttribute(ctx context.Context, s *types.Strain, id types.UUID, n, v string, cid types.CID) error {
	var err error

	deferred, start, l := initVendorFuncs("ChangeAttribute", db.logger, s.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["strainattribute"]["change"], v, n, s.UUID)

	if err != nil {
		if isUniqueViolation(err) {
			return db.ChangeAttribute(ctx, s, id, n, v, cid) // FIXME: infinite loop?
		}
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("attribute was not changed")
	}

	i, j := 0, len(s.Attributes)
	for i < j && s.Attributes[i].UUID != id {
		i++
	}
	s.Attributes[i].Value = v

	return nil
}

func (db *Conn) RemoveAttribute(ctx context.Context, s *types.Strain, id types.UUID, cid types.CID) error {
	var err error

	deferred, start, l := initVendorFuncs("RemoveAttribute", db.logger, s.UUID, cid)
	defer deferred(start, err, l)

	result, err := db.ExecContext(ctx, db.sql["strainattribute"]["remove"], id)

	if err != nil {
		return err
	} else if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows != 1 { // most likely cause is a bad vendor.uuid
		return fmt.Errorf("attribute was not removed")
	}

	i, j := 0, len(s.Attributes)
	for i < j && s.Attributes[i].UUID != id {
		i++
	}
	s.Attributes = append(s.Attributes[:i], s.Attributes[i+1:]...)

	return nil
}

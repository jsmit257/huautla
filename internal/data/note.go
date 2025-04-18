package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jsmit257/huautla/types"
)

func (db *Conn) GetNotes(ctx context.Context, id types.UUID, cid types.CID) ([]types.Note, error) {
	var err error

	deferred, l := initAccessFuncs("GetNotes", db.logger, id, cid)
	defer deferred(&err, l)

	var rows *sql.Rows
	result := []types.Note{}

	rows, err = db.query.QueryContext(ctx, psqls["note"]["get"], id)
	if err != nil {
		return result, err
	}

	defer rows.Close()

	note := types.Note{}

	for rows.Next() {
		if err = rows.Scan(
			&note.UUID,
			&note.Note,
			&note.MTime,
			&note.CTime,
		); err != nil {
			return result, err
		}
		result = append(result, note)
	}

	return result, nil
}

func (db *Conn) AddNote(ctx context.Context, oID types.UUID, notes []types.Note, n types.Note, cid types.CID) ([]types.Note, error) {
	var err error
	deferred, l := initAccessFuncs("AddNote", db.logger, oID, cid)
	defer deferred(&err, l)

	n.UUID = types.UUID(db.generateUUID().String())
	n.MTime = time.Now().UTC()
	n.CTime = n.MTime

	var rows int64
	result, err := db.ExecContext(ctx, psqls["note"]["add"],
		n.UUID,
		n.Note,
		oID,
		n.CTime,
	)
	if err != nil {
		if isPrimaryKeyViolation(err) {
			return db.AddNote(ctx, oID, notes, n, cid)
		}
		return notes, err
	} else if rows, err = result.RowsAffected(); err != nil {
		return notes, err
	} else if rows != 1 {
		return notes, fmt.Errorf("note was not added")
	}

	return append([]types.Note{n}, notes...), err
}

func (db *Conn) ChangeNote(ctx context.Context, notes []types.Note, n types.Note, cid types.CID) ([]types.Note, error) {
	var err error
	deferred, l := initAccessFuncs("ChangeNote", db.logger, n.UUID, cid)
	defer deferred(&err, l)

	n.MTime = time.Now().UTC()

	var rows int64
	result, err := db.ExecContext(ctx, psqls["note"]["change"],
		n.Note,
		n.MTime,
		n.UUID,
	)
	if err != nil {
		return notes, err
	} else if rows, err = result.RowsAffected(); err != nil {
		return notes, err
	} else if rows != 1 {
		return notes, fmt.Errorf("note was not changed")
	}

	i, j := 0, len(notes)
	for i < j && notes[i].UUID != n.UUID {
		i++
	}

	return append(append([]types.Note{n}, notes[:i]...), notes[i+1:]...), nil
}

func (db *Conn) RemoveNote(ctx context.Context, notes []types.Note, id types.UUID, cid types.CID) ([]types.Note, error) {
	var err error
	deferred, l := initAccessFuncs("RemoveNote", db.logger, id, cid)
	defer deferred(&err, l)

	var rows int64
	result, err := db.ExecContext(ctx, psqls["note"]["remove"], id)
	if err != nil {
		return notes, err
	} else if rows, err = result.RowsAffected(); err != nil {
		return notes, err
	} else if rows != 1 {
		err = fmt.Errorf("note could not be removed")
		return notes, err
	}

	i, j := 0, len(notes)
	for i < j && notes[i].UUID != id {
		i++
	}

	return append(notes[:i], notes[i+1:]...), nil
}

func (db *Conn) notesReport(ctx context.Context, id types.UUID, cid types.CID, p *rpttree) ([]types.Entity, error) {
	notes, err := db.GetNotes(ctx, id, cid)
	if err != nil {
		return nil, err
	} else if len(notes) == 0 {
		return nil, nil
	}

	result := make([]types.Entity, len(notes))
	for i, n := range notes {
		rpt, err := db.newRpt(ctx, n, cid, p)
		if err != nil {
			return nil, err
		}
		result[i] = rpt.Data()
	}

	return result, nil
}

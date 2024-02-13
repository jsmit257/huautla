package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var eventtypes = []types.EventType{
	{UUID: "0", Name: "Condensation", Severity: "Warn", Stage: stages[3]},
	{UUID: "1", Name: "Fruiting", Severity: "Info", Stage: stages[1]},
	{UUID: "2", Name: "Crash", Severity: "Error", Stage: stages[1]},
	{UUID: "3", Name: "Sunset", Severity: "RIP", Stage: stages[2]},
}

func Test_SelectAllEventTypes(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		result []types.EventType
		err    error
	}{
		"happy_path": {result: eventtypes},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectAllEventTypes(context.Background(), types.CID(k))
			require.Equal(t, v.err, err)
			require.Subset(t, result, v.result)
		})
	}
}

func Test_SelectEventType(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		result types.EventType
		err    error
	}{
		"happy_path": {
			id:     "0",
			result: eventtypes[0],
		},
		"no_rows_returned": {
			id:     "missing",
			result: types.EventType{UUID: "missing"},
			err:    noRows,
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectEventType(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}

func Test_InsertEventType(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		e   types.EventType
		err error
	}{
		"happy_path": {
			e: types.EventType{Name: "bogus", Severity: "Info", Stage: stages[1]},
		},
		"no_rows_affected_typecheck": {
			e:   types.EventType{Name: "bogus", Stage: stages[0]},
			err: fmt.Errorf(checkConstraintViolation, "event_types", "event_types_severity_check"),
		},
		"no_rows_affected_stage": {
			e:   types.EventType{Name: "bogus", Severity: "Info", Stage: types.Stage{UUID: "missing"}},
			err: fmt.Errorf("eventtype was not added"),
		},
		"unique_key_violation": {
			e:   types.EventType{Name: "Vacation", Stage: stages[1]},
			err: fmt.Errorf(checkConstraintViolation, "event_types", "event_types_severity_check"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.InsertEventType(context.Background(), v.e, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.NotEmpty(t, result.UUID)
		})
	}
}

func Test_UpdateEventType(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id  types.UUID
		e   types.EventType
		err error
	}{
		"happy_path": {
			id: "update me!",
			e:  types.EventType{Name: "renamed"},
		},
		"no_rows_affected": {
			id:  "missing",
			err: fmt.Errorf("eventtype was not updated: 'missing'"),
		},
		// "no_rows_affected_typecheck": { // currently don't update severity
		// 	id:  "update me!",
		// 	e:   types.EventType{Name: "bogus", Stage: stages[0]},
		// 	err: fmt.Errorf(checkConstraintViolation, "event_types", "event_types_severity_check"),
		// },
		"unique_key_violation": { // currently don't update stage_uuid
			id:  "update me!",
			e:   types.EventType{Name: "Fruiting", Stage: stages[1]},
			err: fmt.Errorf(uniqueKeyViolation, "event_types_name_stage_uuid_key"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.UpdateEventType(context.Background(), v.id, v.e, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

func Test_DeleteEventType(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id  types.UUID
		err error
	}{
		"happy_path": {
			id: "delete me!",
		},
		"no_rows_affected": {
			id:  "foobar",
			err: fmt.Errorf("eventtype could not be deleted: 'foobar'"),
		},
		"referential_violation": {
			id: eventtypes[1].UUID,
			err: fmt.Errorf(foreignKeyViolation1toMany,
				"event_types",
				"events_eventtype_uuid_fkey",
				"events"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteEventType(context.Background(), v.id, types.CID(k))
			equalErrorMessages(t, v.err, err)
		})
	}
}

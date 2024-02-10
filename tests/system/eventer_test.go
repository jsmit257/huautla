package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var events = []types.Event{
	{UUID: "0", Temperature: 2, Humidity: 1, MTime: epoch, CTime: epoch, EventType: eventtypes[1]},
	{UUID: "1", Temperature: 0, Humidity: 1, MTime: epoch, CTime: epoch, EventType: eventtypes[0]},
	{UUID: "2", Temperature: 0, Humidity: 8, MTime: epoch, CTime: epoch, EventType: eventtypes[0]},
}

func Test_GetLifecycleEvents(t *testing.T) {
	gle := map[string]struct {
		lc     types.Lifecycle
		result []types.Event
		err    error
	}{
		"happy_path": {
			lc:     types.Lifecycle{UUID: "0"},
			result: []types.Event{events[0], events[2]},
		},
		"no_rows_returned": {
			lc:  types.Lifecycle{UUID: "foobar"},
			err: fmt.Errorf("sql: no rows in result set"),
		},
	}
	for k, v := range gle {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			v.lc.Events = []types.Event{}
			err := db.GetLifecycleEvents(context.Background(), &v.lc, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, v.lc.Events[0:len(v.result)])
		})
	}
}
func Test_SelectByEventType(t *testing.T) {
	set := map[string]struct {
		e      types.EventType
		result []types.Event
		err    error
	}{
		"happy_path": {
			e:      eventtypes[0],
			result: []types.Event{events[1], events[2]},
		},
		"no_rows_returned": {
			e:   types.EventType{UUID: "missing"},
			err: fmt.Errorf("sql: no rows in result set"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			result, err := db.SelectByEventType(context.Background(), v.e, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result[0:len(v.result)])
		})
	}
}
func Test_SelectEvent(t *testing.T) {
	set := map[string]struct {
		id     types.UUID
		result types.Event
		err    error
	}{
		"happy_path": {
			id:     events[0].UUID,
			result: events[0],
		},
		"no_rows_returned": {
			id:  "foobar",
			err: fmt.Errorf("sql: no rows in result set"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			result, err := db.SelectEvent(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, result)
		})
	}
}
func Test_AddEvent(t *testing.T) {
	t.Parallel()

	_, err := db.SelectLifecycle(context.Background(), "add event", types.CID("Test_AddEvent"))
	require.Nil(t, err)

	set := map[string]struct {
		lc  types.Lifecycle
		e   types.Event
		err error
	}{
		"happy_path": {
			lc: types.Lifecycle{},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			err := db.AddEvent(context.Background(), &v.lc, v.e, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
func Test_ChangeEvent(t *testing.T) {
	t.Parallel()

	_, err := db.SelectLifecycle(context.Background(), "change event", types.CID("Test_ChangeEvent"))
	require.Nil(t, err)

	set := map[string]struct {
		lc  types.Lifecycle
		e   types.Event
		err error
	}{
		"happy_path": {
			e: types.Event{},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			err := db.ChangeEvent(context.Background(), &v.lc, v.e, types.CID(k))
			require.Equal(t, v.err, err)
		})
	}
}
func Test_RemoveEvent(t *testing.T) {
	t.Parallel()

	lc, err := db.SelectLifecycle(context.Background(), "remove event", types.CID("Test_RemoveEvent"))
	require.Nil(t, err)

	set := map[string]struct {
		id     types.UUID
		result []types.Event
		err    error
	}{ // NB: the order in which the test cases are declared matters
		"no_rows_affected": {
			id:     "missing event",
			result: lc.Events,
			err:    fmt.Errorf("event could not be removed"),
		},
		"happy_path": {
			id:     "remove event 2",
			result: []types.Event{lc.Events[0], lc.Events[2]},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			// t.Parallel() // XXX: don't do this, the data doesn't support it
			err := db.RemoveEvent(context.Background(), &lc, v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, lc.Events)
		})
	}
}

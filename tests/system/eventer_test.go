package test

// yeah, all these things could be tested in the lifecycler tests,
// but it's easy to separate them this way, and ignore these details
// in livecycler_test

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
	// t.Skip()
	t.Parallel()

	set := map[string]struct {
		lc     types.Lifecycle
		result []types.Event
		err    error
	}{
		"happy_path": {
			lc:     types.Lifecycle{UUID: "0"},
			result: []types.Event{events[0], events[2]},
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			var actual types.Event
			err := db.GetLifecycleEvents(context.Background(), &v.lc, types.CID(k))
			require.Equal(t, v.err, err)
			for i, j := 0, len(v.result); i < j; i++ {
				event := v.result[i]
				actual, err = findEvent(v.lc.Events, v.result[i].UUID)
				require.Nil(t, err)
				require.Equal(t, event.Temperature, actual.Temperature)
				require.Equal(t, event.Humidity, actual.Humidity)
				// require.Truef(t, event.MTime.String() == actual.MTime.String(), "expected\n'%#v'\nactual\n'%#v'", event.MTime.String(), actual.MTime.String()) // get fucked!
				// require.Truef(t, event.CTime == actual.CTime, "expected\n'%#q'\nactual\n'%#q'", event.CTime, actual.CTime)
				require.Equal(t, event.EventType, actual.EventType)
			}
		})
	}
}

func Test_SelectByEventType(t *testing.T) {
	t.Parallel()

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
			e:      types.EventType{UUID: "missing"},
			result: []types.Event{},
			// err: fmt.Errorf("sql: no rows in result set"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			var actual types.Event
			result, err := db.SelectByEventType(context.Background(), v.e, types.CID(k))
			require.Equal(t, v.err, err)
			for i, j := 0, len(v.result); i < j; i++ {
				event := v.result[i]
				actual, err = findEvent(result, v.result[i].UUID)
				require.Nil(t, err)
				require.Equal(t, event.Temperature, actual.Temperature)
				require.Equal(t, event.Humidity, actual.Humidity)
				// require.Equalf(t, event.MTime.UnixMilli(), actual.MTime.UnixMilli(), "expected\n'%#v'\nactual\n'%#v'", event.MTime.String(), actual.MTime.String()) // get fucked!
				// require.Truef(t, event.CTime == actual.CTime, "expected\n'%#q'\nactual\n'%#q'", event.CTime, actual.CTime)
				require.Equal(t, event.EventType, actual.EventType)
			}
		})
	}
}

func Test_SelectEvent(t *testing.T) {
	t.Parallel()

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
			id:     "missing",
			result: types.Event{UUID: "missing"},
			err:    fmt.Errorf("sql: no rows in result set"),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			result, err := db.SelectEvent(context.Background(), v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result.Temperature, result.Temperature)
			require.Equal(t, v.result.Humidity, result.Humidity)
			// require.Truef(t, v.result.MTime.String() == result.MTime.String(), "expected\n'%#v'\nactual\n'%#v'", v.result.MTime.String(), result.MTime.String()) // get fucked!
			// require.Truef(t, v.result.CTime == result.CTime, "expected\n'%#q'\nactual\n'%#q'", v.result.CTime, result.CTime)
			require.Equal(t, v.result.EventType, result.EventType)
		})
	}
}

func Test_AddEvent(t *testing.T) {
	t.Parallel()

	lc, err := db.SelectLifecycle(context.Background(), "add event", types.CID("Test_AddEvent"))
	require.Nil(t, err)

	set := map[string]struct {
		e     types.Event
		count int
		err   error
	}{
		"happy_path": { // happy path has to run first, synchronously
			e:     types.Event{EventType: eventtypes[1]},
			count: 2,
		},
		"no_rows_affected_eventtype": {
			e:     types.Event{EventType: types.EventType{UUID: "missing"}},
			count: 1,
			err:   fmt.Errorf("event was not added"),
		},
	}
	for k, v := range set {
		k, v, lc := k, v, lc
		t.Run(k, func(t *testing.T) {
			// t.Parallel()
			err := db.AddEvent(context.Background(), &lc, v.e, types.CID(k))
			t.Logf("actual: %#v", lc.Events)
			require.Equal(t, v.err, err)
			require.Equalf(t, v.count, len(lc.Events), "actual: %#v", lc.Events)
		})
	}
}

func Test_ChangeEvent(t *testing.T) {
	t.Parallel()

	lc, err := db.SelectLifecycle(context.Background(), "change event", types.CID("Test_ChangeEvent"))
	require.Nil(t, err)

	set := map[string]struct {
		e      types.Event
		result types.EventType
		err    error
	}{
		"happy_path": { // happy path needs to run first, synchronously
			e:      types.Event{UUID: "change event", EventType: eventtypes[1]},
			result: eventtypes[1],
		},
		"no_rows_affected": { // dunno how this would happen, but whatever
			e:      types.Event{UUID: "missing", EventType: eventtypes[0]},
			result: eventtypes[0],
			err:    fmt.Errorf("event was not changed"),
		},
		"no_rows_affected_eventtype": {
			e:      types.Event{EventType: types.EventType{UUID: "missing"}},
			result: eventtypes[0],
			err:    fmt.Errorf("event was not changed"),
		},
	}
	for k, v := range set {
		k, v, lc := k, v, lc
		lc.Events = append([]types.Event{}, lc.Events...)
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.ChangeEvent(context.Background(), &lc, v.e, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, lc.Events[0].EventType)
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
	}{ // NB: the order in which the test cases are declared matters: errors should come first
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
		k, v, lc := k, v, lc
		t.Run(k, func(t *testing.T) {
			t.Parallel() // XXX: don't do this, the data doesn't support it
			err := db.RemoveEvent(context.Background(), &lc, v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, lc.Events)
		})
	}
}

func findEvent(events []types.Event, id types.UUID) (types.Event, error) {
	for i, j := 0, len(events); i < j; i++ {
		if events[i].UUID == id {
			return events[i], nil
		}
	}
	return types.Event{}, fmt.Errorf("id: '%s' was not found in '%#v'", id, events)
}

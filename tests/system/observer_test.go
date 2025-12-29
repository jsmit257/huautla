package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jsmit257/huautla/types"
	"github.com/stretchr/testify/require"
)

var events []types.Event

func init() {
	for _, id := range []types.UUID{
		"0",
		"1",
		"2",
		"add spore event source 0",
		"add spore event source 1",
		"add spore event source 2",
		"add clone event source 0",
		"get gen event 0",
		"get gen event 1",
		"get gen event 2",
	} {
		if e, err := db.SelectEvent(context.Background(), id, "lifecycleevent_init"); err != nil {
			panic(err)
		} else {
			events = append(events, e)
		}
	}
}

func Test_SelectByObservable(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		id     types.UUID
		result []types.Event
		err    error
	}{
		"happy_path": {
			id:     "get gen event",
			result: []types.Event{events[7], events[8], events[9]},
		},
		"no_rows_returned": {
			id:     "missing",
			result: []types.Event{},
		},
	}
	for name, tc := range set {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var actual types.Event
			result, err := db.SelectByObservable(context.Background(), tc.id, types.CID(name))
			require.Equal(t, tc.err, err)
			require.Equal(t, len(tc.result), len(result))
			for i, j := 0, len(tc.result); i < j; i++ {
				event := tc.result[i]
				actual, err = findEvent(result, tc.result[i].UUID)
				require.Nil(t, err)
				require.Equal(t, event.Temperature, actual.Temperature)
				require.Equal(t, event.Humidity, actual.Humidity)
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
			require.Equal(t, v.result.EventType, result.EventType)
		})
	}
}

func Test_InsertEvent(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		e   types.Event
		oid types.UUID
		err error
	}{
		"happy_lc_path": {
			oid: "lc insert event",
			e:   types.Event{UUID: "lc insert event", EventType: eventtypes[1]},
		},
		"happy_gen_path": {
			oid: "gen insert event",
			e:   types.Event{UUID: "gen insert event", EventType: eventtypes[1]},
		},
		"missing_event_type": {
			oid: "lc insert event",
			e:   types.Event{UUID: "lc insert fails", EventType: types.EventType{UUID: "missing"}},
			err: fmt.Errorf("event was not added"),
		},
		"missing observable": { // dunno how this would happen, but whatever
			e:   types.Event{UUID: "lc insert fails", EventType: eventtypes[1]},
			err: fmt.Errorf("event was not added"),
		},
	}
	for name, tc := range set {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := db.InsertEvent(context.Background(), tc.oid, tc.e, types.CID(name))
			require.Equal(t, tc.err, err)
			// require.Equal(t, tc.result, evt.EventType)
		})
	}
}

func Test_UpdateEvent(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		e   types.Event
		oid types.UUID
		err error
	}{
		"happy_path": { // happy path needs to run first, synchronously
			oid: "lc change event",
			e:   types.Event{UUID: "lc change event", EventType: eventtypes[1]},
		},
		"no_observable_affected": { // dunno how this would happen, but whatever
			e:   types.Event{UUID: "lc change event", EventType: eventtypes[1]},
			err: fmt.Errorf("observable was not changed"),
		},
		"no_event_affected": {
			e: types.Event{UUID: "missing", EventType: eventtypes[0]},
			// it fails because the event is missing, but since the observable is
			// updated first, that's the error we get
			err: fmt.Errorf("observable was not changed"),
		},
	}
	for name, tc := range set {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := db.UpdateEvent(context.Background(), tc.oid, tc.e, types.CID(name))
			require.Equal(t, tc.err, err)
			// require.Equal(t, tc.result, evt.EventType)
		})
	}
}

func Test_DeleteEvent(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		evid types.UUID
		oid  types.UUID
		err  error
	}{
		"happy_lc_path": { // happy path needs to run first, synchronously
			oid:  "lc delete event",
			evid: "lc delete event",
		},
		"happy_gen_path": { // happy path needs to run first, synchronously
			oid:  "gen delete event",
			evid: "gen delete event",
		},
		"observable_mismatch": {
			oid:  "gen delete event",
			evid: "get gen event 0",
			err:  fmt.Errorf("observable was not changed"),
		},
		// // other combinations of observable-/event- missing all fail for
		// // the same reason
		// "event_missing": {
		// 	oid:  "gen delete event",
		// 	evid: "get gen event 0",
		// 	err:  fmt.Errorf("observable was not changed"),
		// },
	}
	for name, tc := range set {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := db.DeleteEvent(context.Background(), tc.oid, tc.evid, types.CID(name))
			require.Equal(t, tc.err, err)
			// require.Equal(t, tc.result, evt.EventType)
		})
	}
}

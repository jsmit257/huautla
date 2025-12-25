package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

func Test_GetLifecycleEvents(t *testing.T) {
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
				require.Equal(t, event.EventType, actual.EventType)
			}
		})
	}
}

func Test_AddLifecycleEvent(t *testing.T) {
	t.Parallel()

	lc, err := db.SelectLifecycle(context.Background(), "add event", types.CID("Test_AddLifecycleEvent"))
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
			err := db.AddLifecycleEvent(context.Background(), &lc, v.e, types.CID(k))
			t.Logf("actual: %#v", lc.Events)
			require.Equal(t, v.err, err)
			require.Equalf(t, v.count, len(lc.Events), "actual: %#v", lc.Events)
		})
	}
}

func Test_ChangeLifecycleEvent(t *testing.T) {
	t.Parallel()

	lc, err := db.SelectLifecycle(context.Background(), "change event", types.CID("Test_ChangeLifecycleEvent"))
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
			_, err := db.ChangeLifecycleEvent(context.Background(), &lc, v.e, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, lc.Events[0].EventType)
		})
	}
}

func Test_RemoveLifecycleEvent(t *testing.T) {
	t.Parallel()

	lc, err := db.SelectLifecycle(context.Background(), "remove event", types.CID("Test_RemoveLifecycleEvent"))
	require.Nil(t, err)

	set := map[string]struct {
		id     types.UUID
		result []types.Event
		err    error
	}{
		"no_rows_affected": {
			id:     "missing event",
			result: lc.Events,
			err:    fmt.Errorf("event could not be removed"),
		},
		"happy_path": {
			id:     "remove event 2",
			result: []types.Event{lc.Events[0], lc.Events[2]},
		},
		"used_by_sources": {
			id:     "clone",
			result: lc.Events,
			err:    fmt.Errorf("pq: foreign key violation"),
		},
	}
	for k, v := range set {
		k, v, lc := k, v, lc
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			err := db.RemoveLifecycleEvent(context.Background(), &lc, v.id, types.CID(k))
			equalErrorMessages(t, v.err, err)
			require.Equal(t, v.result, lc.Events)
		})
	}
}

package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jsmit257/huautla/types"
)

func init() {
	for _, id := range []types.UUID{"get gen event 0", "get gen event 1", "get gen event 2"} {
		if e, err := db.SelectEvent(context.Background(), id, "generationevent_init"); err != nil {
			panic(err)
		} else {
			events = append(events, e)
		}
	}
}

func Test_GetGenerationEvents(t *testing.T) {
	t.Parallel()

	set := map[string]struct {
		g      types.Generation
		result []types.Event
		err    error
	}{
		"happy_path": {
			g: types.Generation{UUID: "get gen event"},
			result: func(e []types.Event) []types.Event {
				result := []types.Event{}
				e0, err := findEvent(e, "get gen event 0")
				require.Nil(t, err)
				result = append(result, e0)
				e1, err := findEvent(e, "get gen event 1")
				require.Nil(t, err)
				result = append(result, e1)
				e2, err := findEvent(e, "get gen event 2")
				require.Nil(t, err)
				result = append(result, e2)
				return result
			}(events),
		},
	}
	for k, v := range set {
		k, v := k, v
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			var actual types.Event
			err := db.GetGenerationEvents(context.Background(), &v.g, types.CID(k))
			require.Equal(t, v.err, err)
			for i, j := 0, len(v.result); i < j; i++ {
				event := v.result[i]
				actual, err = findEvent(v.g.Events, v.result[i].UUID)
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

func Test_AddGenerationEvent(t *testing.T) {
	t.Parallel()

	g, err := db.SelectGeneration(context.Background(), "add gen event", types.CID("Test_AddGenerationEvent"))
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
		k, v, g := k, v, g
		t.Run(k, func(t *testing.T) {
			// t.Parallel()
			err := db.AddGenerationEvent(context.Background(), &g, v.e, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equalf(t, v.count, len(g.Events), "actual: %#v", g.Events)
		})
	}
}

func Test_ChangeGenerationEvent(t *testing.T) {
	t.Parallel()

	g, err := db.SelectGeneration(context.Background(), "change event", types.CID("Test_ChangeGenerationEvent"))
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
		k, v, g := k, v, g
		g.Events = append([]types.Event{}, g.Events...)
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			_, err := db.ChangeGenerationEvent(context.Background(), &g, v.e, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, g.Events[0].EventType)
		})
	}
}

func Test_RemoveGenerationEvent(t *testing.T) {
	t.Parallel()

	g, err := db.SelectGeneration(context.Background(), "remove gen event", types.CID("Test_RemoveGenerationEvent"))
	require.Nil(t, err)
	require.Equal(t, 3, len(g.Events), "g: %v", g)

	set := map[string]struct {
		id     types.UUID
		result []types.Event
		err    error
	}{ // NB: the order in which the test cases are declared matters: errors should come first
		"no_rows_affected": {
			id:     "missing event",
			result: g.Events,
			err:    fmt.Errorf("event could not be removed"),
		},
		"notes_foreign_key": {
			id:     "event foreign key",
			result: g.Events,
			err:    fmt.Errorf("event could not be removed"),
			// err: fmt.Errorf("pq: foreign key violation"),
		},
		"photos_foreign_key": {
			id:     "gen photo 3",
			result: g.Events,
			err:    fmt.Errorf("event could not be removed"),
			// err: fmt.Errorf("pq: foreign key violation"),
		},
		"happy_path": {
			id:     "remove gen event 1",
			result: []types.Event{g.Events[0], g.Events[2]},
		},
	}
	for k, v := range set {
		k, v, g := k, v, g
		t.Run(k, func(t *testing.T) {
			t.Parallel() // XXX: don't do this, the data doesn't support it
			err := db.RemoveGenerationEvent(context.Background(), &g, v.id, types.CID(k))
			require.Equal(t, v.err, err)
			require.Equal(t, v.result, g.Events)
		})
	}
}

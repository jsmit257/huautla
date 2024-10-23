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
	} {
		if e, err := db.SelectEvent(context.Background(), id, "lifecycleevent_init"); err != nil {
			panic(err)
		} else {
			events = append(events, e)
		}
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

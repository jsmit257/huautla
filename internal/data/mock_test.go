package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jsmit257/huautla/types"
)

type (
	mocker  struct{ sqlmock.Sqlmock }
	row     []string
	xformer []driver.Value
	xform   map[int]any
)

var (
	attributes = []interface{}{
		mustObject(_attrs[0]),
		mustObject(_attrs[1]),
		mustObject(_attrs[2]),
	}
	album = []types.Entity{
		mustObject(_photos[0]),
		mustObject(_photos[1]),
		mustObject(_photos[2]),
	}
	ingredients = []interface{}{
		mustObject(_ingredients[0]),
		mustObject(_ingredients[1]),
		mustObject(_ingredients[2]),
	}
	events = []interface{}{
		mustObject(_events[0]),
		mustObject(_events[1]),
		mustObject(_events[2]),
	}
	notes = []types.Entity{
		mustEntity(_notes[0]),
		mustEntity(_notes[1]),
		mustEntity(_notes[2]),
	}
)

func mustJSON(o any) []byte {
	js, _ := json.Marshal(o)
	return js
}

func mustEntity(o any) types.Entity {
	result := types.Entity{}
	_ = json.Unmarshal(mustJSON(o), &result)
	return result
}

func mustObject(o any) map[string]interface{} {
	result := make(map[string]interface{})
	_ = json.Unmarshal(mustJSON(o), &result)
	return result
}

func (x xformer) replace(xforms ...xform) []driver.Value {
	result := make([]driver.Value, len(x))
	copy(result, x)

	for _, xform := range xforms {
		for k, v := range xform {
			result[k] = v
		}
	}

	return result
}

// deprecated: use set() in a builder instead
func (r row) mock(mock sqlmock.Sqlmock, rows ...[]driver.Value) {
	mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(r).AddRows(rows...))
}

func (r row) set(rows ...[]driver.Value) func(*mocker) *mocker {
	return func(mock *mocker) *mocker {
		mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(r).AddRows(rows...))
		return mock
	}
}

// this is what fail() throws, easy to test for and descriptive enough if it isn't caught
func (r row) err() error {
	return fmt.Errorf("fail error %v", r)
}

// this takes args so toggling set/fail for debugging doesn't require removing arg(s)
func (r row) fail(...any) func(*mocker) *mocker {
	return func(mock *mocker) *mocker {
		mock.ExpectQuery("").WillReturnError(r.err())
		return mock
	}
}

func newBuilder(mock sqlmock.Sqlmock, fns ...func(*mocker) *mocker) *mocker {
	return (&mocker{mock}).add(fns...)
}

// seems redundant, but you can do conditional branching from a common root (think
// lifecycles and generations: similar but not the same)
func (m *mocker) add(fns ...func(*mocker) *mocker) *mocker {
	for _, fn := range fns {
		fn(m)
	}
	return m
}

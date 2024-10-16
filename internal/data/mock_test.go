package data

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jsmit257/huautla/types"
)

type (
	attr       types.StrainAttribute
	event      types.Event
	ingredient types.Ingredient
	note       types.Note
	photo      types.Photo

	row     []string
	xformer []driver.Value
	xform   map[int]any
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

func (clz row) mock(mock sqlmock.Sqlmock, rows ...[]driver.Value) {
	mock.ExpectQuery("").WillReturnRows(sqlmock.NewRows(clz).AddRows(rows...))
}

package sqlitedb

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
)

var (
	_ driver.Value     = (*JSON[any])(nil)
	_ sql.Scanner      = (*JSON[any])(nil)
	_ json.Marshaler   = (*JSON[any])(nil)
	_ json.Unmarshaler = (*JSON[any])(nil)
	_ yaml.Marshaler   = (*JSON[any])(nil)
	_ yaml.Unmarshaler = (*JSON[any])(nil)
)

// WrapJSON wraps the given data into a JSON envelope, so that it can be stored in database as json format. Callers
// should provide only value types to data.
func WrapJSON[T any](data T) *JSON[T] {
	return &JSON[T]{wrapped: data}
}

// JSON is a generic envelope for wrapped data to be stored in database as json format.
//
// This type also implements yaml.UnmarshalYAML function. However, the implementation expects the node to be a string
// containing json data that can be in turn unmarshalled with json.Unmarshal function. This implementation is intended
// to make it easy when loading db fixtures so that the wrapped data can be read directly as json and do not have to
// define additional yaml field tags.
type JSON[T any] struct {
	wrapped T
}

// Unwrap returns the wrapped data.
func (j JSON[T]) Unwrap() T {
	return j.wrapped
}

// UnwrapRef returns a reference to the wrapped data.
func (j JSON[T]) UnwrapRef() *T {
	return &j.wrapped
}

func (j JSON[T]) Value() (driver.Value, error) {
	return json.Marshal(j.wrapped)
}

func (j *JSON[T]) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(bytes, &j.wrapped)
}

func (j JSON[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.wrapped)
}

func (j *JSON[T]) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &j.wrapped)
}

func (j JSON[T]) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(j.wrapped)
}

func (j *JSON[T]) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(&j.wrapped)
}

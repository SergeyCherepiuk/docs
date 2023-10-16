package internal

import (
	"context"
	"fmt"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// TODO: Review errors
var (
	ErrNoRecords           = fmt.Errorf("no records found")
	ErrNilRecord           = fmt.Errorf("record is nil")
	ErrInvalidValue        = fmt.Errorf("value is invalid")
	ErrCannotSetValue      = fmt.Errorf("value cannot be set")
	ErrVariableNotFound    = fmt.Errorf("variable not found")
	ErrInvalidVariableType = fmt.Errorf("invalid variable type")
	ErrPropertyNotFound    = fmt.Errorf("property not found")
	ErrTypeMismatch        = fmt.Errorf("value and node's property have different types")
)

func GetSingle[T any](ctx context.Context, result neo4j.ResultWithContext, variable string) (T, error) {
	var value T

	record, err := result.Single(ctx)
	if record == nil || err != nil {
		return value, ErrNoRecords
	}

	if reflect.TypeOf(value).Kind() == reflect.Struct {
		return collectStruct[T](ctx, record, variable)
	}
	return collectPrimitive[T](ctx, record, variable)
}

func GetMultiple[T any](ctx context.Context, result neo4j.ResultWithContext, variable string) ([]T, error) {
	records, err := result.Collect(ctx)
	if records == nil || len(records) <= 0 || err != nil {
		return nil, ErrNoRecords
	}

	values := make([]T, len(records))
	for i, record := range records {
		var value T
		var err error
		if reflect.ValueOf(value).Kind() == reflect.Struct {
			value, err = collectStruct[T](ctx, record, variable)
		} else {
			value, err = collectPrimitive[T](ctx, record, variable)
		}
		if err != nil {
			return nil, err
		}
		values[i] = value
	}

	return values, nil
}

func collectPrimitive[T any](ctx context.Context, record *neo4j.Record, variable string) (T, error) {
	var value T

	if record == nil {
		return value, ErrNilRecord
	}

	rv := reflect.ValueOf(&value).Elem()
	if !rv.IsValid() {
		return value, ErrInvalidValue
	} else if !rv.CanSet() {
		return value, ErrCannotSetValue
	}

	primitive, found := record.Get(variable)
	if !found {
		return value, ErrVariableNotFound
	}

	pv := reflect.ValueOf(primitive)
	if rv.Kind() != pv.Kind() {
		return value, ErrTypeMismatch
	}

	rv.Set(pv)

	return value, nil
}

func collectStruct[T any](ctx context.Context, record *neo4j.Record, variable string) (T, error) {
	var value T

	if record == nil {
		return value, ErrNilRecord
	}

	v, found := record.Get(variable)
	if !found {
		return value, ErrVariableNotFound
	}

	rt := reflect.TypeOf(&value).Elem()
	rv := reflect.ValueOf(&value).Elem()

	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("prop")
		if tag == "" {
			continue
		}

		fv := rv.Field(i)
		if !fv.IsValid() {
			return value, ErrInvalidValue
		} else if !fv.CanSet() {
			return value, ErrCannotSetValue
		}

		var property any
		var found bool

		switch v := v.(type) {
		case neo4j.Node:
			property, found = v.Props[tag]
		case map[string]any:
			property, found = v[tag]
		default:
			return value, ErrInvalidVariableType
		}

		if !found {
			return value, ErrPropertyNotFound
		}

		pv := reflect.ValueOf(property)
		if fv.Kind() != pv.Kind() {
			return value, ErrTypeMismatch
		}

		fv.Set(pv)
	}

	return value, nil
}

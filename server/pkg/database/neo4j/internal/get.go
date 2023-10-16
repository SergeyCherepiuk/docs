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
	ErrInvalidVariable     = fmt.Errorf("variable is invalid")
	ErrVariableCannotBeSet = fmt.Errorf("variable cannot be set")
	ErrValueNotFound       = fmt.Errorf("value not found")
	ErrInvalidValueType    = fmt.Errorf("invalid value type")
	ErrTypeMismatch        = fmt.Errorf("value and node's property have different types")
)

func GetSingle[T any](ctx context.Context, result neo4j.ResultWithContext, alias string) (T, error) {
	var variable T

	record, err := result.Single(ctx)
	if record == nil || err != nil {
		return variable, ErrNoRecords
	}

	if reflect.TypeOf(variable).Kind() == reflect.Struct {
		return collectStruct[T](ctx, record, alias)
	}
	return collectPrimitive[T](ctx, record, alias)
}

func GetMultiple[T any](ctx context.Context, result neo4j.ResultWithContext, alias string) ([]T, error) {
	records, err := result.Collect(ctx)
	if records == nil || len(records) <= 0 || err != nil {
		return nil, ErrNoRecords
	}

	variables := make([]T, len(records))
	for i, record := range records {
		var variable T
		var err error
		if reflect.ValueOf(variable).Kind() == reflect.Struct {
			variable, err = collectStruct[T](ctx, record, alias)
		} else {
			variable, err = collectPrimitive[T](ctx, record, alias)
		}
		if err != nil {
			return nil, err
		}
		variables[i] = variable
	}

	return variables, nil
}

func collectPrimitive[T any](ctx context.Context, record *neo4j.Record, alias string) (T, error) {
	var variable T

	if record == nil {
		return variable, ErrNilRecord
	}

	rv := reflect.ValueOf(&variable).Elem()
	if !rv.IsValid() {
		return variable, ErrInvalidVariable
	} else if !rv.CanSet() {
		return variable, ErrVariableCannotBeSet
	}

	value, found := record.Get(alias)
	if !found {
		return variable, ErrValueNotFound
	}

	pv := reflect.ValueOf(value)
	if rv.Kind() != pv.Kind() {
		return variable, ErrTypeMismatch
	}

	rv.Set(pv)

	return variable, nil
}

func collectStruct[T any](ctx context.Context, record *neo4j.Record, alias string) (T, error) {
	var variable T

	if record == nil {
		return variable, ErrNilRecord
	}

	value, found := record.Get(alias)
	if !found {
		return variable, ErrValueNotFound
	}

	rt := reflect.TypeOf(&variable).Elem()
	rv := reflect.ValueOf(&variable).Elem()

	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("prop")
		if tag == "" {
			continue
		}

		fv := rv.Field(i)
		if !fv.IsValid() {
			return variable, ErrInvalidVariable
		} else if !fv.CanSet() {
			return variable, ErrVariableCannotBeSet
		}

		var property any
		var found bool

		switch value := value.(type) {
		case neo4j.Node:
			property, found = value.Props[tag]
		case map[string]any:
			property, found = value[tag]
		default:
			return variable, ErrInvalidValueType
		}

		if !found {
			return variable, ErrValueNotFound
		}

		pv := reflect.ValueOf(property)
		if fv.Kind() != pv.Kind() {
			return variable, ErrTypeMismatch
		}

		fv.Set(pv)
	}

	return variable, nil
}

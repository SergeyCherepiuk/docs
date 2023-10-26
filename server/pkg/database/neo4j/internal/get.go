package internal

import (
	"context"
	"fmt"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type (
	ErrorNoRecords        error
	ErrorNilRecord        error
	ErrorInvalidValue     error
	ErrorValueCannotBeSet error
	ErrorAliasNotFound    error
	ErrorTypeMismatch     error
	ErrorInvalidAliasType error
)

func GetSingle[T any](ctx context.Context, result neo4j.ResultWithContext, alias string) (T, error) {
	var variable T

	record, err := result.Single(ctx)
	if err != nil {
		return variable, ErrorNoRecords(fmt.Errorf("no record was found"))
	}
	if record == nil {
		return variable, ErrorNilRecord(fmt.Errorf("record is nil"))
	}

	if reflect.TypeOf(variable).Kind() == reflect.Struct {
		return collectStruct[T](ctx, record, alias)
	}
	return collectPrimitive[T](ctx, record, alias)
}

func GetMultiple[T any](ctx context.Context, result neo4j.ResultWithContext, alias string) ([]T, error) {
	records, err := result.Collect(ctx)
	if len(records) <= 0 || err != nil {
		return nil, ErrorNoRecords(fmt.Errorf("no records were found"))
	}

	variables := make([]T, len(records))
	for i, record := range records {
		if record == nil {
			return nil, ErrorNilRecord(fmt.Errorf("record is nil"))
		}

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
		return variable, ErrorNilRecord(fmt.Errorf("record is nil"))
	}

	rv := reflect.ValueOf(&variable).Elem()
	if !rv.IsValid() {
		return variable, ErrorInvalidValue(fmt.Errorf("variable is not valid"))
	} else if !rv.CanSet() {
		return variable, ErrorValueCannotBeSet(fmt.Errorf("variable cannot be set"))
	}

	value, found := record.Get(alias)
	if !found {
		return variable, ErrorAliasNotFound(fmt.Errorf("alias \"%s\" not found", alias))
	}

	pv := reflect.ValueOf(value)
	if rv.Kind() != pv.Kind() {
		return variable, ErrorTypeMismatch(fmt.Errorf("cannot set value of a type %s, to a variable of a type %s", pv.Kind(), rv.Kind()))
	}

	rv.Set(pv)

	return variable, nil
}

func collectStruct[T any](ctx context.Context, record *neo4j.Record, alias string) (T, error) {
	var variable T

	if record == nil {
		return variable, ErrorNilRecord(fmt.Errorf("record is nil"))
	}

	value, found := record.Get(alias)
	if !found {
		return variable, ErrorAliasNotFound(fmt.Errorf("alias \"%s\" not found", alias))
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
			return variable, ErrorInvalidValue(fmt.Errorf("struct field in not valid"))
		} else if !fv.CanSet() {
			return variable, ErrorValueCannotBeSet(fmt.Errorf("struct field cannot be set"))
		}

		var property any
		var found bool
		switch value := value.(type) {
		case neo4j.Node:
			property, found = value.Props[tag]
		case map[string]any:
			property, found = value[tag]
		default:
			return variable, ErrorInvalidAliasType(fmt.Errorf("invalid alias type: %T", value))
		}

		if !found {
			return variable, ErrorAliasNotFound(fmt.Errorf("alias \"%s\" not found", alias))
		}

		pv := reflect.ValueOf(property)
		if fv.Kind() != pv.Kind() {
			return variable, ErrorTypeMismatch(fmt.Errorf("cannot set value of a type %s, to a variable of a type %s", pv.Kind(), rv.Kind()))
		}

		fv.Set(pv)
	}

	return variable, nil
}

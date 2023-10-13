package internal

import (
	"context"
	"fmt"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetSingle[T any](ctx context.Context, result neo4j.ResultWithContext) (T, error) {
	record, err := result.Single(ctx)
	if record == nil || err != nil {
		return *new(T), fmt.Errorf("failed to get the record from result")
	}

	return collect[T](ctx, record)
}

func GetMultiple[T any](ctx context.Context, result neo4j.ResultWithContext) ([]T, error) {
	records, err := result.Collect(ctx)
	if records == nil || len(records) <= 0 || err != nil {
		return nil, fmt.Errorf("failed to get the records from result")
	}

	values := make([]T, 0, len(records))
	for i := 0; i < len(records); i++ {
		value, err := collect[T](ctx, records[i])
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}

func collect[T any](ctx context.Context, record *neo4j.Record) (T, error) {
	if record == nil {
		return *new(T), fmt.Errorf("record is nil")
	}

	value := new(T)
	valueType := reflect.TypeOf(value).Elem()
	valueValue := reflect.ValueOf(value).Elem()

	for i := 0; i < valueType.NumField(); i++ {
		tag := valueType.Field(i).Tag.Get("prop")
		if tag == "" {
			continue
		}

		fieldValue := valueValue.Field(i)
		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			return *value, fmt.Errorf("the '%s' field is invalid or cannot be set", tag)
		}

		recordProp, found := record.Get(tag)
		if !found {
			return *value, fmt.Errorf("the '%s' field is not found in node", tag)
		}

		fieldValue.Set(reflect.ValueOf(recordProp))
	}

	return *value, nil
}

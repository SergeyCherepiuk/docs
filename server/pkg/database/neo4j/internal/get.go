package internal

import (
	"context"
	"fmt"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// TODO: Refactor and prepare for being used recursively (if the field is a struct with "prop" tags)
func GetSingle[T any](ctx context.Context, result neo4j.ResultWithContext, value *T) error {
	record, err := result.Single(ctx)
	if err != nil {
		return err
	}

	valueType := reflect.TypeOf(value).Elem()
	valueValue := reflect.ValueOf(value).Elem()

	for i := 0; i < valueValue.NumField(); i++ {
		fieldType := valueType.Field(i)
		fieldValue := valueValue.Field(i)

		tag := fieldType.Tag.Get("prop")
		if tag == "" {
			continue
		}

		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			return fmt.Errorf("one of the field is invalid or cannot be set")
		}

		recordProp, found := record.Get(tag)
		if !found {
			return fmt.Errorf("one of the field is not found")
		}

		fieldValue.Set(reflect.ValueOf(recordProp))
	}

	return nil
}

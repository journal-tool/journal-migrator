package lib

import (
	"database/sql"
	"reflect"
)

func StructFieldNames(T any) []string {
	var value = reflect.ValueOf(T).Elem()
	var fieldNumber = value.Type().NumField()
	var fieldNames = make([]string, fieldNumber)

	for i := range fieldNumber {
		fieldNames[i] = value.Type().Field(i).Name
	}

	return fieldNames
}

func StructFieldPointers(T any) []any {
	var value = reflect.ValueOf(T).Elem()
	var fieldNumber = value.Type().NumField()
	var fieldPointers = make([]any, fieldNumber)

	for i := range fieldNumber {
		fieldPointers[i] = value.Field(i).Addr().Interface()
	}

	return fieldPointers
}

func ParseRow[T any](row *sql.Row) (*T, error) {
	var rowStruct = new(T)
	var fieldsPointers = StructFieldPointers(rowStruct)

	var err = row.Scan(fieldsPointers...)
	if err != nil {
		return nil, err
	}

	return rowStruct, nil
}

func ParseRows[T any](rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	var rowStructs []T
	var rowStruct = new(T)
	var fieldsPointers = StructFieldPointers(rowStruct)

	for rows.Next() {
		var err = rows.Scan(fieldsPointers...)
		if err != nil {
			return rowStructs, err
		}

		rowStructs = append(rowStructs, *rowStruct)
	}

	var err = rows.Err()
	if err != nil {
		return nil, err
	}

	return rowStructs, nil
}

package dataaccess //nolint:revive // established package name, changing would break API

import (
	"reflect"

	"gorm.io/gorm"
)

// NewTypedDataAccess creates a new DataAccess for a specific model type using generics
func NewTypedDataAccess[T any](db *gorm.DB) DataAccess {
	var instance *T
	modelType := reflect.TypeOf(instance)
	return NewDataAccess(db, modelType)
}

// NewTypedDataAccessWithInterface creates a new DataAccess with interface for testing using generics
func NewTypedDataAccessWithInterface[T any](db GormDBInterface) DataAccess {
	var instance *T
	modelType := reflect.TypeOf(instance)
	return NewDataAccessWithInterface(db, modelType)
}

// InsertRow inserts a data row with type safety using generics
func InsertRow[T any](da DataAccess, datarow *T) error {
	return da.Insert(datarow)
}

// SelectRows selects data with type safety using generics
func SelectRows[T any](da DataAccess, datafilter *T, preloads ...string) ([]*T, error) {
	result, err := da.Select(datafilter, preloads...)
	if err != nil {
		return nil, err
	}

	// Type assert the result to []*T
	if typedResult, ok := result.([]*T); ok {
		return typedResult, nil
	}

	// If type assertion fails, return zero value
	var zero []*T
	return zero, err
}

// UpdateRow updates data with type safety using generics
func UpdateRow[T any](da DataAccess, datafilter *T, datarow *T) error {
	return da.Update(datafilter, datarow)
}

// DeleteRows deletes data with type safety using generics
func DeleteRows[T any](da DataAccess, datafilter *T, associateds ...string) error {
	return da.Delete(datafilter, associateds...)
}

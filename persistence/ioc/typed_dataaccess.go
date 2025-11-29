package ioc

import (
	"reflect"

	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
)

// TypedDataAccessModule provides convenience methods for typed data access
type TypedDataAccessModule struct{}

// NewTypedDataAccessModule creates a new typed data access module
func NewTypedDataAccessModule() *TypedDataAccessModule {
	return &TypedDataAccessModule{}
}

// RegisterServices registers convenience services for typed data access
func (m *TypedDataAccessModule) RegisterServices(register dependencyinjection.Register) error {
	// This module provides helper functions but doesn't register services itself
	// Services are registered by the specific database modules
	return nil
}

// Helper functions for creating typed data access instances

// NewTypedDataAccess is a convenience function to create typed data access
func NewTypedDataAccess[T any](resolver dependencyinjection.Resolver) dataaccess.DataAccess {
	modelType := reflect.TypeOf((*T)(nil)).Elem()
	result, ok := resolver.Type((*dataaccess.DataAccess)(nil), map[string]interface{}{
		"modelType": modelType,
	}).(dataaccess.DataAccess)
	if !ok {
		panic("failed to resolve dataaccess.DataAccess")
	}
	return result
}

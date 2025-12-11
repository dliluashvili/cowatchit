package validators

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Repository[T any] interface {
	FindByField(fields map[string]any) (*T, error)
}

// @Todo more generic
// It checks in a db
func Unique[T any](repo Repository[T]) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {

		record, err := repo.FindByField(map[string]any{
			fl.FieldName(): fl.Field().String(),
		})

		if err != nil {
			fmt.Println("err", err)
			panic("error !")
		}

		return record == nil
	}
}

// It checks in slice
func UniqueItem(fl validator.FieldLevel) bool {
	items, ok := fl.Field().Interface().([]string)

	fmt.Println(items)
	if !ok {
		return false
	}

	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		if _, exists := seen[item]; exists {
			return false
		}
		seen[item] = struct{}{}
	}

	return true
}

package upstruct

import (
	"errors"

	"github.com/fatih/structs"
)

func Update(target any, update any) error {
	if !structs.IsStruct(target) || !structs.IsStruct(update) {
		return errors.New("arguments must be structs")
	}

	targetFields := structs.Fields(target)
	updateFields := structs.Fields(update)

	for _, targetField := range targetFields {
		for _, updateField := range updateFields {
			if targetField.Name() != updateField.Name() {
				continue
			}

			if structs.IsStruct(targetField.Value()) && structs.IsStruct(updateField.Value()) {
				targetFieldValue, _ := targetField.Value().(any)
				updateFieldValue, _ := updateField.Value().(any)
				Update(&targetFieldValue, &updateFieldValue)

				targetField.Set(targetFieldValue)
			}

			if !updateField.IsZero() && targetField.Kind() == updateField.Kind() {
				targetField.Set(updateField.Value())
			}

			break
		}
	}

	return nil
}

func UpdateFn(target any, update any, fieldHandler ...func(*structs.Field, *structs.Field)) error {
	if !structs.IsStruct(target) || !structs.IsStruct(update) {
		return errors.New("arguments must be structs")
	}

	targetFields := structs.Fields(target)
	updateFields := structs.Fields(update)

	for _, targetField := range targetFields {
		for _, updateField := range updateFields {
			if targetField.Name() != updateField.Name() {
				continue
			}

			// add own configs
			for _, handlerFunc := range fieldHandler {
				handlerFunc(targetField, updateField)
			}

			if structs.IsStruct(targetField.Value()) && structs.IsStruct(updateField.Value()) {
				targetFieldValue, _ := targetField.Value().(any)
				updateFieldValue, _ := updateField.Value().(any)
				UpdateFn(&targetFieldValue, &updateFieldValue, fieldHandler...)

				targetField.Set(targetFieldValue)
			}

			if !updateField.IsZero() && targetField.Kind() == updateField.Kind() {
				targetField.Set(updateField.Value())
			}

			break
		}
	}

	return nil
}

package upstruct

import (
	"errors"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

type UpdateOption interface {
	IsValid() (bool, error)
	IsConditionMet(targetValue interface{}, updateValue interface{}) bool
}

type DifferentTypesOption struct {
	TargetType string
	UpdateType string
}

func (d DifferentTypesOption) IsValid() (bool, error) {
	if strings.TrimSpace(d.TargetType) == "" {
		return false, errors.New("target type is not specified")
	}
	if strings.TrimSpace(d.UpdateType) == "" {
		return false, errors.New("update type is not specified")
	}
	return true, nil
}

func (d DifferentTypesOption) IsConditionMet(targetValue interface{}, updateValue interface{}) bool {
	return d.TargetType == reflect.TypeOf(targetValue).String() && d.UpdateType == reflect.TypeOf(updateValue).String()
}

type SameKindOption struct {
	Kind reflect.Kind
}

func (s SameKindOption) IsValid() (bool, error) {
	if s.Kind == reflect.Invalid {
		return false, errors.New("invalid kind")
	}

	return true, nil
}

func (d SameKindOption) IsConditionMet(targetValue interface{}, updateValue interface{}) bool {
	return d.Kind == reflect.ValueOf(targetValue).Kind() && d.Kind == reflect.ValueOf(updateValue).Kind()
}

// ------------------------------------------------

type OptionHandler func(target, update *structs.Field)

type UpdateStructOptions struct {
	Option  UpdateOption
	Handler OptionHandler
}

// --------------------------------------------------

func extractFields(x any, isTarget bool) (*reflect.Value, []*structs.Field) {
	y := reflect.ValueOf(x)

	if y.Elem().Kind() == reflect.Interface {
		elem := y.Elem().Elem()

		if isTarget {
			newValue := reflect.New(elem.Type())
			return &newValue, structs.Fields(newValue.Interface())
		}

		return nil, structs.Fields(elem.Interface())
	}

	return nil, structs.Fields(x)
}

func isStruct(s any) bool {

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	// uninitialized zero value of a struct
	if v.Kind() == reflect.Invalid {
		return false
	}

	return v.Kind() == reflect.Struct
}

func Update(target any, update any, fieldHandlers ...UpdateStructOptions) error {
	// logger := getLogger()

	if !isStruct(target) || !isStruct(update) {
		return errors.New("arguments must be structs")
	}

	structVal, targetFields := extractFields(target, true)
	_, updateFields := extractFields(update, false)

	for _, targetField := range targetFields {
	InnerLoop:
		for _, updateField := range updateFields {
			if targetField.Name() != updateField.Name() {
				continue
			}

			for _, handler := range fieldHandlers {
				_, err := handler.Option.IsValid()
				if err != nil {
					return err
				}

				if handler.Option.IsConditionMet(targetField.Value(), updateField.Value()) {
					handler.Handler(targetField, updateField)
					break InnerLoop
				}
			}

			if structs.IsStruct(targetField.Value()) && structs.IsStruct(updateField.Value()) {
				targetFieldValue := targetField.Value()
				updateFieldValue := updateField.Value()
				Update(&targetFieldValue, &updateFieldValue, fieldHandlers...)

				targetFieldReflectValue := reflect.ValueOf(targetFieldValue)

				if targetFieldReflectValue.Kind() == reflect.Ptr {
					targetFieldReflectValue = targetFieldReflectValue.Elem()
				}

				targetField.Set(targetFieldReflectValue.Interface())
				break
			}

			if !updateField.IsZero() && targetField.Kind() == updateField.Kind() {
				targetField.Set(updateField.Value())
			}

			break
		}
	}

	if structVal != nil {
		targetElem := reflect.ValueOf(target).Elem()

		targetElem.Set(*structVal)
	}

	return nil
}

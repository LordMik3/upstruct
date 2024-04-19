package upstruct

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/rs/zerolog"
)

type UpdateOption interface {
	IsValid() (bool, error)
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

type SameKindOption struct {
	Kind reflect.Kind
}

func (s SameKindOption) IsValid() (bool, error) {
	if s.Kind == reflect.Invalid {
		return false, errors.New("invalid kind")
	}

	return true, nil
}

// ------------------------------------------------

type OptionHandler func(target, update *structs.Field)

type UpdateStructOptions struct {
	Option  UpdateOption
	Handler OptionHandler
}

func NewUpdateStructOption[T UpdateOption](updateOption T, handler OptionHandler) *UpdateStructOptions {
	return &UpdateStructOptions{
		Option:  updateOption,
		Handler: handler,
	}
}

// --------------------------------------------------

func getLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stdout,
	}).
		With().
		Caller().
		Timestamp().
		Logger()
}

func Update(target any, update any, fieldHandler ...UpdateStructOptions) error {
	logger := getLogger()
	if !structs.IsStruct(target) || !structs.IsStruct(update) {
		return errors.New("arguments must be structs")
	}

	targetFields := structs.Fields(target)
	updateFields := structs.Fields(update)

	for _, targetField := range targetFields {
	InnerLoop:
		for _, updateField := range updateFields {
			if targetField.Name() != updateField.Name() {
				continue
			}

			for _, handler := range fieldHandler {
				_, err := handler.Option.IsValid()
				if err != nil {
					return err
				}

				conditionMet := false

				value, ok := handler.Option.(SameKindOption)
				if !ok {
					value, _ := handler.Option.(DifferentTypesOption)
					if value.TargetType == reflect.TypeOf(targetField.Value()).String() &&
						value.UpdateType == reflect.TypeOf(updateField.Value()).String() {
						conditionMet = true
					}
				}

				if value.Kind == targetField.Kind() && value.Kind == updateField.Kind() {
					conditionMet = true
				}

				if conditionMet {
					handler.Handler(targetField, updateField)
					break InnerLoop
				}
			}

			if structs.IsStruct(targetField.Value()) && structs.IsStruct(updateField.Value()) {
				targetFieldAny, _ := targetField.Value().(any)
				updateFieldAny, _ := updateField.Value().(any)

				// TODO: fix the interface type shit 
				Update(&targetFieldAny, &updateFieldAny, fieldHandler...)
				logger.Info().Msgf("targetField new value=%+v", targetFieldAny)
				logger.Info().Msgf("updateField value=%+v", updateFieldAny)
				targetField.Set(targetFieldAny)
				break
			}

			if !updateField.IsZero() && targetField.Kind() == updateField.Kind() {
				logger.Info().Msgf("targetField name=%s", targetField.Name())
				logger.Info().Msgf("updateField name=%s", updateField.Name())
				targetField.Set(updateField.Value())
			}

			break
		}
	}

	return nil
}

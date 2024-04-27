package upstruct_test

import (
	"database/sql"
	"testing"

	"github.com/LordMik3/upstruct"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
)

type address struct {
	StreetNumber uint16
	StreetName   string
}

type addressDTO struct {
	StreetNumber uint16
	StreetName   string
}

type target struct {
	Email   sql.NullString
	Address address
}

type update struct {
	Email   string
	Address addressDTO
}

func TestUpdateFn(t *testing.T) {

	target := target{
		Email: sql.NullString{
			String: "",
			Valid:  false,
		},
		Address: address{
			StreetNumber: 105,
			StreetName:   "cpu street",
		},
	}

	update := update{
		Email: "test@gmail.com",
		Address: addressDTO{
			StreetName: "gpu street",
		},
	}

	upstruct.Update(&target, &update, upstruct.UpdateStructOption{
		Option: upstruct.DifferentTypesOption{
			TargetType: "sql.NullString",
			UpdateType: "string",
		},
		Handler: func(target, update *structs.Field) {
			target.Set(sql.NullString{
				String: update.Value().(string),
				Valid:  update.Value().(string) != "",
			})
		},
	})

	assert.Equal(t, update.Email, target.Email.String)
	assert.Equal(t, update.Address.StreetName, target.Address.StreetName)
}

package upstruct_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/LordMik3/upstruct"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
)

type address struct {
	StreetNumber uint16
	StreetName   string
}

type target struct {
	Name    string
	Age     uint8
	Email   sql.NullString
	Address address
}

type update struct {
	Name    string
	Age     uint8
	Email   string
	Address address
}

func TestUpdateFn(t *testing.T) {

	target := target{
		Name: "upstruct",
		Age:  18, // legal
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
		Name:  "up struct",
		Email: "test@gmail.com",
		Address: address{
			StreetName: "gpu street",
		},
	}

	upstruct.UpdateFn(&target, &update, func(f1, f2 *structs.Field) {
		if reflect.TypeOf(f1.Value()).String() == "sql.NullString" {
			f1.Set(sql.NullString{
				String: f2.Value().(string),
				Valid:  true,
			})
		}
	})

	assert.Equal(t, update.Email, target.Email.String)
	assert.Equal(t, update.Name, target.Name)
	assert.Equal(t, update.Address.StreetName, target.Address.StreetName)
}

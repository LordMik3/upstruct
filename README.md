# upstruct [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/LordMik3/upstruct)

upstruct is an utility package to update a golang struct with another one

This is a fork of [upstruct](https://github.com/hackirby/upstruct), go show some love if you like it.

## Install

```bash
go get github.com/LordMik3/upstruct
```

## Example Update

```go
import "github.com/hackirby/upstruct"

type User struct {
    Username string
    Password string
}

type UserPatch struct {
    Username string
    Password string
}

var user = User{
    Username: "user",
    Password: "password",
}

var userPatch = UserPatch{
    Username: "newuser",
    Password: "newpassword",
}

func main() {
    upstruct.Update(&user, userPatch)
}
```

## Example UpdateFn

```go
import (
    "reflect"

    "github.com/hackirby/upstruct"
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

var target := target{
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

var update := update{
        Name:  "up struct",
        Email: "test@gmail.com",
        Address: address{
            StreetName: "gpu street",
        },
    }

func main(){
    upstruct.UpdateFn(&target, &update, func(f1, f2 *structs.Field) {
        if reflect.TypeOf(f1.Value()).String() == "sql.NullString" {
            f1.Set(sql.NullString{
                String: f2.Value().(string),
                Valid:  true,
            })
        }
    })
}
```

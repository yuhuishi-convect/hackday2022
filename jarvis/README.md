# JARVIS - Just A Rather Very Intelligent System

A dynamic data platform to support easily modeling schema and database queries (through CRUD APIs) on the fly.

## Get started
### Schema Modeling
* install `go` through `Homebrew` or download it from https://golang.org/dl/ (version > 1.15)
* go to repo folder
* run `go get entgo.io/ent/cmd/ent`
* run `go run entgo.io/ent/cmd/ent init User`
    > The command above will generate the schema for User under <repo>/ent/schema/ directory.
* run `go run generator.go` under <repo>/generator folder
    > The command above will generate CRUD APIs for the User schema
* run `go run proxy.go` under <repo>/proxy folder
    > The command above will start a HTTP server
* goto url http://localhost:8080/uploadSchema (either on browser or on Postman)
    > The endpoint will automatically create User schema in a sqlite database (ent.db) under <repo>/proxy folder. Do `$ sqlite3 ent.db` to verify it.
* update `ent/schema/user.go` with following changes,
```
// update import statement
import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// add age field and name field
// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int("age").
            Positive(),
        field.String("name").
            Default("unknown"),
    }
}
```
* run `go run generator.go` under <repo>/generator folder
* rerun `go run proxy.go` under <repo>/proxy folder
> The new User schema will be reflected in sqlite database

### CRUD APIs
* keep `proxy.go` server running after above steps
* Send POST request to http://localhost;8080/api/user with body (set age value and name value) to create a new User record
* Send GET request to http://localhost;8080/api/user to get all User records
* Send GET request to http://localhost;8080/api/user/1 to get the `id == 1` User record
* Send DELETE request to http://localhost;8080/api/user/1 to delete the `id == 1` User record

## TODO
[ ] remove the dependency of entgo



# JARVIS - Just A Rather Very Intelligent System

A dynamic data platform to support easily modeling schema and querying and manipulating data from a database on the fly.

## Goal

[ORM](https://en.wikipedia.org/wiki/Object%E2%80%93relational_mapping) is a technique that lets developers query and manipulate data from a database using an object-oriented paradigm. Using ORM saves a lot of time because:

- [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself) - it's easier to update, maintain, and reuse the code
- No need to write poorly-formed SQL
- Easy to use with a lot of stuff that is done automatically - data operations.

Bu the traditional ORM library can be a pain:

- Model definitions is done programmatically
- Unable to update model definitions once the program starts running

```
engine = create_engine('sqlite:///:memory:')

metadata_obj = MetaData()

employees = Table('employees', metadata_obj,
    Column('employee_id', Integer, primary_key=True),
    Column('employee_name', String(60), nullable=False, key='name'),
    Column('employee_dept', Integer, ForeignKey("departments.department_id"))
)
sql
employees.create(engine)

```

The JARVIS aims to address these two pains while keeping the pros from the traditional ORM library.


## Get started
### Schema Modeling
* install `go` through `Homebrew` or download it from [https://golang.org/dl/](https://golang.org/dl/) (version > 1.15)
* go to repo folder
* run `go get entgo.io/ent/cmd/ent`
* run `go run entgo.io/ent/cmd/ent init User`
    > The command above will generate the schema for the User under <repo>/ent/schema/ directory.
* run `go run generator.go` under <repo>/generator folder
    > The command above will generate CRUD APIs for the User schema
* run `go run proxy.go` under <repo>/proxy folder
    > The command above will start a HTTP server
* goto url [http://localhost:8080/uploadSchema](http://localhost:8080/uploadSchema) (either on the browser or on Postman)
    > The endpoint will automatically create User schema in a sqlite database (ent.db) under <repo>/proxy folder. Do $ sqlite3 ent.db to verify it.
* update `ent/schema/user.go` with the following changes,
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
[ ] support model definition in metadata file
```
type Employee mix Entity schema name "employee" {
  id: string
  name: string
  department: Department
}

```



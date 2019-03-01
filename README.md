# Medium example

Este repositorio es utilizado en un nuevo post de medium que explica el uso de
GraphQL, PostgreSQL, Goose y Go.


## Instalar goose

```
$ go get -tags nosqlite3 github.com/steinbacher/goose/cmd/goose
```

### Migraciones para postgres

```
$ goose -path database/postgres -env docker create base
```


## Run server

```
$ go run cmd/server/main.go
```

# queries

### To get all users

```
query {
    users {
        id,
        first_name,
        last_name,
        email,
        phone,
        birthdate
    }
}
```

### To get an specific user

```
query {
    user(id: "<id>") {
        id,
        first_name,
        last_name,
        email,
        phone,
        birthdate,
        address {
            address_line,
            city,
            locality,
            region,
            country,
            postal_code
        }
    }
}
```


R.
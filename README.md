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

# Mutation

```
mutation {
  createUser(
    input: {
      email: "john.doe@example.org"
      first_name: "John"
      last_name: "Doe"
      phone: "912345678"
      birthdate: "1990-01-06T11:30:00+00:00"
      address: {
        address_line: "4129 Sycamore Lake Road"
        city: "Green Bay"
        locality: "Wisconsin"
        region: "Wisconsin"
        country: "USA"
        postal_code: 54304
      }
    }
  ) {
    id
    email
    first_name
    last_name
    phone
    address {
      address_line
      city
      locality
      region
      country
      postal_code
    }
  }
}
```

R.

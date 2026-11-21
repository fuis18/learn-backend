# 03-api-json

## Comandos

```sh
go run .
```

## CLI Tests

```sh
# Read the page and the health check
xh GET http://localhost:8080/
xh GET http://localhost:8080/health

# See the products (JSON)
xh GET http://localhost:8080/api/products

# Create a product: JSON → 201 + Location
xh POST http://localhost:8080/api/products name=Keyboard price:=50.5 stock:=3

# Get one, patch only the sent fields, replace the whole product
xh GET http://localhost:8080/api/products/1
xh PATCH http://localhost:8080/api/products/1 stock:=10
xh PUT http://localhost:8080/api/products/1 name=Keyboard price:=80 stock:=5

# Delete → 204
xh DELETE http://localhost:8080/api/products/1

# Errors: bad id → 400, missing product → 404, unknown field → 400
xh GET http://localhost:8080/api/products/abc
xh GET http://localhost:8080/api/products/99
xh POST http://localhost:8080/api/products nme=Mouse
```

## Tests

```sh
hurl --test test.hurl
```

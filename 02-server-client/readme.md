# 02-server-client

## Commands

```sh
go run .
```

## CLI Tests

```sh
# Read the page and the health check
xh GET http://localhost:8080/
xh GET http://localhost:8080/health

# See the products
xh GET http://localhost:8080/products

# Create a product: POST form → 303 redirect (PRG)
xh --form POST http://localhost:8080/products name=Keyboard price=50.5 stock=3

# Create and follow the redirect to see the message and the new row
xh -F --form POST http://localhost:8080/products name=Mouse price=25.99 stock=5

# Validation error: empty name → 303 with the error message
xh -F --form POST http://localhost:8080/products name= price=10 stock=1
```

## Tests

```sh
hurl --test test.hurl
```

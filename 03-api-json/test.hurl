# Fase 3: API JSON + middlewares

# The server is alive
GET http://localhost:8080/health
HTTP/1.1 200
[Asserts]
body == "OK\n"

# The page loads the client
GET http://localhost:8080/
HTTP/1.1 200
[Asserts]
header "Content-Type" == "text/html; charset=utf-8"
body contains "src=\"/static/js/app.js\""

# The store starts empty
GET http://localhost:8080/api/products
HTTP/1.1 200
[Asserts]
header "Content-Type" == "application/json; charset=utf-8"
body contains "[]"

# Post creates and answers 201 + Location
POST http://localhost:8080/api/products
{
    "name": "Teclado",
    "price": 50.5,
    "stock": 3
}
HTTP/1.1 201
[Asserts]
header "Location" == "/api/products/1"
jsonpath "$.id" == 1
jsonpath "$.name" == "Teclado"

# Get list 
GET http://localhost:8080/api/products
HTTP/1.1 200
[Asserts]
jsonpath "$[0].name" == "Teclado"

# Get one
GET http://localhost:8080/api/products/1
HTTP/1.1 200
[Asserts]
jsonpath "$.price" == 50.5
jsonpath "$.stock" == 3

# Bad id and missing product
GET http://localhost:8080/api/products/abc
HTTP/1.1 400
[Asserts]
jsonpath "$.error" == "id must be a positive integer"

GET http://localhost:8080/api/products/99
HTTP/1.1 404
[Asserts]
jsonpath "$.error" == "product not found"

# Patch changes only the sent fields
PATCH http://localhost:8080/api/products/1
{
    "stock": 10
}
HTTP/1.1 200
[Asserts]
jsonpath "$.stock" == 10
jsonpath "$.name" == "Teclado"

# Put replaces the whole product
PUT http://localhost:8080/api/products/1
{
    "name": "Teclado RGB",
    "price": 80,
    "stock": 5
}
HTTP/1.1 200
[Asserts]
jsonpath "$.name" == "Teclado RGB"
jsonpath "$.price" == 80

# Unknown fields are rejected
POST http://localhost:8080/api/products
{
    "nme": "Ratón"
}
HTTP/1.1 400
[Asserts]
jsonpath "$.error" contains "unknown field"

# Delete answers 204 and the product is gone
DELETE http://localhost:8080/api/products/1
HTTP/1.1 204

DELETE http://localhost:8080/api/products/1
HTTP/1.1 404
[Asserts]
jsonpath "$.error" == "product not found"

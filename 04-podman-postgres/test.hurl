# Fase 4: PostgreSQL + Podman

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

# Static files are served with the right content type
GET http://localhost:8080/static/js/app.js
HTTP/1.1 200
[Asserts]
header "Content-Type" == "text/javascript; charset=utf-8"

GET http://localhost:8080/static/css/style.css
HTTP/1.1 200
[Asserts]
header "Content-Type" == "text/css; charset=utf-8"

# The store starts empty (the Tests section truncates the table first)
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

# Get list and Get one
GET http://localhost:8080/api/products
HTTP/1.1 200
[Asserts]
jsonpath "$[0].name" == "Teclado"

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

# Empty patch is rejected
PATCH http://localhost:8080/api/products/1
{}
HTTP/1.1 400
[Asserts]
jsonpath "$.error" == "no fields to update (send name, price and/or stock)"

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

# Validations: name, price and stock
POST http://localhost:8080/api/products
{
    "name": "",
    "price": 10,
    "stock": 1
}
HTTP/1.1 400
[Asserts]
jsonpath "$.error" == "name is required"

POST http://localhost:8080/api/products
{
    "name": "Barato",
    "price": -1
}
HTTP/1.1 400
[Asserts]
jsonpath "$.error" == "price must be >= 0"

POST http://localhost:8080/api/products
{
    "name": "Sin stock",
    "price": 1,
    "stock": -1
}
HTTP/1.1 400
[Asserts]
jsonpath "$.error" == "stock must be >= 0"

# Delete answers 204 and the product is gone
DELETE http://localhost:8080/api/products/1
HTTP/1.1 204

DELETE http://localhost:8080/api/products/1
HTTP/1.1 404
[Asserts]
jsonpath "$.error" == "product not found"

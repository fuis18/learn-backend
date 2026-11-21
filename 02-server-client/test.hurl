# Phase 2: MVC with forms (GET + POST)

# The server is alive
GET http://localhost:8080/health
HTTP/1.1 200
[Asserts]
body == "OK\n"

# The main page renders the form and the empty table
GET http://localhost:8080/
HTTP/1.1 200
[Asserts]
header "Content-Type" == "text/html; charset=utf-8"
body contains "action=\"/products\""
body contains "No products"

# Create a product: real POST with form data, follow the PRG redirect
# and check that the final page shows the message and the new row
POST http://localhost:8080/products
[Form]
name: Keyboard
price: 50.5
stock: 3
[Options]
location: true
HTTP/1.1 200
[Asserts]
body contains "Product created"
body contains "Keyboard"
body contains "50.5"

# Empty name → the server answers with the PRG redirect (303 + Location),
# so a browser reload never repeats the POST
POST http://localhost:8080/products
[Form]
name:
price: 10
stock: 1
HTTP/1.1 303
[Asserts]
header "Location" == "/?error=The+name+is+required"

# Invalid price → follow the redirect and check the error message,
# and that nothing was saved (the table still only has the Keyboard)
POST http://localhost:8080/products
[Form]
name: Mouse
price: abc
stock: 1
[Options]
location: true
HTTP/1.1 200
[Asserts]
body contains "Invalid price"
body contains "Keyboard"
body not contains "Mouse"

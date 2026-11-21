# Test Home (GET /)
GET http://localhost:8080/
HTTP/1.1 200

# Test Health (GET /health)
GET http://localhost:8080/health
HTTP/1.1 200
[Asserts]
body contains "OK"

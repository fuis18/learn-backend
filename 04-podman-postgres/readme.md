# 04-podman-postgres

- **Development**: The Go server runs on the host (`go run .`) and connects to PostgreSQL at `localhost:5432` using the `DB_*` variables.
- **Stack**: Podman Compose spins up PostgreSQL and the containerized server; the server applies migrations upon startup, and the data persists in the `pgdata` volume.

## Comandos

```sh
# development: run the server on the host
go run .

# start everything
podman-compose up

# build the app
podman-compose build

# stop
podman-compose down
# stop + delete the data volume
podman-compose down -v

podman logs -f learn-backend-04-app

# Para entrar a psql:
podman exec -it learn-backend-04-postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"'
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

# Validations: empty name and empty PATCH → 400
xh POST http://localhost:8080/api/products name= price:=10 stock:=1
xh PATCH http://localhost:8080/api/products/1
```

## Tests

```sh
podman exec learn-backend-04-postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "TRUNCATE products RESTART IDENTITY";'
hurl --test test.hurl
```

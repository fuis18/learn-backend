# 05-react-compose

- **Development**: Vite runs on the host (`:5173`) and proxies `/api` to `http://localhost:8080`.
- **Stack**: Nginx serves the built `dist/` directory and acts as a reverse proxy for `/api/` to the `backend` service.

## Backend

```sh
cd backend

go mod init react-compose
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/joho/godotenv

go run .
```

## Frontend

```sh
cd frontend
pnpm init

pnpm add -E react react-dom react-router
pnpm add -E tailwindcss @tailwindcss/vite
pnpm add -E zustand zod react-hook-form @hookform/resolvers
pnpm add -E @tanstack/react-table tanstack/react-query

pnpm add -D -E typescript vite @types/react @types/react-dom @vitejs/plugin-react
pnpm add -D -E oxlint

pnpm dlx shadcn@latest init -t vite

pnpm dlx shadcn@latest add form input button label
pnpm dlx shadcn@latest add table pagination
pnpm dlx shadcn@latest add dialog button-group

pnpm dev
pnpm lint
```

## Commands

```sh
# start everything
podman-compose up

# build the app
podman-compose build

# stop
podman-compose down
# stop + delete the data volume
podman-compose down -v
```

## API

JSON only (via `/api` on the Nginx or Vite proxy):

| Método | Ruta                 | Descripción                                               |
| ------ | -------------------- | --------------------------------------------------------- |
| GET    | `/health`            | Healthcheck → `OK` (text/plain)                           |
| GET    | `/api/products`      | List all products (ordered by ID)                         |
| POST   | `/api/products`      | Create → `201` + `Location`                               |
| GET    | `/api/products/{id}` | Product detail                                            |
| PATCH  | `/api/products/{id}` | Update only the fields provided (`code`, `name`, `price`) |
| DELETE | `/api/products/{id}` | Delete → `204`                                            |

Errors: `400` (invalid JSON, validation issues, non-numeric ID, PATCH request with no fields), `404` (product not found), `409` (duplicate `code`), `500`.

## CLI Tests

```sh
# Read the SPA served by Nginx and the health check (direct: /health is not proxied)
xh GET http://localhost:3000/
xh GET http://localhost:8080/health

# See the products through the /api proxy
xh GET http://localhost:3000/api/products

# Create a product: JSON → 201 + Location (code is optional, the DB generates it)
xh POST http://localhost:3000/api/products name=Keyboard price:=50.5

# Get one and patch only the sent fields (there is no PUT)
xh GET http://localhost:3000/api/products/1
xh PATCH http://localhost:3000/api/products/1 price:=60

# Delete → 204
xh DELETE http://localhost:3000/api/products/1
```

## Tests

```sh
podman-compose down -v
hurl --test test.hurl
```

# 🗺️ Roadmap for Learning Backend Development with Go

## 🛠️ Testing and Development Tools

- **[Hurl](https://hurl.dev/):** For HTTP integration testing using plain-text files.
- **[xh](https://github.com/ducaale/xh):** A fast, user-friendly alternative to `curl` for quick testing.

## Stage 1 - Routing and Handlers

Building a simple server: a single binary that listens for requests and sends responses. The goal is to understand the process that occurs between the arrival of a request and the delivery of a response.

Folder: `01-server/`

- Server (Go)

![Stage 01](./assets/01-general.png)

## Stage 2 - MVC, HTML Forms, and Templates

The project grows beyond a single file: logic is organized into models, controllers, and views, while responses are structured using templates.

Folder: `02-server-client/`

- Server (Go)
- Client (views)

![Stage 02](./assets/02-general.png)

## Stage 3 - JSON API, Middlewares, and Consumption via `fetch`

The server moves beyond simply rendering pages; it exposes a JSON-based API using the five standard HTTP verbs, which the browser consumes via `fetch`.

Folder: `03-api-json/`

- Server (Go)
  - JSON (API consumption)
  - Middlewares (Logger, Recover)
- Client (views)

![Stage 03](./assets/03-general.png)


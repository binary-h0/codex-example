# codex-example

This repository contains a simple Go bulletin board server using the `gin` web framework and `gorm` for ORM. The database expects an H2 instance running in PostgreSQL compatibility mode.

## Running the server

1. Start an H2 database in PostgreSQL server mode (example port `5435`).
2. Set `DATABASE_URL` to the PostgreSQL connection string. For local testing the server defaults to:
   `host=localhost port=5435 user=sa dbname=test sslmode=disable`
3. Build and run the server:

```bash
go run .
```

Endpoints:
- `GET /posts` – list posts
- `GET /posts/:id` – get a single post
- `POST /posts` – create a post (expects JSON `{"title": "...", "content": "..."}`)



# Chirpy

Boot.dev backend path — [Learn HTTP Servers in Golang](https://www.boot.dev/courses/learn-http-servers-golang)

## Environment Variables

Create a `.env` file in the project root:

```
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
SECRET_KEY=<random-secret>
```

| Variable   | Description                          |
|------------|--------------------------------------|
| `DB_URL`   | PostgreSQL connection string         |
| `PLATFORM` | `dev` or `prod`                      |
| `SECRET_KEY` | JWT signing secret                 |

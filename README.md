# Learn Go

## Running the application

### Todo api

```
go run ./cmd/todo-api/main.go
```

#### Drop or Start with fresh DB

Start server with a fresh database

```
go run ./cmd/todo-api/main.go clean
```

Only drop the database

```
go run ./cmd/todo-api/main.go drop
```
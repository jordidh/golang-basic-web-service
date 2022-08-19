# golang-basic-web-service

Web service in GO with the most common features:

1. Accept command-line arguments
2. Print web server routes
3. Serve static content
4. Connect to Databases

## Project startup

```bash
mkdir golang-webserver
cd golang-webserver
go mod init golang-webserver
touch webserver.go
mkdir routes
touch routes/books.go
touch routes/authors.go
touch routes/publishers.go
mkdir static
touch static/error.txt
```

## Run in development

Commands to run de app

```bash
# Run all go modules (default env, port and graceful timeout)
go run .

# Run main program (default env, port and graceful timeout)
go run webserver.go

# Display command-line arguments help
go run webserver.go -help

# Run with specific env, port and graceful timeout
go run webserver.go -env production -port 3100 -graceful-timeout 15s

# Print all routes
go run webserver.go -print-routes
```

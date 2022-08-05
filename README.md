# golang-basic-web-service

Web service in GO with the most common features:

1. Accept command-line arguments
2. Print web server routes
3. Serve static content
4. Connect to Databases

## Run in development

Commands to run de app

```bash
# Run all go modules (default env, port and graceful timeout)
go run .

# Run main program (default env, port and graceful timeout)
go run basic-web-server.go

# Display command-line arguments help
go run basic-web-server.go -help

# Run with specific env, port and graceful timeout
go run basic-web-server.go -env production -port 3100 -graceful-timeout 15s

# Print all routes
go run basic-web-server.go -print-routes
```


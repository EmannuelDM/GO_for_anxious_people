# Go Hello World (01_hello_world)

A simple Hello World program written in Go.

## Installing Go

https://go.dev/doc/install

```bash
sudo apt  install golang-go
```

## Running the Program

To run the program, navigate to this directory (`01_hello_world`) and execute the following commands:

### Option 1: Run directly
```bash
go run main.go
```

### Option 2: Compile and run
1. Build the executable binary:
   ```bash
   go build main.go
   ```
2. Run the compiled binary:
   ```bash
   ./main
   ```


## what is the file go.mod?

In Go, a **`go.mod`** file is the root file of a **Go module**. It is the core of Go's built-in dependency management system.

When you create a `go.mod` file in a directory, that directory becomes the root of a Go module. The file primarily defines three things:

1. **The Module Path**: The import path prefix for all packages within the module (e.g., `module 01_hello_world`). This tells Go how other packages or modules can import code from this module.
2. **The Go Version**: The minimum version of Go required to run this module (e.g., `go 1.20`).
3. **Dependencies**: A list of external packages (and their specific versions) that your code depends on. These are added automatically when you run commands like `go get` or `go mod tidy`.

### Structure of a `go.mod` file:
```go
    module 01_hello_world // Module identity

    go 1.20              // Go version

    require (
        // External dependencies would appear here, for example:
        // github.com/gin-gonic/gin v1.9.0 
    )
```
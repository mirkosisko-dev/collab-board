# Makefile Targets Documentation

The `Makefile` provides a set of commands to streamline development tasks for this Go project. Below is a list of available targets with their descriptions.

## How to Use

To execute any of the targets, use the following command:

```bash
make [target]
```

For example, to build the project, run:

```bash
make build
```

## Available Targets

| Target  | Description                                |
| ------- | ------------------------------------------ |
| `help`  | Display this help message                  |
| `build` | Build the Go application                   |
| `run`   | Build and run the Go application           |
| `test`  | Run tests                                  |
| `clean` | Clean up generated files                   |
| `setup` | Set up the project (download dependencies) |
| `start` | Run the application                        |
| `env`   | Display current environment variables      |
| `lint`  | Run linting on the Go code                 |
| `fmt`   | Format the Go code                         |
| `check` | Run formatting, linting, and tests         |

## Setting Up the Project

To set up the project for the first time, run:

```bash
make setup
```

This command will download all Go dependencies.

Make sure to set `DATABASE_URL`:

```bash
export DATABASE_URL="your-database-url"
```

## Code Quality Checks

To run formatting, linting, and tests in one go, use:

```bash
make check
```

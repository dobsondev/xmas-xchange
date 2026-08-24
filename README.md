# Christmas Gift Exchange

A Go application for organizing Secret Santa / gift exchange events with support for participant restrictions.

## Features

- Read participant data from TOML configuration files
- Support for participant restrictions (people who shouldn't exchange gifts)
- Automatic gift exchange assignment with retry logic for constraint satisfaction
- Debug logging to trace the assignment process
- Configurable maximum retry attempts

## Building

```bash
make build
```

The binary will be created at `bin/xmas-xchange`.

## Running

The CLI is built with [cobra](https://github.com/spf13/cobra) and has two commands: `exchange` (compute a gift exchange and save it to a TOML file) and `print` (load and display a previously saved exchange).

```bash
# Compute an exchange, writing to a local file (--filename is required)
go run . exchange --toml=custom_participants.toml --filename=exchange-toml/my-exchange.toml

# Compute an exchange without dry-run mode (actually send messages)
go run . exchange --filename=exchange-toml/my-exchange.toml --dry-run=false

# Compute an exchange with a custom max attempts
go run . exchange --filename=exchange-toml/my-exchange.toml --max-attempts=100

# Write the exchange to S3 instead of the local filesystem
go run . exchange --filename=my-exchange.toml --output=s3 --s3-bucket=my-bucket --aws-region=us-east-1

# S3_BUCKET/AWS_REGION env vars work as a fallback for the flags above
S3_BUCKET=my-bucket AWS_REGION=us-east-1 go run . exchange --filename=my-exchange.toml --output=s3

# Print a previously saved local exchange
go run . print exchange-toml/my-exchange.toml

# Print a previously saved exchange from S3
go run . print s3://my-bucket/my-exchange.toml --aws-region=us-east-1
```

### Exchange TOML output

Every `exchange` run writes the computed giver/receiver assignments as TOML to the file you name with the required `--filename` flag. `--output` chooses the destination:

- `local` (default) — `--filename` is used as-is as a local file path (parent directories are created automatically).
- `s3` — `--filename` is used as the S3 object key. The bucket and region come from `--s3-bucket`/`--aws-region`, falling back to the `S3_BUCKET`/`AWS_REGION` environment variables if the flags aren't passed; one or the other is required.

`print` reads whatever location string `exchange` printed: a local path, or an `s3://bucket/key` URI (in which case `--aws-region`, falling back to `AWS_REGION`, is required).

## Configuration

Create a `participants.toml` file with your participant information:

```toml
[[participant]]
name = "Alice"
number = "+15551234567"
restrictions = ["Bob"]

[[participant]]
name = "Bob"
number = "+15559876543"
restrictions = ["Alice"]
```

See `participants.example.toml` for a complete example.

## Testing

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage summary
make test-coverage

# Generate HTML coverage report
make test-coverage-report
```

### Test Coverage

The project maintains high test coverage for business logic:

- **Exchange package**: 95.9% coverage - Core gift exchange algorithm
- **Participants package**: 63.6% coverage - TOML parsing and data structures
- **Main package**: 33.3% coverage - Entry point and configuration

Functions intentionally not covered by unit tests:
- `main()` and cobra `RunE` handlers - Entry points with I/O operations and flag parsing
- `setSlogLevel()` - Modifies global state
- `PrintParticipants()` and `PrintExchange()` - Console output formatting only

These functions are documented with comments explaining why they're excluded from test coverage.

## Environment Variables

- `LOG_LEVEL` - Set logging level (DEBUG, INFO, WARN, ERROR). Defaults to INFO.
- `S3_BUCKET` - Fallback for `exchange --s3-bucket` when writing to S3.
- `AWS_REGION` - Fallback for `exchange --aws-region` when writing to S3.

## Project Structure

```
.
├── main.go              # Entry point, delegates to cmd
├── cmd/                 # Cobra CLI commands
│   ├── root.go
│   ├── exchange.go
│   ├── print.go
│   └── root_test.go
├── exchange/            # Gift exchange logic
│   ├── exchange.go
│   ├── exchange_test.go
│   ├── exchange_toml.go
│   └── exchange_toml_test.go
├── exchange-toml/        # Suggested local exchange TOML output location (gitignored)
├── output/               # Local file / S3 writers for exchange TOML output
│   ├── local.go
│   ├── local_test.go
│   ├── s3.go
│   └── s3_test.go
└── participants/        # Participant data structures
    ├── participants.go
    └── participants_test.go
```

## How It Works

1. Participants are loaded from a TOML configuration file
2. The algorithm attempts to assign gift exchanges with these constraints:
   - No one can give to themselves
   - Restrictions are respected (configurable per participant)
   - Each person receives exactly one gift
   - Each person gives exactly one gift
3. If constraints can't be satisfied, the algorithm retries with a different random ordering
4. After max attempts are exhausted, the program panics with an error

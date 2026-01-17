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

```bash
# Run with default settings
make run

# Run with debug logging
make debug

# Run with custom participants file
go run . -toml=custom_participants.toml

# Run without dry-run mode (actually send messages)
go run . -dry-run=false

# Run with custom max attempts
go run . -max-attempts=100
```

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
- `main()` - Entry point with I/O operations and flag parsing
- `setSlogLevel()` - Modifies global state
- `PrintParticipants()` and `PrintExchange()` - Console output formatting only

These functions are documented with comments explaining why they're excluded from test coverage.

## Environment Variables

- `LOG_LEVEL` - Set logging level (DEBUG, INFO, WARN, ERROR). Defaults to INFO.

## Project Structure

```
.
├── main.go              # Entry point and configuration
├── exchange/            # Gift exchange logic
│   ├── exchange.go
│   └── exchange_test.go
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

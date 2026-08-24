# Christmas Gift Exchange

A Go application for organizing Secret Santa / gift exchange events with support for participant restrictions.

## Yearly Gift Exchange Run

This is the actual yearly process, run via GitHub Actions:

1. **Test run**: repo's **Actions** tab → **Run Gift Exchange** → **Run workflow**, leave `send-sms` unchecked (the default), and run it. This uploads a real exchange to S3 but sends no texts. Open the run's log and copy the filename from its `Exchange written to s3://...` line.

2. **Print the test results locally**, to sanity-check participants/restrictions:
   ```bash
   go run . print s3://<bucket>/<filename-from-step-1> --aws-region=<region>
   ```

3. **Dry-run the SMS messages locally**, to see exactly what each text will say (safe to do locally with real names — it's only your machine, never the Actions log):
   ```bash
   go run . sendsms s3://<bucket>/<filename-from-step-1> --aws-region=<region>
   ```

4. **Run it for real**: back in the **Actions** tab, run **Run Gift Exchange** again, this time *checking* `send-sms`. This computes a fresh exchange and actually texts everyone via Twilio.

> Every `exchange` run generates new random pairings, so the assignments you check in steps 2–3 are **not** the ones sent in step 4 — they're a different random draw. Steps 1–3 confirm the pipeline works end-to-end (S3 upload, AWS auth, message formatting), not a preview of the final assignments.

See [Local Quick Start](#local-quick-start) to run everything from your own machine instead, and [Running via GitHub Actions](#running-via-github-actions) for full workflow details.

## Local Quick Start

Running a full gift exchange is four steps:

```bash
# 1. Create a participants.toml file with everyone's name, number, and restrictions
cp participants.example.toml participants.toml
# ...then edit participants.toml with your real participants

# 2. Compute the exchange and save the assignments to a TOML file
go run . exchange --filename=exchange-toml/exchange.toml

# 3. Set the Twilio environment variables sendsms needs to actually send texts
export TWILIO_ACCOUNT_SID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export TWILIO_AUTH_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export TWILIO_PHONE_NUMBER=+15551234567

# 4. Text everyone their gift recipient
go run . sendsms exchange-toml/exchange.toml --dry-run=false
```

Step 3 needs to set **real environment variables** in the shell/process that runs step 4 — a `.env` file by itself does nothing here. This project doesn't load `.env` files automatically (no `godotenv` or similar); a `.env` is just a plain text file unless something sources it into the environment first (e.g. `set -a; source .env; set +a`, `direnv`, or your CI setting them as secrets). If you skip step 3 and run step 4 with `--dry-run=false`, `sendsms` fails fast with a clear error rather than attempting to send.

See [Configuration](#configuration) for the `participants.toml` format, [Running](#running) for the full set of flags (S3 output, max attempts, etc), and [Environment Variables](#environment-variables) for all of them.

### Troubleshooting

**Someone didn't get their text** (e.g. their carrier blocked the Twilio number): resend to just that one person with `--name`, instead of texting everyone again:

```bash
go run . sendsms exchange-toml/exchange.toml --dry-run=false --name="Alice"
```

**Still no luck / worst case**: print just their assignment locally and relay it to them yourself (in person, a different messaging app, etc.):

```bash
go run . print exchange-toml/exchange.toml --name="Alice"
```

**Checking a message's actual delivery status**: `sendsms` prints a Twilio message SID for every send (e.g. `Sent to Alice (+1 555 123 4567) [sid=SMxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx, status=queued]`), but that status is only the *initial* one — Twilio confirms real delivery asynchronously. Check the [Twilio Console](https://console.twilio.com) → Monitor → Logs → Messaging, or use the [Twilio CLI](https://www.twilio.com/docs/twilio-cli/quickstart):

```bash
# Install (macOS)
brew tap twilio/brew && brew install twilio

# Look up a specific message by SID
twilio api:core:messages:fetch --sid SMxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Or list recent messages without needing a SID
twilio api:core:messages:list --limit 20
```

The CLI picks up `TWILIO_ACCOUNT_SID`/`TWILIO_AUTH_TOKEN` from the environment automatically if they're already set (e.g. the same ones exported for `sendsms`) — otherwise run `twilio login` once to store a profile.

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

The CLI is built with [cobra](https://github.com/spf13/cobra) and has three commands: `exchange` (compute a gift exchange and save it to a TOML file), `print` (load and display a previously saved exchange), and `sendsms` (text each giver who their recipient is via Twilio).

```bash
# Compute an exchange, writing to a local file (--filename is required)
go run . exchange --toml=custom_participants.toml --filename=exchange-toml/my-exchange.toml

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

# Print just one participant's assignment (e.g. to relay it manually if their SMS was blocked)
go run . print exchange-toml/my-exchange.toml --name=Alice

# Simulate sending SMS notifications for everyone in a saved exchange (--dry-run defaults to true)
go run . sendsms exchange-toml/my-exchange.toml

# Actually send SMS notifications via Twilio (requires TWILIO_ACCOUNT_SID/TWILIO_AUTH_TOKEN/TWILIO_PHONE_NUMBER)
go run . sendsms exchange-toml/my-exchange.toml --dry-run=false

# Resend to just one participant by name (case-insensitive)
go run . sendsms exchange-toml/my-exchange.toml --dry-run=false --name=Alice

# sendsms also reads from S3 the same way print does
go run . sendsms s3://my-bucket/my-exchange.toml --aws-region=us-east-1
```

### Exchange TOML output

Every `exchange` run writes the computed giver/receiver assignments as TOML to the file you name with the required `--filename` flag. `--output` chooses the destination:

- `local` (default) — `--filename` is used as-is as a local file path (parent directories are created automatically).
- `s3` — `--filename` is used as the S3 object key. The bucket and region come from `--s3-bucket`/`--aws-region`, falling back to the `S3_BUCKET`/`AWS_REGION` environment variables if the flags aren't passed; one or the other is required.

`print` reads whatever location string `exchange` printed: a local path, or an `s3://bucket/key` URI (in which case `--aws-region`, falling back to `AWS_REGION`, is required). `--name` (case-insensitive) filters output to a single participant's assignment — handy if their SMS never arrived and you need to relay it another way.

### SMS notifications

`sendsms` reads a saved exchange (same local path / `s3://` input as `print`) and texts each giver: *"Hello \<name\>! Your gift recipient is \<recipient-name\>. Merry Christmas!"*

- `--dry-run` defaults to `true` — it prints what would be sent without calling Twilio, and doesn't require any `TWILIO_*` environment variables.
- `--dry-run=false` sends real messages and requires `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN`, and `TWILIO_PHONE_NUMBER` to all be set (Twilio credentials are environment-variable-only — there's no CLI flag for them, since they're secrets).
- `--name` limits the run to a single participant (case-insensitive match on their name), useful for resending to just one person.

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
- `AWS_REGION` - Fallback for `exchange --aws-region` / `print --aws-region` / `sendsms --aws-region` when using S3.
- `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN`, `TWILIO_PHONE_NUMBER` - Required by `sendsms --dry-run=false` to send real SMS messages via Twilio.

## Running via GitHub Actions

The `Run Gift Exchange` workflow (`.github/workflows/run-gift-exchange.yml`) is manually triggered (`workflow_dispatch`) and always computes an exchange and uploads it to S3, and — only if `send-sms` is checked — sends the SMS notifications too. Its `sendsms` step always runs with `--quiet`, so no participant names or numbers ever appear in the Actions log; leaving `send-sms` unchecked only creates the S3 file.

It takes two inputs when triggered:
- `send-sms` (boolean, default `false`) — check this to actually text participants; leave it unchecked for a test run.
- `filename` (string, optional) — the S3 key to save the exchange as. Leave it blank to default to `<current-year>_xchange.toml` (e.g. `2026_xchange.toml`), computed at run time — GitHub Actions doesn't support dynamic expressions in a `workflow_dispatch` input's declared default, so the input itself just shows blank in the trigger form. The `nosms_` prefix is only added to this auto-generated default when `send-sms` is unchecked; if you type a filename yourself, it's used exactly as given either way.

Before it can run, configure these in the repo's **Settings → Secrets and variables → Actions** (all as **Variables**, except the Twilio/participant values which are **Secrets** — none of this is hardcoded in the workflow, so forking the repo makes it obvious what needs to be set up):

- **Secrets**: `PARTICIPANTS_TOML` (the full contents of your `participants.toml`, pasted as-is — GitHub secrets support multi-line values, no encoding needed), `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN`, `TWILIO_PHONE_NUMBER`.
- **Variables**: `S3_BUCKET`, `AWS_ROLE_ARN` (the IAM role to assume via OIDC — must trust this repo), `AWS_REGION`.

## Project Structure

```
.
├── main.go              # Entry point, delegates to cmd
├── cmd/                 # Cobra CLI commands
│   ├── root.go
│   ├── exchange.go
│   ├── print.go
│   ├── sendsms.go
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
├── sms/                  # Twilio SMS sending
│   ├── twilio.go
│   └── twilio_test.go
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

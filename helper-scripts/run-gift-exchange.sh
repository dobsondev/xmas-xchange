#!/usr/bin/env bash
set -euo pipefail

# Runs a full gift exchange from your local machine: compute -> upload to S3 -> confirm ->
# send SMS. This is the local equivalent of triggering the "Run Gift Exchange" GitHub
# Actions workflow with send-sms checked. Assignments are never printed here so the
# organizer doesn't spoil their own surprise - use `print --name=<person>` later if
# someone's SMS didn't arrive and you need to relay their assignment manually.
#
# Requires: go, make, the 1Password CLI (op) signed in, and the "Twilio" item in the
# "Private" vault with "Account SID" / "Auth Token" / "Gift Exchange Phone Number" fields.

S3_BUCKET="${S3_BUCKET:-dobsondev-family-xmas-xchange}"
AWS_REGION="${AWS_REGION:-ca-central-1}"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${REPO_ROOT}"

FILENAME="$(date -u +%Y)_xchange.toml"
LOCATION="s3://${S3_BUCKET}/${FILENAME}"

echo "==> Checking AWS access to s3://${S3_BUCKET}..."
if ! aws s3 ls "s3://${S3_BUCKET}" --region "${AWS_REGION}" >/dev/null 2>&1; then
  echo "Can't reach the bucket - running 'aws sso login'..."
  aws sso login
  if ! aws s3 ls "s3://${S3_BUCKET}" --region "${AWS_REGION}" >/dev/null 2>&1; then
    echo "Still can't reach s3://${S3_BUCKET} after logging in. Aborting." >&2
    exit 1
  fi
fi

echo "==> Pulling Twilio credentials from 1Password..."
export TWILIO_ACCOUNT_SID
export TWILIO_AUTH_TOKEN
export TWILIO_PHONE_NUMBER
TWILIO_ACCOUNT_SID="$(op read 'op://Private/Twilio/Account SID')"
TWILIO_AUTH_TOKEN="$(op read 'op://Private/Twilio/Auth Token')"
TWILIO_PHONE_NUMBER="$(op read 'op://Private/Twilio/Gift Exchange Phone Number')"

echo "==> Building..."
make build

echo "==> Computing exchange and uploading to ${LOCATION}..."
./bin/xmas-xchange exchange --filename="${FILENAME}" --output=s3 --s3-bucket="${S3_BUCKET}" --aws-region="${AWS_REGION}"

echo
read -r -p "Exchange computed and uploaded to ${LOCATION}. Send SMS notifications to everyone now? [y/N] " CONFIRM
if [[ ! "${CONFIRM}" =~ ^[Yy]$ ]]; then
  echo "Aborted. No messages sent. Exchange is saved at ${LOCATION} if you want to send later."
  exit 0
fi

echo "==> Sending SMS notifications..."
./bin/xmas-xchange sendsms "${LOCATION}" --aws-region="${AWS_REGION}" --dry-run=false

echo "==> Done!"

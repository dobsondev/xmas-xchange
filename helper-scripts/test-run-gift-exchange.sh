#!/usr/bin/env bash
set -euo pipefail

# Test version of run-gift-exchange.sh: computes a real exchange and uploads it to S3
# (using a test_-prefixed, timestamped filename so it never collides with the real
# yearly file), prints the assignments, then optionally sends ONE real test SMS to a
# single participant you choose - never to everyone. Useful for sanity-checking the
# whole pipeline (including a real Twilio send) without risking texting the full list.
# The test send always uses --test, so the message is prefixed with "TEST: ".
#
# Requires: go, make, the 1Password CLI (op) signed in, and the "Twilio" item in the
# "Private" vault with "Account SID" / "Auth Token" / "Gift Exchange Phone Number" fields.

S3_BUCKET="${S3_BUCKET:-dobsondev-family-xmas-xchange}"
AWS_REGION="${AWS_REGION:-ca-central-1}"

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${REPO_ROOT}"

FILENAME="test_$(date -u +%Y%m%d%H%M%S)_xchange.toml"
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
echo "==> Assignments:"
./bin/xmas-xchange print "${LOCATION}" --aws-region="${AWS_REGION}"

echo
read -r -p "Send a real test SMS to one participant? [y/N] " SEND_TEST
if [[ ! "${SEND_TEST}" =~ ^[Yy]$ ]]; then
  echo "No test SMS sent. Exchange is saved at ${LOCATION}."
  exit 0
fi

read -r -p "Participant name to send the test SMS to: " PARTICIPANT_NAME
if [[ -z "${PARTICIPANT_NAME}" ]]; then
  echo "No name entered, aborting."
  exit 1
fi

echo "==> Sending test SMS to ${PARTICIPANT_NAME}..."
./bin/xmas-xchange sendsms "${LOCATION}" --aws-region="${AWS_REGION}" --dry-run=false --name="${PARTICIPANT_NAME}" --test

echo "==> Done!"

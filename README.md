# Christmas Gift Exchange Script

[![Test Run Gift Exchange](https://github.com/dobsondev/xmas-xchange/actions/workflows/test-build.yml/badge.svg)](https://github.com/dobsondev/xmas-xchange/actions/workflows/test-build.yml)

This project contains a containerized Python script that will create a Christmas gift exchange and send SMS messages (via [Twilio](https://www.twilio.com/en-us)) to the participants notifying them of who they are buying gifts for. This way no one has to be the "keeper of the secrets" and know who got assigned to who, rather it will be a surprise for everyone involved in the gift exchange including the organizer.

The script also has built in options for constraints so that certain people cannot be assigned other participants in the gift exchange. This is useful to ensure that couples are not paired together or other similar situations.

The results of the gift exchange will also be uploaded to [AWS S3](https://aws.amazon.com/s3/) storage so that if needed you can, as the organizer of the gift exchange, double check who was assigned to who and solve any issues that may arise.

> **🎄 For Annual Use:** 99% of the time you'll just use the GitHub Actions workflows. The local Docker commands are mainly for development or troubleshooting. If anything you will run the local test workflow to ensure your `json/data.json` is valid locally before pushing it up to GitHub.

## 🎯 Quick Start (Recommended - GitHub Actions)

For annual gift exchange - this is what you'll use most:

### Step 1: Test Run (Always do this first!)
1. Go to **Actions** tab in your GitHub repository
2. Click **Test Run Gift Exchange** workflow
3. Click **Run workflow** → Run on `main` branch
4. Wait for completion and check that test passes
5. Review the test output to verify everything looks correct

### Step 2: Real Run (After test passes)
1. Go to **Actions** tab in your GitHub repository  
2. Click **Run Gift Exchange** workflow
3. Click **Run workflow** → Run on `main` branch
4. SMS messages will be sent to participants! 🎉

**⚠️ Important:** Always run the test first each year to verify everything works!

## 🔄 Workflow Overview

```
Update Participants → Test Run → Verify Results → Real Run → SMS Sent!
                          ↓              ↓
                    (GitHub Action)  (Check S3 file)
```

## Prerequisites

**Required Services:**
- **Twilio:** SMS service (~$0.10 per message sent)
- **AWS S3:** File storage (~$0.01/month for small files)

**Setup Time:** ~30 minutes first time, ~5 minutes annually

You will need the following accounts in order to get this project to work:

1. Twilio account with a number capable of sending SMS
2. Amazon Web Services account with S3 storage and IAM user setup to interact with S3

## 📋 Annual Checklist

Before running each year:

- [ ] Update participant list in `json/data.json` (locally)
- [ ] Update GitHub secret `DATA_JSON` with new participant data
- [ ] Run **Test Run Gift Exchange** workflow first
- [ ] Verify test results look correct
- [ ] Run **Run Gift Exchange** workflow for real SMS sending

**Pro tip:** Keep your local `json/data.json` updated year-round, then just update the GitHub secret when ready to run.

## GitHub Setup

### Setting Up GitHub Secrets

You need two secrets in your GitHub repository:

**`ENV` Secret:**
1. Copy your working local `.env` file content (see Local Setup section for details)
2. Go to your repository → Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Name: `ENV`
5. Paste your `.env` file content directly

**`DATA_JSON` Secret:**
1. Run: `cat json/data.json | base64`
2. Copy the output
3. Go to your repository → Settings → Secrets and variables → Actions
4. Click "New repository secret" 
5. Name: `DATA_JSON`
6. Paste the base64 output

**💡 Tip:** You only need to update `DATA_JSON` when participants change.

### Branch Protection (Optional)

While not strictly required, I would recommend setting up branch protection for your `main` branch. I like to include the following settings:

- Require a pull request before merging
- Require status checks to pass before merging
  - Require branches to be up to date before merging
  - `Test run xmas-change.py` status check is required

### Test Run Gift Exchange on GitHub

The **Test Run Gift Exchange** (`.github/test-build.yml`) workflow builds the Docker image and runs a test where no SMS messages are sent. You can think of this like the "UAT" run of the script.

The workflow will create a file prefixed with `github_` and suffixed with `_dryrun` that will be uploaded to S3 so that you can inspect the results. This allows you to ensure everything is working as intended for both the constraints and participant phone numbers. It will NOT send SMS messages to participants.

### Run Gift Exchange on GitHub

The **Run Gift Exchange** (`.github/run-script.yml`) workflow builds the Docker image and runs the script with the intention of sending SMS messages to the participants for the real gift exchange. You can think of this like the "production" run of the script.

The workflow will create a file with the gift exchange results and upload that to S3 in case you need to verify what participant got assigned to what other participant. The idea is this file will not be viewed unless needed so that no one knows who was assigned each other. No prefixes or suffixes will be added to the upload. This run will send out SMS messages to the participants to let them know who they were assigned in the gift exchange.

## Local Setup

There are two files you need to setup in order for this project to work on your local machine:

1. `.env`
2. `json/data.json`

### `.env`

The `.env` file needs to contain all your credentials for Twilio and AWS. See `.env.example` to see what the variables should be setup. You will need to provide:

1. Twilio Account SID
2. Twilio Auth Token
3. Twilio Phone number in the format of `+1##########`
4. AWS IAM User Access Key
5. AWS IAM User Secret Access Key
6. AWS Region of S3 Bucket
7. S3 Bucket Name

To get your AWS credentials, you'll need to create an IAM user and access keys in the AWS console:

1. Go to IAM Service:
  - Log into AWS Console
  - Search for "IAM" or go to https://console.aws.amazon.com/iam/
2. Create IAM User:
  - Click "Users" in left sidebar
  - Click "Create user"
  - Enter a username
  - Click "Next"
3. Set Permissions:
  - Choose "Attach policies directly"
  - Select the AWS services your app needs (either S3Full or only access to the bucket you want to use)
  - Click "Next" then "Create user"
4. Create Access Keys:
  - Click on your new user
  - Go to "Security credentials" tab
  - Click "Create access key"
  - Choose "Application running outside AWS"
  - Click "Next" then "Create access key"
5. Copy Credentials:
  - Copy the Access key ID (this is your `AWS_ACCESS_KEY_ID`)
  - Copy the Secret access key (this is your `AWS_SECRET_ACCESS_KEY`)
  - Store them securely in your password vault - you can't view the secret key again

### `./json/data.json`

This file needs to contain all the required information of the participants in the gift exchange. This should include their names, their mobile phone numbers (in the format of `+1##########`) and any constraints (people they cannot match with for the gift exchange). See `./json/data.example.json` for an example of how this should be formatted.

A participants entry should look like this:

```json
"Participant": {
    "phone_number": "+15556667777",
    "constraints": ["Wife", "Brother"]
}
```

In this example, "Participant" is not allowed to match with "Wife" or "Brother" who would also have their own entries in the file. This way you can make sure that participants aren't matched up with their partners or whatever other constraints you might choose to have. Note that in this example I am using "Participant", "Wife" and "Brother" as sample names just to make it clear what each person is in relation to each other.

## Run Locally with Docker

Below you will find a summary of helpful Docker commands that might need to be used for this project.

### Build the Docker Image

```bash
docker build -t xmas-xchange .
```

To build and ensure there is no caching, use:

```bash
docker build --no-cache -t xmas-xchange .
```

### Run the Docker Container

```bash
docker run --env-file .env --rm xmas-xchange
```

### Dry Run the Docker Container

```bash
docker run --env-file .env --rm xmas-xchange --dry-run
```

If you want to hide sensitive output (names and phone numbers), then use the following (this is used in the GitHub actions workflow to ensure nothing sensitive gets posted on GitHub.com):

```bash
docker run --env-file .env --rm xmas-xchange --dry-run --hide-sensitive-output
```

There is also an option specifically for when the script is run on a GitHub runner for testing:

```bash
docker run --env-file .env --rm xmas-xchange --github-test
```

This is the equivalent of running `docker run --env-file .env --rm xmas-xchange --dry-run --hide-sensitive-output` and does a little extra output formatting to make it clear it's running on GitHub.

### Using the `helper.py` Script

I've added a helper script to retrieve the gift giver and recipient based on the S3 file name as well as the gift givers name. This can be used in case someone's carrier blocks the SMS message or something to that effect.

```bash
docker run --env-file .env --rm --entrypoint python xmas-xchange helper.py "<S3_FILE_NAME>" "<GIFT_GIVER_NAME>"
```

The output from this should look like:

```bash
✅ S3 connection successful!
Adam -> Beatrice
```

### Using the `test.py` Script

The comprehensive test script validates your entire gift exchange workflow and provides confidence that everything works correctly before sending real SMS messages.

#### What it Tests

1. **Service Connections** - Validates both S3 and Twilio connectivity
2. **Constraint Display** - Shows all constraints for transparency  
3. **Dry-Run Execution** - Runs the main script and captures the S3 filename
4. **Helper Script Validation** - Tests helper queries for **all participants**
5. **Assignment Download** - Downloads and parses the full assignment from S3
6. **Constraint Validation** - Ensures no constraint violations occurred
7. **Completeness Check** - Verifies everyone gives and receives exactly once
8. **Cross-Validation** - Confirms helper results match the full assignment perfectly

#### Usage

```bash
docker run --env-file .env --rm --entrypoint python xmas-xchange test.py
```

#### Sample Output
```
=== Gift Exchange Workflow Test ===

Loaded 7 people with constraints
Constraints:
  Adam cannot give to: Beatrice
  Beatrice cannot give to: Adam
  Carole cannot give to: Danielle
  Danielle cannot give to: Carole
  Edgar cannot give to: Frank
  Frank cannot give to: Edgar, Danielle
  Gray has no constraints

Testing service connections...
✅ S3 connection successful!
✅ Twilio connection successful!
✅ All service connections successful

Running dry-run...
✅ Dry-run completed successfully

Testing helper queries for all 7 people
  ✅ Adam -> Danielle
  ✅ Beatrice -> Edgar
  ✅ Carole -> Adam
  ✅ Danielle -> Beatrice
  ✅ Edgar -> Gray
  ✅ Frank -> Carole
  ✅ Gray -> Frank

Downloading and validating full assignment...
✅ Downloaded assignment file with 7 assignments

Validating constraints...
✅ All constraints satisfied
Validating assignment completeness...
✅ Assignment is complete and valid

Cross-validating all 7 helper results...
  ✅ Adam -> Danielle (matches)
  ✅ Beatrice -> Edgar (matches)
  ✅ Carole -> Adam (matches)
  ✅ Danielle -> Beatrice (matches)
  ✅ Edgar -> Gray (matches)
  ✅ Frank -> Carole (matches)
  ✅ Gray -> Frank (matches)

🎉 ALL TESTS PASSED! 🎉
Generated assignment file: 2025-08-15_20-29-02_gift_assignments_dryrun.txt
```

This comprehensive testing gives you full confidence in your gift exchange system before the real run!

## 🚨 Troubleshooting

**Test workflow fails?**
- Check GitHub secrets are set correctly
- Verify Twilio/AWS credentials haven't expired
- Make sure your S3 bucket exists and is accessible

**Participant doesn't receive SMS?**
- Use helper script: `docker run --env-file .env --rm --entrypoint python xmas-xchange helper.py "<filename>" "<name>"`
- Check phone number format in `json/data.json` (must be `+1##########`)
- Verify the participant's phone can receive SMS from Twilio numbers

**Need to check who got who?**
- Download assignment file from S3 bucket
- Only view if absolutely necessary to preserve surprise!

**GitHub Action fails with secrets error?**
- Verify both `ENV` and `DATA_JSON` secrets are set in repository settings
- For `DATA_JSON`: ensure you used `base64` encoding: `cat json/data.json | base64`

**Local setup not working?**
- Verify `.env` file exists and has all required variables
- Check `json/data.json` exists and follows the correct format
- Test Docker connectivity: `docker run hello-world`
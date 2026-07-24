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
- **AWS IAM:** OIDC authentication for GitHub Actions (no cost)

**Setup Time:** ~30 minutes first time, ~5 minutes annually

You will need the following accounts in order to get this project to work:

1. Twilio account with a number capable of sending SMS
2. Amazon Web Services account with S3 storage and IAM configured for GitHub Actions OIDC

## 📋 Annual Checklist

Before running each year:

- [ ] Update participant list in `json/data.json` (locally)
- [ ] Update GitHub secret `DATA_JSON` with new participant data
- [ ] Run **Test Run Gift Exchange** workflow first
- [ ] Verify test results look correct
- [ ] Run **Run Gift Exchange** workflow for real SMS sending

**Pro tip:** Keep your local `json/data.json` updated year-round, then just update the GitHub secret when ready to run.

## AWS Setup (One-Time)

The AWS infrastructure for this project is managed with **Terraform** (see `infrastructure/`) instead of manual console setup. It provisions:

- An **S3 bucket** for storing gift exchange results (encrypted, versioning disabled, public access blocked)
- An **IAM role** (`github-actions-xmas-xchange-role`) that GitHub Actions assumes via **OpenID Connect (OIDC)** for secure, temporary AWS credentials, scoped to `s3:PutObject` / `s3:GetObject` / `s3:ListBucket` on that bucket
- An **IAM user** (`dobsondev-family-xmas-xchange`) with S3 access, used for [local development](#for-local-development-creating-aws-iam-user)

### Prerequisites

- [Terraform](https://developer.hashicorp.com/terraform/install) installed locally
- AWS credentials configured locally (e.g. via `aws configure` or environment variables) with permission to manage IAM and S3
- A GitHub OIDC identity provider already registered in IAM. This is an account-level resource (`https://token.actions.githubusercontent.com`, audience `sts.amazonaws.com`) and is **not** created by `main.tf`, since most accounts only ever need one. If your account doesn't have one yet, create it once via **IAM Console** → **Identity providers** → **Add provider** → **OpenID Connect**.

### Step 1: Review Configuration

`infrastructure/main.tf` hardcodes values specific to this deployment (AWS account ID, GitHub repo `dobsondev/xmas-xchange`, S3 bucket name `dobsondev-family-xmas-xchange`). If you're forking this project for your own AWS account, update these values before applying.

### Step 2: Configure the Backend

Terraform state is stored remotely in S3.

1. Copy `infrastructure/backend.tfvars.example` to `infrastructure/backend.tfvars`
2. Update `bucket`, `key`, and `region` to point at your own state bucket

`backend.tfvars` is gitignored since state backend configuration is environment-specific.

### Step 3: Apply

```bash
cd infrastructure
terraform init -backend-config=backend.tfvars
terraform plan
terraform apply
```

**✅ AWS Setup Complete!** GitHub Actions will now use temporary credentials via OIDC, and the IAM user for local development will be ready for you to generate access keys for (see below).

## GitHub Setup

### Setting Up GitHub Secrets

You need the following secrets in your GitHub repository:

**Twilio Secrets:**
1. Go to your repository → Settings → Secrets and variables → Actions
2. Create these secrets:
   - `TWILIO_ACCOUNT_SID` - Your Twilio Account SID
   - `TWILIO_AUTH_TOKEN` - Your Twilio Auth Token
   - `TWILIO_PHONE_NUMBER` - Your Twilio phone number (format: `+15556667777`)

**S3 Secret:**
1. `S3_BUCKET` - Your S3 bucket name

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
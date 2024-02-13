# Christmas Gift Exchange Script

[![Test Run Gift Exchange](https://github.com/dobsondev/xmas-xchange/actions/workflows/test-build.yml/badge.svg)](https://github.com/dobsondev/xmas-xchange/actions/workflows/test-build.yml)

This project contains a containerized Python script that will create a Christmas gift exchange and send SMS messages (via [Twilio](https://www.twilio.com/en-us)) to the participants notifying them of who they are buying gifts for. This way no one has to be the "keeper of the secrets" and know who got assigned to who, rather it will be a surprise for everyone involved in the gift exchange including the organizer.

The script also has built in options for constraints so that certain people cannot be assigned other participants in the gift exchange. This is useful to ensure that couples are not paired together or other similar situations.

The results of the gift exchange will also be uploaded to [AWS S3](https://aws.amazon.com/s3/) storage so that if needed you can, as the organizer of the gift exchange, double check who was assigned to who and solve any issues that may arise.

There are two ways to setup and run this project:

1. Run it on your local machine via Docker
2. Run it on GitHub Actions via Docker

## Prerequisites

You will need the following accounts in order to get this project to work:

1. Twilio account with a number capable of sending SMS
2. Amazon Web Services account with S3 storage and IAM user setup to interact with S3

## Local Setup

There are two files you need to setup in order for this project to work on your local machine:

1. `.env`
2. `json/data.json`

### `.env`

The `.env` file needs to contain all your credentials for Twilio and AWS. See `.env.example` to see what the variables should be setup. You will need to provide:

1. Twilio Account SID
2. Twilo Auth Token
3. Twilo Phone number in the format of `+1##########`
4. AWS IAM User Access Key
5. AWS IAM User Secret Access Key
6. AWS Region of S3 Bucket
7. S3 Bucket Name

### `./json/data.json`

This file needs to contain all the required information of the participants in the gift exchange. This should include their names, their mobile phone numbers (in the format of `+1##########`) and any constraints (people they cannot match with for the gift exchange). See `./json/data.example.json` for an example of this should be formatted.

A participants entry should look like this:

```json
"Participant": {
    "phone_number": "+15556667777",
    "constraints": ["Wife", "Brother"]
}
```

In this example, "Participant" is not allowed to match with "Wife" or "Brother" who would also have their own entries in the file. This way you can make sure that participants aren't matched up with their partners or whatever other constraints you might choose to have. Note that in this example I am using "Participant", "Wife" and "Brother" as sample names just to make it clear what each person is in relation to each other.

### GitHub Actions Setup

In order to setup your repository so that you can use the GitHub actions in this repository, you will need to add two secrets:

1. `ENV`
2. `DATA_JSON`

If you do not want to run this project using GitHub Actions - simply remove the `.github` folder from your fork of the project. Then you do not have to setup any secrets and the actions will not run.

### Branch Protection

While not strictly required, I would recommend setting up branch protection for your `main` branch. I like to include the following settings:

- Require a pull request before merging
- Require status checks to pass before merging
  - Require branches to be up to date before merging
  - `Test run xmas-change.py` status check is required

### `ENV` Secret

These map directly to the above `.env` file and `json/data.json` files which are used to run this project on your local.

The `ENV` secret is a direct copy paste of your complete `.env` file. Once you have a working `.env` on your local, run the following command to copy the file's content to your clipboard:

```bash
pbcopy < .env
```

### `DATA_JSON` Secret

You can then paste the `.env` file's content into your repository as a secret named `ENV`.

For the `DATA_JSON`, you must first base64 encode the data of the file and then upload that as your secret. Once you have a working `data/data.json` on your local, run the following command to base64 encode it's content:

```bash
cat json/data.json | base64
```

Copy the output of this command and then paste it into your repository as a secret named `DATA_JSON`.

### Test Run Gift Exchange

Once you have your secrets setup - everything should work with the GitHub workflows in this project. The **Test Run Gift Exchange** (`.github/test-build.yml`) workflow build the Docker image and a test run will be ran on the script where no SMS messages will be set. You can think of this like the "UAT" run of the script.

The **Test Run Gift Exchange** (`.github/test-build.yml`) workflow run will create a file prefixed with `github_` and suffixed with `_dryrun` that will be uploaded to S3 so that you can inspect the results of the run. This allows you to ensure everything is working as intended for both the constraints and participant phone numbers. It will NOT send SMS messages to participants.

### Run Gift Exchange

The **Run Gift Exchange** (`.github/run-script.yml`) workflow will build the Docker image and then run the script with the intention of sending SMS messages to the participants for the real gift exchange. You can think of this like the "production" run of the script.

The **Run Gift Exchange** (`.github/run-script.yml`) workflow run will create a file with the gift exchange results and upload that to S3 in case you need to verify what participant got assigned to what other participant. The idea is this file will not be viewed unless needed so that no one knows who was assigned each other. No prefixes or suffixes will be added to the upload. This run will send out SMS messages to the participants to let them know who they were assigned in the gift exchange.

## Docker

Below you will find a summary of helpful Docker commands that might need to be used for this project.

### Build the Docker Image

```bash
docker build -t xmas-xchange .
```

### Run the Docker Container

```bash
docker run --rm xmas-xchange
```

### Dry Run the Docker Container

```bash
docker run --rm xmas-xchange --dry-run
```

If you want to hide sensitive output (names and phone numbers), then use the following (this is used in the GitHub actions workflow to ensure nothing sensitive gets posted on GitHub.com):

```bash
docker run --rm xmas-xchange --dry-run --hide-sensitive-output
```

There is also an option specifically for when the script is run on a GitHub runner for testing:

```bash
docker run --rm xmas-xchange --github-test
```

This is the equivalent of running `docker run --rm xmas-xchange --dry-run --hide-sensitive-output` and does a little extra output formatting to make it clear it's running on GitHub.
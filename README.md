# Christmas Gift Exchange Script

This project contains a containerized Python script that will create a Christmas gift exchange and send txt messages to the participants notifying them of who they are buying gifts for.

## Prerequisites

You will need the following in order to get this project to work:

1. Twilio account with a number capable of sending SMS
2. Amazon Web Services account with S3 storage and IAM user setup to interact with S3

## Setup

There are two files you need to setup in order for this project to work:

1. `.env`
2. `/json/data.json`

### `.env`

This needs to contain all your credentials for Twilio and AWS. See `.env.example` to see what the envionrment variable names should be. You will need to provide:

1. Twilio Account SID
2. Twilo Auth Token
3. Twilo Phone number in the format of `+1##########`
4. AWS IAM User Access Key
5. AWS IAM User Secret Access Key
6. AWS Region of S3 Bucket

### `./json/data.json`

This file needs to contain all the required information of the participants in the gift exchange. This should include their names, their mobile phone numbers (in the format of `+1##########`) and any constraints (people they cannot match with for the gift exchange). See `./json/data.example.json` for an example of this should be formatted.

A participants entry should look like this:

```json
"Participant": {
    "phone_number": "+15556667777",
    "constraints": ["Wife", "Brother"]
}
```

In this example, "Participant" is not allowed to match with "Wife" or "Brother" who would also have entries in the file. This way you can make sure that participants aren't matched up with their partners or whatever other constraints you might choose to have.

### GitHub Action Secrets Setup

In order to setup your repository so that you can use the GitHub actions in this repository, you will need to add two secrets:

1. `ENV`
2. `DATA_JSON`

The `ENV` secret is a simple copy paste of your complete `.env` file. If you have a working `.env` on your local, copy it up to GitHub as a secret and it will work there.

For the `DATA_JSON`, you must first base64 encode the data of the file and put that up as a secret. To do this, run the following command when you have everything working as expected on your local:

```bash
cat json/data.json | base64
```

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
bash
```
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
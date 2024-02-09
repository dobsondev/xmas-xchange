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

## Docker

Below you will find a summary of helpful Docker commands that might need to be used for this project.

### Build the Docker Image

```
docker build -t xmas-xchange .
```

### Run the Docker Container

```
docker run --rm xmas-xchange
```

### Dry Run the Docker Container

```
docker run --rm xmas-xchange --dry-run
```

If you want to hide sensitive output (names and phone numbers), then use the following (this is used in the GitHub actions workflow to ensure nothing sensitive gets posted on GitHub.com):

```
docker run --rm xmas-xchange --dry-run --hide-sensitive-output
```
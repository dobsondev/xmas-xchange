# Local Development & Setup

There are two files you need to setup in order for this project to work on your local machine:

1. `.env`
2. `json/data.json`

### `.env`

The `.env` file needs to contain all your credentials for Twilio and AWS. See `.env.example` to see what the variables should be setup. You will need to provide:

1. Twilio Account SID
2. Twilio Auth Token
3. Twilio Phone number in the format of `+1##########`
4. AWS IAM User Access Key (for local development only)
5. AWS IAM User Secret Access Key (for local development only)
6. AWS Session Token (optional - only needed for temporary credentials)
7. AWS Region of S3 Bucket
8. S3 Bucket Name

#### For Local Development: Creating AWS IAM User

The IAM user for local development (`dobsondev-family-xmas-xchange`) is created by [Terraform](#aws-setup-one-time) as part of the AWS Setup above, so you only need to generate access keys for it:

1. Go to IAM Service:
   - Log into AWS Console
   - Search for "IAM" or go to https://console.aws.amazon.com/iam/
2. Find the User:
   - Click "Users" in left sidebar
   - Click on `dobsondev-family-xmas-xchange`
3. Create Access Keys:
   - Go to "Security credentials" tab
   - Click "Create access key"
   - Choose "Application running outside AWS"
   - Click "Next" then "Create access key"
4. Copy Credentials:
   - Copy the Access key ID (this is your `AWS_ACCESS_KEY_ID`)
   - Copy the Secret access key (this is your `AWS_SECRET_ACCESS_KEY`)
   - Store them securely in your password vault - you can't view the secret key again

**Note:** GitHub Actions uses OIDC and doesn't need these access keys. The access keys are only for local development.

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
- Verify Twilio credentials haven't expired
- Ensure AWS OIDC role exists and has correct trust policy
- Make sure your S3 bucket exists and is accessible
- Verify the IAM role has the correct S3 permissions

**AWS OIDC authentication fails?**
- Verify OIDC identity provider exists in IAM
- Check the trust policy on your GitHub Actions role
- Ensure the `repo:` condition matches your repository exactly
- Confirm the role has the inline S3 policy attached

**Participant doesn't receive SMS?**
- Use helper script: `docker run --env-file .env --rm --entrypoint python xmas-xchange helper.py "<filename>" "<name>"`
- Check phone number format in `json/data.json` (must be `+1##########`)
- Verify the participant's phone can receive SMS from Twilio numbers

**Need to check who got who?**
- Download assignment file from S3 bucket
- Only view if absolutely necessary to preserve surprise!

**GitHub Action fails with secrets error?**
- Verify all Twilio and S3 secrets are set in repository settings
- For `DATA_JSON`: ensure you used `base64` encoding: `cat json/data.json | base64`

**Local setup not working?**
- Verify `.env` file exists and has all required variables
- For local development, ensure you have IAM user access keys (not OIDC)
- Check `json/data.json` exists and follows the correct format
- Test Docker connectivity: `docker run hello-world`

**S3 403 Forbidden errors?**
- Ensure the IAM role (for GitHub) or IAM user (for local) has the S3 permissions
- Verify bucket name matches in both the policy and your secrets/env file
- Check that `AWS_SESSION_TOKEN` is being passed for OIDC credentials
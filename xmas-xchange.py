import argparse, boto3, io, json, random
from datetime import datetime
from decouple import config
from twilio.rest import Client

#
# SETUP
#

# Parse command-line arguments
parser = argparse.ArgumentParser(description='Send Christmas Gift Exchange text messages')
parser.add_argument('--dry-run', action='store_true', help='Do a dry-run without sending SMS messages to the recipients.')
parser.add_argument('--hide-sensitive-output', action='store_true', help='Hide output messages that contain names and phone numbers (useful in GitHub actions).')
parser.add_argument('--github-test', action='store_true', help='This is used to tell the script that it is being run from GitHub actions so it will do some things differently.')
args = parser.parse_args()

if not args.dry_run and not args.github_test:
    # Configure Twilio credentials and phone number
    TWILIO_ACCOUNT_SID = config('TWILIO_ACCOUNT_SID')
    TWILIO_AUTH_TOKEN = config('TWILIO_AUTH_TOKEN')
    TWILIO_PHONE_NUMBER = config('TWILIO_PHONE_NUMBER')

# Configure AWS credentials and region
AWS_ACCESS_KEY_ID = config('AWS_ACCESS_KEY_ID')
AWS_SECRET_ACCESS_KEY = config('AWS_SECRET_ACCESS_KEY')
AWS_REGION = config('AWS_REGION')
S3_BUCKET = config('S3_BUCKET')
# Create an S3 client
s3 = boto3.client(
    's3', 
    aws_access_key_id=AWS_ACCESS_KEY_ID, 
    aws_secret_access_key=AWS_SECRET_ACCESS_KEY, 
    region_name=AWS_REGION
)

# ANSI escape codes for text color
R = '\033[91m'
G = '\033[92m'
Y = '\033[93m'
B = '\033[94m'
# ANSI escape codes for background color
BG_R = '\033[41m'
BG_G = '\033[42m'
BG_Y = '\033[43m'
BG_B = '\033[44m'
# ANSI escape code to reset color to default
EC = '\033[0m'

#
# SCRIPT
#

# Read combined data from the JSON file
with open('json/data.json') as json_file:
    people_info = json.load(json_file)

# Extract people_info and constraints from the combined data
constraints = {person: info['constraints'] for person, info in people_info.items()}
# List of people
people = list(people_info.keys())

# Check if the assignment meets the constraints
def check_constraints(assignment):
    for person, recipient in assignment.items():
        if recipient in constraints.get(person, []):
            return False
    return True

# Keep shuffling until a valid assignment is found
while True:
    random.shuffle(people)
    assignment = {people[i]: people[(i + 1) % len(people)] for i in range(len(people))}
    if check_constraints(assignment):
        break

def upload_assigment_data_to_s3(assigment_data):
    # Create a BytesIO object
    file_object = io.BytesIO(assignment_data.encode())
    # Upload file object to S3
    current_datetime = datetime.now().strftime("%Y-%m-%d_%H-%M-%S")
    file_prefix = 'github_' if args.github_test else ''
    file_suffix = '_dryrun' if args.dry_run or args.github_test else ''
    file_name = f"{file_prefix}{current_datetime}_gift_assignments{file_suffix}.txt"
    s3.upload_fileobj(file_object, S3_BUCKET, file_name)
    # Output what we did
    if not args.dry_run and not args.github_test:
        print(f"Gift assignments sent successfully and file has been written to S3 as {file_name}.")
    else:
        if not args.github_test:
            new_line = '\n' if not args.hide_sensitive_output else ''
            print(f"{new_line}{R}DRY RUN{EC} of gift assignments performed and file has been written to S3 as {G}{file_name}{EC}.")
        else:
            print(f"GitHub test of gift assignments performed and file has been written to S3 as {file_name}.")

# Sort the assignments alphabetically by name
assignment = dict(sorted(assignment.items(), key=lambda x: x[0]))
assignment_data = ""

if not args.dry_run and not args.github_test:
    # Send messages using Twilio
    client = Client(TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN)
else:
    assignment_data += f"---------------------------- DRY RUN ----------------------------\n\n"
    if not args.github_test:
        new_line = '\n' if not args.hide_sensitive_output else ''
        print(f"{new_line}{BG_R}---------------------------- DRY RUN ----------------------------{EC}{new_line}")
    else:
        print("GitHub test DRY RUN")

for person, recipient in assignment.items():
    person_phone = people_info[person]['phone_number']
    recipient_phone = people_info[recipient]['phone_number']
    message = f"Hello {person}! Your gift recipient is {recipient}. Merry Christmas!"
   
    assignment_data += f"{person} -> {recipient}\n"
    assignment_data += f"  Preview message to {person} ({person_phone}):\n  {message}\n\n"

    if not args.dry_run and not args.github_test:
        # Send a message via Twilo
        client.messages.create(
            to=person_phone,
            from_=TWILIO_PHONE_NUMBER,
            body=message
        )
    else:
        if not args.hide_sensitive_output and not args.github_test:
            # Output information to the terminal
            print(f"{G}{person}{EC} -> {Y}{recipient}{EC}")
            print(f"  Preview message to {G}{person}{EC} ({G}{person_phone}{EC}):\n  {B}{message}{EC}")

upload_assigment_data_to_s3(assignment_data)

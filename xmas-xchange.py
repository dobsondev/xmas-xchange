import argparse, json, random
from decouple import config
from helper import (
    setup_s3_client, setup_twilio_client, test_s3_connection, test_twilio_connection,
    upload_assignment_data_to_s3, Colors
)

def parse_arguments():
    """Parse command-line arguments"""
    parser = argparse.ArgumentParser(description='Send Christmas Gift Exchange text messages')
    parser.add_argument('--dry-run', action='store_true', help='Do a dry-run without sending SMS messages to the recipients.')
    parser.add_argument('--hide-sensitive-output', action='store_true', help='Hide output messages that contain names and phone numbers (useful in GitHub actions).')
    parser.add_argument('--github-test', action='store_true', help='This is used to tell the script that it is being run from GitHub actions so it will do some things differently.')
    return parser.parse_args()

def load_people_data():
    """Load people information from JSON file"""
    with open('json/data.json') as json_file:
        return json.load(json_file)

def check_constraints(assignment, constraints):
    """Check if the assignment meets the constraints"""
    for person, recipient in assignment.items():
        if recipient in constraints.get(person, []):
            return False
    return True

def generate_assignment(people, constraints):
    """Generate a valid gift exchange assignment"""
    while True:
        random.shuffle(people)
        assignment = {people[i]: people[(i + 1) % len(people)] for i in range(len(people))}
        if check_constraints(assignment, constraints):
            return assignment

def print_dry_run_header(is_github_test, hide_sensitive_output):
    """Print the dry run header with appropriate formatting"""
    if is_github_test:
        print("GitHub test DRY RUN")
    else:
        new_line = '\n' if not hide_sensitive_output else ''
        print(f"{new_line}{Colors.BG_RED}---------------------------- DRY RUN ----------------------------{Colors.END}{new_line}")

def create_message(person, recipient):
    """Create the gift exchange message"""
    return f"Welcome to the Gift Exchange!\n\nHello {person}! Your gift recipient is {recipient}. Merry Christmas!\n\nReply STOP to unsubscribe."

def format_assignment_data(person, recipient, person_phone, message):
    """Format assignment data for output"""
    return f"{person} -> {recipient}\n  Preview message to {person} ({person_phone}):\n  {message}\n\n"

def print_assignment_info(person, recipient, person_phone, message, hide_sensitive_output, is_github_test):
    """Print assignment information to console (dry run only)"""
    if not hide_sensitive_output and not is_github_test:
        print(f"{Colors.GREEN}{person}{Colors.END} -> {Colors.YELLOW}{recipient}{Colors.END}")
        print(f"  Preview message to {Colors.GREEN}{person}{Colors.END} ({Colors.GREEN}{person_phone}{Colors.END}):\n  {Colors.BLUE}{message}{Colors.END}")

def send_message(twilio_client, person_phone, message):
    """Send SMS message via Twilio"""
    twilio_phone = config('TWILIO_PHONE_NUMBER')
    twilio_client.messages.create(
        to=person_phone,
        from_=twilio_phone,
        body=message
    )

def print_upload_result(file_name, is_dry_run, is_github_test, hide_sensitive_output):
    """Print the S3 upload result message"""
    if not is_dry_run and not is_github_test:
        print(f"Gift assignments sent successfully and file has been written to S3 as {file_name}.")
    elif is_github_test:
        print(f"GitHub test of gift assignments performed and file has been written to S3 as {file_name}.")
    else:
        new_line = '\n' if not hide_sensitive_output else ''
        print(f"{new_line}{Colors.RED}DRY RUN{Colors.END} of gift assignments performed and file has been written to S3 as {Colors.GREEN}{file_name}{Colors.END}.")

def main():
    # Parse arguments and load data
    args = parse_arguments()
    people_info = load_people_data()
    
    # Extract constraints and people list
    constraints = {person: info['constraints'] for person, info in people_info.items()}
    people = list(people_info.keys())
    
    # Setup and test connections
    s3_client = setup_s3_client()
    if not test_s3_connection(s3_client):
        exit(1)
    
    if not test_twilio_connection():
        exit(1)
    
    # Generate assignment
    assignment = generate_assignment(people, constraints)
    assignment = dict(sorted(assignment.items(), key=lambda x: x[0]))  # Sort alphabetically
    
    # Setup for message sending or dry run
    is_dry_run = args.dry_run or args.github_test
    twilio_client = None if is_dry_run else setup_twilio_client()
    
    # Print dry run header if needed
    assignment_data = ""
    if is_dry_run:
        assignment_data += "---------------------------- DRY RUN ----------------------------\n\n"
        print_dry_run_header(args.github_test, args.hide_sensitive_output)
    
    # Process each assignment
    for person, recipient in assignment.items():
        person_phone = people_info[person]['phone_number']
        message = create_message(person, recipient)
        
        # Add to assignment data for S3 upload
        assignment_data += format_assignment_data(person, recipient, person_phone, message)
        
        # Send message or print info
        if is_dry_run:
            print_assignment_info(person, recipient, person_phone, message, args.hide_sensitive_output, args.github_test)
        else:
            send_message(twilio_client, person_phone, message)
    
    # Upload to S3 and print result
    file_name = upload_assignment_data_to_s3(s3_client, assignment_data, args.dry_run, args.github_test)
    print_upload_result(file_name, args.dry_run, args.github_test, args.hide_sensitive_output)

if __name__ == '__main__':
    main()
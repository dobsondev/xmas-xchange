#!/usr/bin/env python3
import argparse
import boto3
import re
import io
from datetime import datetime
from decouple import config
from twilio.rest import Client

# ANSI color constants
class Colors:
    RED = '\033[91m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    BLUE = '\033[94m'
    BG_RED = '\033[41m'
    BG_GREEN = '\033[42m'
    BG_YELLOW = '\033[43m'
    BG_BLUE = '\033[44m'
    END = '\033[0m'

def setup_s3_client():
    """Setup and return S3 client using environment variables"""
    aws_access_key_id = config('AWS_ACCESS_KEY_ID')
    aws_secret_access_key = config('AWS_SECRET_ACCESS_KEY')
    aws_region = config('AWS_REGION')
    
    return boto3.client(
        's3',
        aws_access_key_id=aws_access_key_id,
        aws_secret_access_key=aws_secret_access_key,
        region_name=aws_region
    )

def setup_twilio_client():
    """Setup and return Twilio client using environment variables"""
    account_sid = config('TWILIO_ACCOUNT_SID')
    auth_token = config('TWILIO_AUTH_TOKEN')
    return Client(account_sid, auth_token)

def test_s3_connection(s3_client):
    """Test S3 connection and return True if successful"""
    try:
        bucket_name = config('S3_BUCKET')
        s3_client.head_bucket(Bucket=bucket_name)
        print("✅ S3 connection successful!")
        return True
    except Exception as e:
        print(f"❌ S3 connection failed: {e}")
        return False

def test_twilio_connection():
    """Test Twilio connection and return True if successful"""
    try:
        client = setup_twilio_client()
        account_sid = config('TWILIO_ACCOUNT_SID')
        client.api.accounts.get(account_sid)
        print("✅ Twilio connection successful!")
        return True
    except Exception as e:
        print(f"❌ Twilio connection failed: {e}")
        return False

def upload_assignment_data_to_s3(s3_client, assignment_data, is_dry_run=False, is_github_test=False):
    """Upload assignment data to S3 with appropriate filename"""
    bucket_name = config('S3_BUCKET')
    
    # Create a BytesIO object
    file_object = io.BytesIO(assignment_data.encode())
    
    # Generate filename
    current_datetime = datetime.now().strftime("%Y-%m-%d_%H-%M-%S")
    file_prefix = 'github_' if is_github_test else ''
    file_suffix = '_dryrun' if is_dry_run or is_github_test else ''
    file_name = f"{file_prefix}{current_datetime}_gift_assignments{file_suffix}.txt"
    
    # Upload file object to S3
    s3_client.upload_fileobj(file_object, bucket_name, file_name)
    
    return file_name

def get_assignment_file_content(s3_client, bucket_name, file_name):
    """Download and return the content of an assignment file from S3"""
    try:
        response = s3_client.get_object(Bucket=bucket_name, Key=file_name)
        return response['Body'].read().decode('utf-8')
    except Exception as e:
        print(f"Error downloading file {file_name}: {e}")
        return None

def parse_assignments(content):
    """Parse assignment content and return a dictionary of person -> recipient"""
    assignments = {}
    lines = content.split('\n')
    
    for line in lines:
        # Look for lines in format "Person -> Recipient"
        match = re.match(r'^(\w+)\s*->\s*(\w+)$', line)
        if match:
            person, recipient = match.groups()
            assignments[person] = recipient
    
    return assignments

def main():
    parser = argparse.ArgumentParser(description='Query gift exchange assignments from S3')
    parser.add_argument('file_name', help='Name of the assignment file in S3 (e.g., 2024-12-01_12-00-00_gift_assignments.txt)')
    parser.add_argument('person_name', help='Name of the person to query')
    
    args = parser.parse_args()
    
    # Setup S3 client and test connection
    s3_client = setup_s3_client()
    if not test_s3_connection(s3_client):
        return 1
    
    # Download and parse the assignment file
    bucket_name = config('S3_BUCKET')
    content = get_assignment_file_content(s3_client, bucket_name, args.file_name)
    if content is None:
        return 1
    
    assignments = parse_assignments(content)
    
    # Look up the person's assignment
    if args.person_name in assignments:
        recipient = assignments[args.person_name]
        print(f"{args.person_name} -> {recipient}")
    else:
        print(f"Person '{args.person_name}' not found in assignments.")
        print("Available people:")
        for person in sorted(assignments.keys()):
            print(f"  - {person}")
        return 1
    
    return 0

if __name__ == '__main__':
    exit(main())
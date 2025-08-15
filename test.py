#!/usr/bin/env python3
"""
Test script for the gift exchange workflow.
This script validates the complete process including constraint checking.
"""

import json
import subprocess
import sys
import tempfile
import os
import re
from datetime import datetime
from helper import setup_s3_client, get_assignment_file_content, parse_assignments, Colors

def load_people_data():
    """Load people information and constraints from JSON file"""
    with open('json/data.json') as json_file:
        data = json.load(json_file)
    
    people = list(data.keys())
    constraints = {person: info['constraints'] for person, info in data.items()}
    return people, constraints, data

def validate_constraints(assignments, constraints):
    """Validate that no assignments violate constraints"""
    violations = []
    
    for person, recipient in assignments.items():
        if recipient in constraints.get(person, []):
            violations.append(f"{person} cannot give to {recipient} (constraint violation)")
    
    return violations

def validate_assignment_completeness(assignments, people):
    """Ensure everyone gives and receives exactly once"""
    givers = set(assignments.keys())
    recipients = set(assignments.values())
    people_set = set(people)
    
    issues = []
    
    # Check if everyone is a giver
    missing_givers = people_set - givers
    if missing_givers:
        issues.append(f"Missing givers: {', '.join(missing_givers)}")
    
    # Check if everyone is a recipient
    missing_recipients = people_set - recipients
    if missing_recipients:
        issues.append(f"Missing recipients: {', '.join(missing_recipients)}")
    
    # Check for duplicate recipients
    recipient_counts = {}
    for recipient in assignments.values():
        recipient_counts[recipient] = recipient_counts.get(recipient, 0) + 1
    
    duplicates = [name for name, count in recipient_counts.items() if count > 1]
    if duplicates:
        issues.append(f"Duplicate recipients: {', '.join(duplicates)}")
    
    return issues

def run_dry_run():
    """Run the main script in dry-run mode and capture the output"""
    try:
        # Import and run the main script directly (since we're already in the container)
        import sys
        import io
        from contextlib import redirect_stdout, redirect_stderr
        
        # Capture stdout to get the filename
        captured_output = io.StringIO()
        
        # Set up command line arguments for dry-run
        original_argv = sys.argv
        sys.argv = ['xmas-xchange.py', '--dry-run', '--hide-sensitive-output']
        
        try:
            # Import and run the main script using importlib
            import importlib.util
            import types
            
            # Load the module dynamically
            spec = importlib.util.spec_from_file_location("xmas_xchange", "xmas-xchange.py")
            xmas_xchange = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(xmas_xchange)
            
            with redirect_stdout(captured_output):
                xmas_xchange.main()
            
            output = captured_output.getvalue()
            
            # Extract filename from output
            filename_match = re.search(r'(\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}_gift_assignments_dryrun\.txt)', output)
            if filename_match:
                filename = filename_match.group(1)
                return filename
            else:
                print(f"{Colors.RED}❌ Could not extract filename from output{Colors.END}")
                print("Output:", output)
                return None
                
        finally:
            sys.argv = original_argv
            
    except Exception as e:
        print(f"{Colors.RED}❌ Dry-run failed: {e}{Colors.END}")
        return None

def test_helper_queries(filename, test_people):
    """Test the helper script with specific people"""
    results = {}
    for person in test_people:
        try:
            # Import and use helper functions directly
            import sys
            import io
            from contextlib import redirect_stdout
            
            # Capture stdout to get the result
            captured_output = io.StringIO()
            
            # Set up command line arguments for helper script
            original_argv = sys.argv
            sys.argv = ['helper.py', filename, person]
            
            try:
                import importlib
                import helper
                importlib.reload(helper)  # Reload to ensure fresh run
                
                with redirect_stdout(captured_output):
                    result_code = helper.main()
                
                if result_code == 0:
                    output = captured_output.getvalue()
                    lines = output.strip().split('\n')
                    
                    # Find the assignment line
                    assignment_line = None
                    for line in lines:
                        if ' -> ' in line and person in line:
                            assignment_line = line
                            break
                    
                    if assignment_line:
                        # Extract recipient from "Person -> Recipient"
                        parts = assignment_line.split(' -> ')
                        if len(parts) == 2:
                            recipient = parts[1].strip()
                            results[person] = recipient
                            print(f"  {Colors.GREEN}✅ {person} -> {recipient}{Colors.END}")
                        else:
                            print(f"  {Colors.RED}❌ Could not parse assignment for {person}{Colors.END}")
                    else:
                        print(f"  {Colors.RED}❌ No assignment found for {person}{Colors.END}")
                else:
                    print(f"  {Colors.RED}❌ Helper query failed for {person} (exit code: {result_code}){Colors.END}")
                    
            finally:
                sys.argv = original_argv
                
        except Exception as e:
            print(f"  {Colors.RED}❌ Helper query failed for {person}: {e}{Colors.END}")
    
    return results

def download_and_validate_full_assignment(s3_client, filename):
    """Download the full assignment file and validate it completely"""
    print(f"{Colors.BLUE}Downloading and validating full assignment...{Colors.END}")
    
    try:
        from decouple import config
        bucket_name = config('S3_BUCKET')
        
        content = get_assignment_file_content(s3_client, bucket_name, filename)
        if content is None:
            print(f"{Colors.RED}❌ Failed to download assignment file{Colors.END}")
            return None
        
        assignments = parse_assignments(content)
        if not assignments:
            print(f"{Colors.RED}❌ No assignments found in file{Colors.END}")
            return None
        
        print(f"{Colors.GREEN}✅ Downloaded assignment file with {len(assignments)} assignments{Colors.END}")
        return assignments
        
    except Exception as e:
        print(f"{Colors.RED}❌ Error downloading/parsing assignment: {e}{Colors.END}")
        return None

def run_comprehensive_test():
    """Run the complete test workflow"""
    print(f"{Colors.BG_BLUE}=== Gift Exchange Workflow Test ==={Colors.END}\n")
    
    # Load test data
    people, constraints, people_data = load_people_data()
    print(f"Loaded {len(people)} people with constraints")
    
    # Display constraints for clarity
    print(f"{Colors.YELLOW}Constraints:{Colors.END}")
    for person, person_constraints in constraints.items():
        if person_constraints:
            print(f"  {person} cannot give to: {', '.join(person_constraints)}")
        else:
            print(f"  {person} has no constraints")
    print()
    
    # Test connections before proceeding
    print(f"{Colors.BLUE}Testing service connections...{Colors.END}")
    
    # Test S3 connection
    s3_client = setup_s3_client()
    from helper import test_s3_connection, test_twilio_connection
    if not test_s3_connection(s3_client):
        print(f"{Colors.RED}❌ Test failed: S3 connection failed{Colors.END}")
        return False
    
    # Test Twilio connection
    if not test_twilio_connection():
        print(f"{Colors.RED}❌ Test failed: Twilio connection failed{Colors.END}")
        return False
    
    print(f"{Colors.GREEN}✅ All service connections successful{Colors.END}")
    print()
    
    # Step 1: Run dry-run
    print(f"{Colors.BLUE}Running dry-run...{Colors.END}")
    filename = run_dry_run()
    if not filename:
        print(f"{Colors.RED}❌ Test failed at dry-run step{Colors.END}")
        return False
    print(f"{Colors.GREEN}✅ Dry-run completed successfully{Colors.END}")
    print()
    
    # Step 2: Test helper queries for all people
    print(f"{Colors.BLUE}Testing helper queries for all {len(people)} people{Colors.END}")
    helper_results = test_helper_queries(filename, people)
    
    if len(helper_results) != len(people):
        print(f"{Colors.RED}❌ Helper queries failed - got {len(helper_results)} results for {len(people)} people{Colors.END}")
        return False
    print()
    
    # Step 3: Download and validate complete assignment  
    full_assignments = download_and_validate_full_assignment(s3_client, filename)
    if not full_assignments:
        return False
    print()
    
    # Step 4: Validate constraints
    print(f"{Colors.BLUE}Validating constraints...{Colors.END}")
    constraint_violations = validate_constraints(full_assignments, constraints)
    
    if constraint_violations:
        print(f"{Colors.RED}❌ Constraint violations found:{Colors.END}")
        for violation in constraint_violations:
            print(f"  - {violation}")
        return False
    else:
        print(f"{Colors.GREEN}✅ All constraints satisfied{Colors.END}")
    
    print(f"{Colors.BLUE}Validating assignment completeness...{Colors.END}")
    completeness_issues = validate_assignment_completeness(full_assignments, people)
    
    if completeness_issues:
        print(f"{Colors.RED}❌ Assignment completeness issues:{Colors.END}")
        for issue in completeness_issues:
            print(f"  - {issue}")
        return False
    else:
        print(f"{Colors.GREEN}✅ Assignment is complete and valid{Colors.END}")
    print()
    
    # Step 6: Cross-validate helper results with full assignment
    print(f"{Colors.BLUE}Cross-validating all {len(helper_results)} helper results...{Colors.END}")
    for person, helper_recipient in helper_results.items():
        full_recipient = full_assignments.get(person)
        if helper_recipient == full_recipient:
            print(f"  {Colors.GREEN}✅ {person} -> {helper_recipient} (matches){Colors.END}")
        else:
            print(f"  {Colors.RED}❌ {person}: helper says {helper_recipient}, full assignment says {full_recipient}{Colors.END}")
            return False
    
    print(f"\n{Colors.BG_GREEN}🎉 ALL TESTS PASSED! 🎉{Colors.END}")
    print(f"Generated assignment file: {filename}")
    return True

def main():
    """Main test function"""
    success = run_comprehensive_test()
    sys.exit(0 if success else 1)

if __name__ == '__main__':
    main()
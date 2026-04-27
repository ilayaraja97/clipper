#!/bin/bash

queries=(
    "Analyze the testdata/users.json file and show me the active users"
    "Parse the testdata/employees.csv and calculate average salary by department"
    "Extract all email addresses from the testdata/special_characters.txt file"
    "List all processes listening on port 8080"
    "Show the disk usage of Docker resources"
    "Find files larger than 100MB in /var/log"
    "Check open network connections and their owning processes"
    "Generate a shell command to clean up old log files"
    "Review the testdata/sample_script.py for potential improvements"
    "Check the testdata/sample_script.sh for security issues to fix problems"
    "Convert the testdata/sample_config.ini to YAML format"
    "Check the testdata/application.log for errors in the last hour"
    "Count occurrences of each log level in testdata/application.log"
    "Summarize the testdata/large_document.txt"
    "Extract all code blocks from testdata/sample_docs.md"
    "Explain what the testdata/docker-compose.yml does"
    "Generate sample data for the tables in testdata/sample_schema.sql"
    "Parse testdata/special_characters.txt and identify all Unicode characters"
    "Find all emojis in testdata/special_characters.txt"
)

for query in "${queries[@]}"; do
    echo "Running: $query"
    clipper -c "$query"
    echo "----------------------------------------" 
done
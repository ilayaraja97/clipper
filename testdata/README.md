# Test Data Directory

This directory contains mock files with synthetic content designed to test clipper's capabilities for processing various types of files and content.

## Files Overview

### Configuration Files
- `sample_config.ini` - INI-style configuration file with database, API, and logging settings
- `docker-compose.yml` - Docker Compose configuration for multi-service application

### Data Files
- `users.json` - JSON file with user data including nested objects and arrays
- `employees.csv` - CSV file with employee information and various data types

### Code Files
- `sample_script.py` - Python script demonstrating classes, functions, and data processing
- `sample_script.sh` - Bash shell script with logging, error handling, and argument parsing
- `sample_schema.sql` - SQL file with table creation, data insertion, and sample queries

### Documentation
- `sample_docs.md` - Markdown file with various formatting elements (tables, lists, code blocks, etc.)
- `large_document.txt` - Large text file (2000+ words) for testing performance with bigger content

### Logs and Special Content
- `application.log` - Sample application log with various log levels and timestamps
- `special_characters.txt` - File with Unicode characters, emojis, special symbols, and edge cases

## Usage for Testing

These files can be used to test clipper's ability to:

1. **Parse different file formats** - JSON, CSV, INI, YAML, SQL, Markdown
2. **Handle various content types** - code, configuration, documentation, logs
3. **Process special characters** - Unicode, emojis, symbols, international text
4. **Work with large files** - performance testing with bigger content
5. **Execute commands on files** - analysis, transformation, searching

## Integration Test Examples

You can use these files to test commands like:

- "Analyze the users.json file and show me the active users"
- "Parse the employees.csv and calculate average salary by department"
- "Review the sample_script.py for potential improvements"
- "Check the application.log for errors in the last hour"
- "Summarize the large_document.txt"
- "Convert the sample_config.ini to YAML format"
- "Find all email addresses in the special_characters.txt file"

## Adding More Test Files

When adding new test files, consider:
- File size variations (small, medium, large)
- Different encodings (UTF-8, ASCII, etc.)
- Various data structures and formats
- Edge cases and special characters
- Real-world scenarios clipper might encounter
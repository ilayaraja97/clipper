#!/usr/bin/env python3
"""
Sample Python script for testing clipper capabilities
This script demonstrates various Python features and patterns
"""

import json
import sys
from datetime import datetime
from typing import List, Dict, Optional

class DataProcessor:
    """A sample class for processing data"""

    def __init__(self, config: Dict[str, any]):
        self.config = config
        self.processed_count = 0

    def process_data(self, data: List[Dict]) -> List[Dict]:
        """Process a list of data items"""
        results = []
        for item in data:
            try:
                processed = self._process_single_item(item)
                results.append(processed)
                self.processed_count += 1
            except Exception as e:
                print(f"Error processing item {item.get('id', 'unknown')}: {e}")
                continue
        return results

    def _process_single_item(self, item: Dict) -> Dict:
        """Process a single data item"""
        # Add timestamp
        item['processed_at'] = datetime.now().isoformat()

        # Validate required fields
        required_fields = ['id', 'name', 'value']
        for field in required_fields:
            if field not in item:
                raise ValueError(f"Missing required field: {field}")

        # Calculate derived values
        if 'value' in item and isinstance(item['value'], (int, float)):
            item['value_doubled'] = item['value'] * 2
            item['value_category'] = 'high' if item['value'] > 100 else 'low'

        return item

def main():
    """Main function"""
    # Sample configuration
    config = {
        'batch_size': 100,
        'output_format': 'json',
        'debug': True
    }

    # Sample data
    sample_data = [
        {'id': 1, 'name': 'Item A', 'value': 50},
        {'id': 2, 'name': 'Item B', 'value': 150},
        {'id': 3, 'name': 'Item C', 'value': 75},
    ]

    # Process data
    processor = DataProcessor(config)
    results = processor.process_data(sample_data)

    # Output results
    if config['output_format'] == 'json':
        print(json.dumps(results, indent=2))
    else:
        for result in results:
            print(f"Processed: {result['name']} (ID: {result['id']})")

    print(f"\nProcessed {processor.processed_count} items successfully")

if __name__ == '__main__':
    main()
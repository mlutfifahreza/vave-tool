#!/usr/bin/env python3
"""
Traffic Simulator for API Load Testing
Generates HTTP requests with configurable duration and QPS (queries per second).
"""

import argparse
import time
import requests
from datetime import datetime
from concurrent.futures import ThreadPoolExecutor, as_completed
from collections import defaultdict
import sys


class TrafficSimulator:
    def __init__(self, url, duration, qps, timeout=10):
        self.url = url
        self.duration = duration
        self.qps = qps
        self.timeout = timeout
        self.results = defaultdict(int)
        self.start_time = None
        self.end_time = None
        # Create session for connection pooling to reduce connection overhead
        self.session = requests.Session()
        # Pre-warm connection pool
        adapter = requests.adapters.HTTPAdapter(
            pool_connections=min(qps, 200),
            pool_maxsize=min(qps, 200),
            max_retries=0
        )
        self.session.mount('http://', adapter)
        self.session.mount('https://', adapter)
        
    def make_request(self):
        """Make a single HTTP request and return the result."""
        try:
            # Use session for connection reuse, measures only actual request/response time
            response = self.session.get(self.url, timeout=self.timeout)
            return {
                'status': response.status_code,
                'success': response.status_code == 200,
                'latency': response.elapsed.total_seconds()
            }
        except requests.exceptions.Timeout:
            return {'status': 'timeout', 'success': False, 'latency': self.timeout}
        except requests.exceptions.RequestException as e:
            return {'status': 'error', 'success': False, 'latency': 0}
    
    def print_progress(self, completed, total, success, failed):
        """Print progress bar and stats."""
        bar_length = 50
        progress = completed / total
        filled = int(bar_length * progress)
        bar = '█' * filled + '░' * (bar_length - filled)
        
        success_rate = (success / completed * 100) if completed > 0 else 0
        
        sys.stdout.write(f'\r[{bar}] {completed}/{total} | '
                        f'✓ {success} | ✗ {failed} | '
                        f'Success: {success_rate:.1f}%')
        sys.stdout.flush()
    
    def run(self):
        """Run the traffic simulation."""
        print("🚀 Starting traffic simulation...")
        print(f"   Duration: {self.duration} seconds")
        print(f"   QPS:      {self.qps} requests/second")
        print(f"   Target:   {self.url}")
        print(f"   Started:  {datetime.now().strftime('%H:%M:%S')}")
        print()
        
        total_requests = self.duration * self.qps
        interval = 1.0 / self.qps
        
        self.start_time = time.time()
        request_times = []
        completed = 0
        success_count = 0
        failed_count = 0
        latencies = []
        status_codes = defaultdict(int)
        
        # Use thread pool for concurrent requests
        # Calculate workers based on QPS and expected latency to prevent queuing
        # Assuming average latency of 100ms, we need QPS * 0.1 workers minimum
        # Adding 50% buffer and capping at reasonable maximum
        min_workers = max(10, int(self.qps * 0.15))
        max_workers_limit = min(min_workers, 2000)
        
        with ThreadPoolExecutor(max_workers=max_workers_limit) as executor:
            futures = []
            
            for i in range(total_requests):
                # Schedule request
                scheduled_time = self.start_time + (i * interval)
                current_time = time.time()
                
                if scheduled_time > current_time:
                    time.sleep(scheduled_time - current_time)
                
                # Check if duration exceeded
                if time.time() - self.start_time >= self.duration:
                    break
                
                future = executor.submit(self.make_request)
                futures.append(future)
                request_times.append(time.time())
            
            # Collect results
            for future in as_completed(futures):
                result = future.result()
                completed += 1
                
                if result['success']:
                    success_count += 1
                else:
                    failed_count += 1
                
                latencies.append(result['latency'])
                status_codes[result['status']] += 1
                
                # Update progress every 10 requests or at the end
                if completed % 10 == 0 or completed == len(futures):
                    self.print_progress(completed, len(futures), success_count, failed_count)
        
        # Close session
        self.session.close()
        
        self.end_time = time.time()
        actual_duration = self.end_time - self.start_time
        
        print("\n")
        print("📊 Traffic Simulation Complete!")
        print(f"   Total Requests:    {completed}")
        print(f"   Successful:        {success_count}")
        print(f"   Failed:            {failed_count}")
        print(f"   Success Rate:      {(success_count/completed*100):.2f}%")
        print(f"   Actual Duration:   {actual_duration:.2f}s")
        print(f"   Actual QPS:        {(completed/actual_duration):.2f}")
        
        if latencies:
            latencies.sort()
            print()
            print("📈 Latency Statistics:")
            print(f"   Min:     {min(latencies)*1000:.2f}ms")
            print(f"   Max:     {max(latencies)*1000:.2f}ms")
            print(f"   Mean:    {(sum(latencies)/len(latencies))*1000:.2f}ms")
            print(f"   Median:  {latencies[len(latencies)//2]*1000:.2f}ms")
            print(f"   P95:     {latencies[int(len(latencies)*0.95)]*1000:.2f}ms")
            print(f"   P99:     {latencies[int(len(latencies)*0.99)]*1000:.2f}ms")
        
        print()
        print("📋 Status Code Distribution:")
        for status, count in sorted(status_codes.items()):
            percentage = (count / completed * 100)
            print(f"   {status}:  {count} ({percentage:.1f}%)")
        
        print()
        print("✨ View results in Grafana: http://localhost:3000")


def main():
    parser = argparse.ArgumentParser(
        description='Generate HTTP traffic for load testing',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s --duration 30 --qps 10
  %(prog)s --duration 60 --qps 100 --url http://localhost:8080/api/products/1
  %(prog)s -d 120 -q 50
        """
    )
    
    parser.add_argument(
        '-u', '--url',
        default='http://localhost:8080/api/products',
        help='Target URL (default: http://localhost:8080/api/products)'
    )
    parser.add_argument(
        '-d', '--duration',
        type=int,
        default=30,
        help='Duration in seconds (default: 30)'
    )
    parser.add_argument(
        '-q', '--qps',
        type=int,
        default=5,
        help='Queries per second (default: 5)'
    )
    parser.add_argument(
        '-t', '--timeout',
        type=int,
        default=10,
        help='Request timeout in seconds (default: 10)'
    )
    
    args = parser.parse_args()
    
    if args.duration <= 0:
        print("Error: Duration must be positive")
        sys.exit(1)
    
    if args.qps <= 0:
        print("Error: QPS must be positive")
        sys.exit(1)
    
    try:
        simulator = TrafficSimulator(
            url=args.url,
            duration=args.duration,
            qps=args.qps,
            timeout=args.timeout
        )
        simulator.run()
    except KeyboardInterrupt:
        print("\n\n⚠️  Simulation interrupted by user")
        sys.exit(0)
    except Exception as e:
        print(f"\n\n❌ Error: {e}")
        sys.exit(1)


if __name__ == '__main__':
    main()

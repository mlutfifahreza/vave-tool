#!/usr/bin/env python3
"""
User Behavior Simulator for API Load Testing
Simulates concurrent users performing realistic workflows:
1. Get list of products
2. View 5 individual products from the list
"""

import argparse
import time
import requests
import random
from datetime import datetime
from concurrent.futures import ThreadPoolExecutor, as_completed
from collections import defaultdict
import sys
import threading


class UserBehaviorSimulator:
    def __init__(self, base_url, concurrent_users, duration, think_time_min=0.5, think_time_max=2.0, timeout=10):
        self.base_url = base_url.rstrip('/')
        self.concurrent_users = concurrent_users
        self.duration = duration
        self.think_time_min = think_time_min
        self.think_time_max = think_time_max
        self.timeout = timeout
        self.start_time = None
        self.end_time = None
        self.stop_flag = threading.Event()
        
        self.stats_lock = threading.Lock()
        self.total_flows = 0
        self.successful_flows = 0
        self.failed_flows = 0
        self.total_requests = 0
        self.successful_requests = 0
        self.failed_requests = 0
        self.latencies = []
        self.status_codes = defaultdict(int)
        self.error_types = defaultdict(int)
        
    def make_request(self, method, url, description):
        """Make a single HTTP request and return the result."""
        try:
            start_time = time.time()
            response = requests.request(method, url, timeout=self.timeout)
            latency = time.time() - start_time
            
            return {
                'success': response.status_code == 200,
                'status': response.status_code,
                'latency': latency,
                'data': response.json() if response.status_code == 200 else None,
                'description': description
            }
        except requests.exceptions.Timeout:
            return {
                'success': False,
                'status': 'timeout',
                'latency': self.timeout,
                'data': None,
                'description': description
            }
        except requests.exceptions.RequestException as e:
            return {
                'success': False,
                'status': 'error',
                'latency': 0,
                'data': None,
                'description': description,
                'error': str(e)
            }
    
    def user_session(self, user_id):
        """Simulate a single user session performing the workflow repeatedly."""
        flows_completed = 0
        flows_succeeded = 0
        flows_failed = 0
        
        while not self.stop_flag.is_set():
            flow_start = time.time()
            flow_success = True
            
            # Step 1: Get list of products with random pagination
            random_page = random.randint(1, 5)
            random_size = random.randint(5, 20)
            list_result = self.make_request(
                'GET',
                f'{self.base_url}/api/products?page={random_page}&size={random_size}',
                f'List Products (page={random_page}, size={random_size})'
            )
            
            with self.stats_lock:
                self.total_requests += 1
                self.latencies.append(list_result['latency'])
                self.status_codes[list_result['status']] += 1
                
                if list_result['success']:
                    self.successful_requests += 1
                else:
                    self.failed_requests += 1
                    self.error_types[f"List: {list_result['status']}"] += 1
                    flow_success = False
            
            # Step 2: Extract product IDs and view individual products
            product_ids = []
            if list_result['success'] and list_result['data']:
                data = list_result['data'].get('data', {})
                products = data.get('products', []) if isinstance(data, dict) else []
                if isinstance(products, list) and len(products) > 0:
                    # Get up to 5 random products
                    sample_size = min(5, len(products))
                    sampled_products = random.sample(products, sample_size)
                    product_ids = [p.get('id') for p in sampled_products if p.get('id')]
            
            # If we got product IDs, view them
            if product_ids:
                for product_id in product_ids:
                    if self.stop_flag.is_set():
                        break
                    
                    # Think time - simulate user reading the list before clicking
                    think_time = random.uniform(self.think_time_min, self.think_time_max)
                    time.sleep(think_time)
                    
                    # Get individual product
                    product_result = self.make_request(
                        'GET',
                        f'{self.base_url}/api/products/{product_id}',
                        f'Get Product {product_id}'
                    )
                    
                    with self.stats_lock:
                        self.total_requests += 1
                        self.latencies.append(product_result['latency'])
                        self.status_codes[product_result['status']] += 1
                        
                        if product_result['success']:
                            self.successful_requests += 1
                        else:
                            self.failed_requests += 1
                            self.error_types[f"Detail: {product_result['status']}"] += 1
                            flow_success = False
            else:
                flow_success = False
            
            # Record flow completion
            flows_completed += 1
            if flow_success:
                flows_succeeded += 1
            else:
                flows_failed += 1
            
            with self.stats_lock:
                self.total_flows += 1
                if flow_success:
                    self.successful_flows += 1
                else:
                    self.failed_flows += 1
            
            # Small delay before starting next flow
            if not self.stop_flag.is_set():
                time.sleep(random.uniform(0.5, 1.5))
        
        return {
            'user_id': user_id,
            'flows_completed': flows_completed,
            'flows_succeeded': flows_succeeded,
            'flows_failed': flows_failed
        }
    
    def print_progress(self):
        """Print periodic progress updates."""
        while not self.stop_flag.is_set():
            elapsed = time.time() - self.start_time
            remaining = max(0, self.duration - elapsed)
            
            with self.stats_lock:
                if self.total_requests > 0:
                    success_rate = (self.successful_requests / self.total_requests * 100)
                    avg_latency = sum(self.latencies) / len(self.latencies) if self.latencies else 0
                    if self.latencies:
                        sorted_latencies = sorted(self.latencies)
                        p99 = sorted_latencies[int(len(sorted_latencies) * 0.99)]
                    else:
                        p99 = 0
                else:
                    success_rate = 0
                    avg_latency = 0
                    p99 = 0
            
            sys.stdout.write(
                f'\r⏱️  {int(elapsed)}s / {self.duration}s | '
                f'👥 {self.concurrent_users} users | '
                f'🔄 {self.total_flows} flows | '
                f'📨 {self.total_requests} reqs | '
                f'✓ {success_rate:.1f}% | '
                f'⚡ {avg_latency*1000:.0f}ms avg, {p99*1000:.0f}ms p99'
            )
            sys.stdout.flush()
            
            time.sleep(1)
    
    def run(self):
        """Run the user behavior simulation."""
        print("🚀 Starting user behavior simulation...")
        print(f"   Concurrent Users:  {self.concurrent_users}")
        print(f"   Duration:          {self.duration} seconds")
        print(f"   Base URL:          {self.base_url}")
        print(f"   Think Time:        {self.think_time_min}s - {self.think_time_max}s")
        print(f"   Started:           {datetime.now().strftime('%H:%M:%S')}")
        print()
        print("📋 User Flow: (random page 1-5, size 5-20)") 
        print("   1. Get list of products")
        print("   2. View 5 random products from the list")
        print("   3. Repeat until duration expires")
        print()
        
        self.start_time = time.time()
        
        # Start progress printer thread
        progress_thread = threading.Thread(target=self.print_progress, daemon=True)
        progress_thread.start()
        
        # Start user sessions
        with ThreadPoolExecutor(max_workers=self.concurrent_users) as executor:
            # Submit user sessions
            futures = [
                executor.submit(self.user_session, user_id)
                for user_id in range(1, self.concurrent_users + 1)
            ]
            
            # Wait for duration
            time.sleep(self.duration)
            
            # Signal all users to stop
            self.stop_flag.set()
            
            # Wait for all users to finish their current flow
            user_results = []
            for future in as_completed(futures):
                try:
                    result = future.result(timeout=30)
                    user_results.append(result)
                except Exception as e:
                    print(f"\n⚠️  User session error: {e}")
        
        self.end_time = time.time()
        actual_duration = self.end_time - self.start_time
        
        # Print final results
        print("\n")
        print("=" * 70)
        print("📊 User Behavior Simulation Complete!")
        print("=" * 70)
        print()
        
        print("👥 User Statistics:")
        print(f"   Concurrent Users:     {self.concurrent_users}")
        print(f"   Total User Flows:     {self.total_flows}")
        print(f"   Successful Flows:     {self.successful_flows}")
        print(f"   Failed Flows:         {self.failed_flows}")
        print(f"   Flow Success Rate:    {(self.successful_flows/self.total_flows*100):.2f}%")
        print()
        
        print("📨 Request Statistics:")
        print(f"   Total Requests:       {self.total_requests}")
        print(f"   Successful:           {self.successful_requests}")
        print(f"   Failed:               {self.failed_requests}")
        print(f"   Request Success Rate: {(self.successful_requests/self.total_requests*100):.2f}%")
        print(f"   Actual Duration:      {actual_duration:.2f}s")
        print(f"   Requests/sec:         {(self.total_requests/actual_duration):.2f}")
        print()
        
        if self.latencies:
            self.latencies.sort()
            print("📈 Latency Statistics:")
            print(f"   Min:     {min(self.latencies)*1000:.2f}ms")
            print(f"   Max:     {max(self.latencies)*1000:.2f}ms")
            print(f"   Mean:    {(sum(self.latencies)/len(self.latencies))*1000:.2f}ms")
            print(f"   Median:  {self.latencies[len(self.latencies)//2]*1000:.2f}ms")
            print(f"   P95:     {self.latencies[int(len(self.latencies)*0.95)]*1000:.2f}ms")
            print(f"   P99:     {self.latencies[int(len(self.latencies)*0.99)]*1000:.2f}ms")
            print()
        
        print("📋 Status Code Distribution:")
        for status, count in sorted(self.status_codes.items()):
            percentage = (count / self.total_requests * 100)
            print(f"   {status}:  {count} ({percentage:.1f}%)")
        
        if self.error_types:
            print()
            print("❌ Error Types:")
            for error_type, count in sorted(self.error_types.items(), key=lambda x: x[1], reverse=True):
                print(f"   {error_type}:  {count}")
        
        print()
        print("✨ View results in Grafana: http://localhost:3000")


def main():
    parser = argparse.ArgumentParser(
        description='Simulate concurrent user behavior for load testing',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s --users 10 --duration 60
  %(prog)s -u 50 -d 120 --url http://localhost:8080
  %(prog)s -u 100 -d 300 --think-time-min 1 --think-time-max 3
        """
    )
    
    parser.add_argument(
        '--url',
        default='http://localhost:8080',
        help='Base URL (default: http://localhost:8080)'
    )
    parser.add_argument(
        '-u', '--users',
        type=int,
        default=10,
        help='Number of concurrent users (default: 10)'
    )
    parser.add_argument(
        '-d', '--duration',
        type=int,
        default=60,
        help='Duration in seconds (default: 60)'
    )
    parser.add_argument(
        '--think-time-min',
        type=float,
        default=0.5,
        help='Minimum think time between requests in seconds (default: 0.5)'
    )
    parser.add_argument(
        '--think-time-max',
        type=float,
        default=2.0,
        help='Maximum think time between requests in seconds (default: 2.0)'
    )
    parser.add_argument(
        '-t', '--timeout',
        type=int,
        default=10,
        help='Request timeout in seconds (default: 10)'
    )
    
    args = parser.parse_args()
    
    if args.users <= 0:
        print("Error: Number of users must be positive")
        sys.exit(1)
    
    if args.duration <= 0:
        print("Error: Duration must be positive")
        sys.exit(1)
    
    if args.think_time_min < 0 or args.think_time_max < 0:
        print("Error: Think times must be non-negative")
        sys.exit(1)
    
    if args.think_time_min > args.think_time_max:
        print("Error: Minimum think time must be less than or equal to maximum")
        sys.exit(1)
    
    try:
        simulator = UserBehaviorSimulator(
            base_url=args.url,
            concurrent_users=args.users,
            duration=args.duration,
            think_time_min=args.think_time_min,
            think_time_max=args.think_time_max,
            timeout=args.timeout
        )
        simulator.run()
    except KeyboardInterrupt:
        print("\n\n⚠️  Simulation interrupted by user")
        sys.exit(0)
    except Exception as e:
        print(f"\n\n❌ Error: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == '__main__':
    main()

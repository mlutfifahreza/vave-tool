#!/usr/bin/env python3

import os
import redis
import sys

def get_env(key, default_value):
    return os.getenv(key, default_value)

def get_int_env(key, default_value):
    value = os.getenv(key)
    if value:
        try:
            return int(value)
        except ValueError:
            pass
    return default_value

def scan_and_delete_keys(client, pattern, description):
    print(f"Invalidating {description} cache...")
    keys = []
    cursor = 0
    while True:
        cursor, batch = client.scan(cursor=cursor, match=pattern, count=1000)
        keys.extend(batch)
        if cursor == 0:
            break

    if keys:
        deleted_count = client.delete(*keys)
        print(f"Deleted {deleted_count} {description} cache keys")
    else:
        print(f"No {description} cache keys found")

    return len(keys)

def main():
    redis_host = get_env("REDIS_HOST", "localhost")
    redis_port = get_int_env("REDIS_PORT", 6379)
    redis_password = get_env("REDIS_PASSWORD", "")
    redis_db = get_int_env("REDIS_DB", 0)

    try:
        client = redis.Redis(
            host=redis_host,
            port=redis_port,
            password=redis_password if redis_password else None,
            db=redis_db,
            decode_responses=True
        )

        # Test connection
        client.ping()
        print("Connected to Redis successfully")

        # Invalidate product list cache
        scan_and_delete_keys(client, "products:list:*", "product list")

        # Invalidate individual product cache
        scan_and_delete_keys(client, "product:*", "individual product")

        print("Cache invalidation completed successfully")

    except redis.ConnectionError as e:
        print(f"Failed to connect to Redis: {e}", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"An error occurred: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
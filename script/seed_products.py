#!/usr/bin/env python3
import argparse
import random
import sys
from typing import List, Tuple
import psycopg2
from psycopg2.extras import execute_batch

CATEGORIES = ['Electronics', 'Furniture', 'Appliances', 'Stationery', 'Accessories', 'Sports', 'Kitchen', 'Home']

PRODUCT_TEMPLATES = [
    ('Laptop Pro', 'High-performance laptop with {}GB RAM and {}GB SSD', (999.99, 2499.99), 'Electronics'),
    ('Wireless Mouse', 'Ergonomic wireless mouse with {} buttons', (19.99, 79.99), 'Electronics'),
    ('Office Chair', '{} ergonomic office chair with lumbar support', (199.00, 699.00), 'Furniture'),
    ('Coffee Maker', '{}-cup coffee maker with timer', (49.99, 299.99), 'Appliances'),
    ('Notebook Set', 'Set of {} ruled notebooks, 200 pages each', (9.99, 29.99), 'Stationery'),
    ('USB-C Cable', '{}-meter USB-C charging cable', (14.99, 39.99), 'Electronics'),
    ('Desk Lamp', '{} desk lamp with adjustable brightness', (29.99, 149.99), 'Furniture'),
    ('Water Bottle', 'Insulated water bottle, {}oz capacity', (19.99, 49.99), 'Accessories'),
    ('Bluetooth Speaker', 'Portable Bluetooth speaker with {}-hour battery', (39.99, 299.99), 'Electronics'),
    ('Standing Desk', '{} height-adjustable standing desk', (399.00, 1299.00), 'Furniture'),
    ('Keyboard', '{} mechanical keyboard with RGB lighting', (79.99, 249.99), 'Electronics'),
    ('Monitor', '{}-inch {} resolution monitor', (199.99, 899.99), 'Electronics'),
    ('Backpack', '{}-liter {} backpack with laptop compartment', (49.99, 149.99), 'Accessories'),
    ('Headphones', '{} headphones with noise cancellation', (99.99, 399.99), 'Electronics'),
    ('Yoga Mat', '{} yoga mat with carrying strap', (24.99, 79.99), 'Sports'),
    ('Blender', '{}-speed blender with {} cups capacity', (39.99, 199.99), 'Appliances'),
    ('Pen Set', 'Set of {} professional pens', (14.99, 49.99), 'Stationery'),
    ('Webcam', '{}P webcam with {} field of view', (59.99, 199.99), 'Electronics'),
    ('Bookshelf', '{}-tier {} bookshelf', (89.99, 299.99), 'Furniture'),
    ('Microwave', '{} watt microwave oven with {} presets', (79.99, 299.99), 'Appliances'),
]

ADJECTIVES = ['Premium', 'Professional', 'Deluxe', 'Standard', 'Basic', 'Advanced', 'Compact', 'Portable']
SIZES = ['Small', 'Medium', 'Large', 'Extra Large']
COLORS = ['Black', 'White', 'Silver', 'Blue', 'Red', 'Green']

def generate_product_data(count: int) -> List[Tuple]:
    products = []
    for i in range(count):
        template = random.choice(PRODUCT_TEMPLATES)
        name_base, desc_template, price_range, category = template
        
        variant = random.choice(ADJECTIVES) if random.random() > 0.5 else ''
        size_or_num = random.choice([8, 12, 16, 24, 32, 48, 64])
        
        if '{}' in desc_template:
            if desc_template.count('{}') == 1:
                description = desc_template.format(size_or_num)
            else:
                description = desc_template.format(size_or_num, size_or_num * 2)
        else:
            description = desc_template
        
        name = f"{variant} {name_base} {i+1}".strip() if variant else f"{name_base} {i+1}"
        price = round(random.uniform(price_range[0], price_range[1]), 2)
        stock = random.randint(5, 500)
        sku = f"SKU-{category[:3].upper()}-{i+1:06d}"
        is_active = random.random() > 0.1
        
        products.append((name, description, price, stock, category, sku, is_active))
    
    return products

def batch_insert_products(products: List[Tuple], batch_size: int = 1000):
    try:
        conn = psycopg2.connect(
            host="localhost",
            database="vave_db",
            user="postgres",
            password=""
        )
        
        cursor = conn.cursor()
        
        insert_query = """
            INSERT INTO products (name, description, price, stock_quantity, category, sku, is_active)
            VALUES (%s, %s, %s, %s, %s, %s, %s)
            ON CONFLICT (sku) DO NOTHING
        """
        
        total_batches = (len(products) + batch_size - 1) // batch_size
        inserted = 0
        
        for batch_num in range(total_batches):
            start_idx = batch_num * batch_size
            end_idx = min((batch_num + 1) * batch_size, len(products))
            batch = products[start_idx:end_idx]
            
            execute_batch(cursor, insert_query, batch, page_size=batch_size)
            inserted += len(batch)
            
            print(f"Progress: {inserted}/{len(products)} products inserted ({(inserted/len(products)*100):.1f}%)")
        
        conn.commit()
        
        cursor.execute("SELECT COUNT(*) FROM products")
        total_count = cursor.fetchone()[0]
        
        cursor.close()
        conn.close()
        
        return inserted, total_count
        
    except psycopg2.Error as e:
        print(f"Database error: {e}")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(
        description='Seed products into the database with batch insert',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python script/seed_products.py 100           # Insert 100 products
  python script/seed_products.py 10000 -b 500 # Insert 10k products with batch size 500
        """
    )
    parser.add_argument('count', type=int, help='Number of products to insert')
    parser.add_argument('-b', '--batch-size', type=int, default=1000,
                        help='Batch size for insert operations (default: 1000)')
    
    args = parser.parse_args()
    
    if args.count <= 0:
        print("Error: Count must be positive")
        sys.exit(1)
    
    if args.batch_size <= 0:
        print("Error: Batch size must be positive")
        sys.exit(1)
    
    print(f"Generating {args.count} products...")
    products = generate_product_data(args.count)
    
    print(f"Inserting products in batches of {args.batch_size}...")
    inserted, total = batch_insert_products(products, args.batch_size)
    
    print(f"\n✓ Successfully processed {inserted} products")
    print(f"✓ Total products in database: {total}")

if __name__ == '__main__':
    main()

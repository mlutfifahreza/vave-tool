#!/usr/bin/env python3
import argparse
import random
import sys
from typing import List, Tuple, Dict
import psycopg2
from psycopg2.extras import execute_batch

CATEGORIES_DATA = [
    {
        'id': 'electronics',
        'name': 'Electronics',
        'description': 'Electronic devices and gadgets'
    },
    {
        'id': 'furniture',
        'name': 'Furniture',
        'description': 'Home and office furniture'
    },
    {
        'id': 'appliances',
        'name': 'Appliances',
        'description': 'Household appliances'
    },
    {
        'id': 'stationery',
        'name': 'Stationery',
        'description': 'Office and school supplies'
    },
    {
        'id': 'accessories',
        'name': 'Accessories',
        'description': 'Fashion and personal accessories'
    },
    {
        'id': 'sports',
        'name': 'Sports',
        'description': 'Sports and fitness equipment'
    },
    {
        'id': 'kitchen',
        'name': 'Kitchen',
        'description': 'Kitchenware and cooking tools'
    },
    {
        'id': 'home',
        'name': 'Home',
        'description': 'Home improvement and decor'
    }
]

SUBCATEGORIES_DATA = [
    # Electronics
    {'id': 'electronics-laptops', 'category_name': 'Electronics', 'name': 'Laptops', 'description': 'Portable computers'},
    {'id': 'electronics-accessories', 'category_name': 'Electronics', 'name': 'Accessories', 'description': 'Computer accessories'},
    {'id': 'electronics-audio', 'category_name': 'Electronics', 'name': 'Audio', 'description': 'Audio equipment'},
    {'id': 'electronics-monitors', 'category_name': 'Electronics', 'name': 'Monitors', 'description': 'Display devices'},
    
    # Furniture
    {'id': 'furniture-chairs', 'category_name': 'Furniture', 'name': 'Chairs', 'description': 'Seating furniture'},
    {'id': 'furniture-desks', 'category_name': 'Furniture', 'name': 'Desks', 'description': 'Work surfaces'},
    {'id': 'furniture-storage', 'category_name': 'Furniture', 'name': 'Storage', 'description': 'Storage solutions'},
    
    # Appliances
    {'id': 'appliances-kitchen', 'category_name': 'Appliances', 'name': 'Kitchen', 'description': 'Kitchen appliances'},
    {'id': 'appliances-cleaning', 'category_name': 'Appliances', 'name': 'Cleaning', 'description': 'Cleaning devices'},
    
    # Stationery
    {'id': 'stationery-writing', 'category_name': 'Stationery', 'name': 'Writing', 'description': 'Writing instruments'},
    {'id': 'stationery-paper', 'category_name': 'Stationery', 'name': 'Paper', 'description': 'Paper products'},
    
    # Accessories
    {'id': 'accessories-bags', 'category_name': 'Accessories', 'name': 'Bags', 'description': 'Bags and luggage'},
    {'id': 'accessories-personal', 'category_name': 'Accessories', 'name': 'Personal', 'description': 'Personal items'},
    
    # Sports
    {'id': 'sports-fitness', 'category_name': 'Sports', 'name': 'Fitness', 'description': 'Fitness equipment'},
    {'id': 'sports-outdoor', 'category_name': 'Sports', 'name': 'Outdoor', 'description': 'Outdoor sports'},
    
    # Kitchen
    {'id': 'kitchen-cookware', 'category_name': 'Kitchen', 'name': 'Cookware', 'description': 'Cooking utensils'},
    {'id': 'kitchen-utensils', 'category_name': 'Kitchen', 'name': 'Utensils', 'description': 'Kitchen tools'},
    
    # Home
    {'id': 'home-decor', 'category_name': 'Home', 'name': 'Decor', 'description': 'Home decoration'},
    {'id': 'home-tools', 'category_name': 'Home', 'name': 'Tools', 'description': 'Home improvement tools'}
]

PRODUCT_TEMPLATES = [
    ('Laptop Pro', 'High-performance laptop with {}GB RAM and {}GB SSD', (999.99, 2499.99), 'Electronics', 'Laptops'),
    ('Wireless Mouse', 'Ergonomic wireless mouse with {} buttons', (19.99, 79.99), 'Electronics', 'Accessories'),
    ('Office Chair', '{} ergonomic office chair with lumbar support', (199.00, 699.00), 'Furniture', 'Chairs'),
    ('Coffee Maker', '{}-cup coffee maker with timer', (49.99, 299.99), 'Appliances', 'Kitchen'),
    ('Notebook Set', 'Set of {} ruled notebooks, 200 pages each', (9.99, 29.99), 'Stationery', 'Paper'),
    ('USB-C Cable', '{}-meter USB-C charging cable', (14.99, 39.99), 'Electronics', 'Accessories'),
    ('Desk Lamp', '{} desk lamp with adjustable brightness', (29.99, 149.99), 'Home', 'Decor'),
    ('Water Bottle', 'Insulated water bottle, {}oz capacity', (19.99, 49.99), 'Accessories', 'Personal'),
    ('Bluetooth Speaker', 'Portable Bluetooth speaker with {}-hour battery', (39.99, 299.99), 'Electronics', 'Audio'),
    ('Standing Desk', '{} height-adjustable standing desk', (399.00, 1299.00), 'Furniture', 'Desks'),
    ('Keyboard', '{} mechanical keyboard with RGB lighting', (79.99, 249.99), 'Electronics', 'Accessories'),
    ('Monitor', '{}-inch {} resolution monitor', (199.99, 899.99), 'Electronics', 'Monitors'),
    ('Backpack', '{}-liter {} backpack with laptop compartment', (49.99, 149.99), 'Accessories', 'Bags'),
    ('Headphones', '{} headphones with noise cancellation', (99.99, 399.99), 'Electronics', 'Audio'),
    ('Yoga Mat', '{} yoga mat with carrying strap', (24.99, 79.99), 'Sports', 'Fitness'),
    ('Blender', '{}-speed blender with {} cups capacity', (39.99, 199.99), 'Appliances', 'Kitchen'),
    ('Pen Set', 'Set of {} professional pens', (14.99, 49.99), 'Stationery', 'Writing'),
    ('Webcam', '{}P webcam with {} field of view', (59.99, 199.99), 'Electronics', 'Accessories'),
    ('Bookshelf', '{}-tier {} bookshelf', (89.99, 299.99), 'Furniture', 'Storage'),
    ('Microwave', '{} watt microwave oven with {} presets', (79.99, 299.99), 'Appliances', 'Kitchen'),
]

ADJECTIVES = ['Premium', 'Professional', 'Deluxe', 'Standard', 'Basic', 'Advanced', 'Compact', 'Portable']
SIZES = ['Small', 'Medium', 'Large', 'Extra Large']
COLORS = ['Black', 'White', 'Silver', 'Blue', 'Red', 'Green']

def get_category_map() -> Dict[str, str]:
    """Return mapping of category names to IDs"""
    return {cat['name']: cat['id'] for cat in CATEGORIES_DATA}

def get_subcategory_map() -> Dict[Tuple[str, str], str]:
    """Return mapping of (category_name, subcategory_name) to subcategory ID"""
    return {(sub['category_name'], sub['name']): sub['id'] for sub in SUBCATEGORIES_DATA}

def seed_categories(cursor) -> None:
    """Seed categories table"""
    insert_query = """
        INSERT INTO categories (id, name, description, is_active)
        VALUES (%s, %s, %s, true)
        ON CONFLICT (id) DO NOTHING
    """
    
    category_data = [(cat['id'], cat['name'], cat['description']) for cat in CATEGORIES_DATA]
    execute_batch(cursor, insert_query, category_data, page_size=100)
    print(f"Seeded {len(CATEGORIES_DATA)} categories")

def seed_subcategories(cursor) -> None:
    """Seed subcategories table"""
    category_map = get_category_map()
    
    insert_query = """
        INSERT INTO subcategories (id, category_id, name, description, is_active)
        VALUES (%s, %s, %s, %s, true)
        ON CONFLICT (id) DO NOTHING
    """
    
    subcategory_data = []
    for sub in SUBCATEGORIES_DATA:
        category_id = category_map.get(sub['category_name'])
        if category_id:
            subcategory_data.append((sub['id'], category_id, sub['name'], sub['description']))
    
    execute_batch(cursor, insert_query, subcategory_data, page_size=100)
    print(f"Seeded {len(subcategory_data)} subcategories")

def generate_product_data(count: int) -> List[Tuple]:
    category_map = get_category_map()
    subcategory_map = get_subcategory_map()
    
    products = []
    for i in range(count):
        template = random.choice(PRODUCT_TEMPLATES)
        name_base, desc_template, price_range, category_name, subcategory_name = template
        
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
        
        # Get category and subcategory IDs
        category_id = category_map.get(category_name)
        subcategory_id = subcategory_map.get((category_name, subcategory_name))
        
        sku = f"SKU-{category_name[:3].upper()}-{i+1:06d}"
        is_active = random.random() > 0.1
        
        products.append((name, description, price, stock, category_id, subcategory_id, sku, is_active))
    
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
        
        # Seed categories and subcategories first
        seed_categories(cursor)
        seed_subcategories(cursor)
        
        insert_query = """
            INSERT INTO products (name, description, price, stock_quantity, category_id, subcategory_id, sku, is_active)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
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

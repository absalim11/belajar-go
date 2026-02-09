-- Create Categories Table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Products Table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    harga INTEGER NOT NULL,
    stok INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

-- Create Index for better join performance
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);

-- Insert Sample Categories
INSERT INTO categories (name, description) VALUES
('Makanan', 'Kategori produk makanan'),
('Minuman', 'Kategori produk minuman'),
('Bumbu', 'Kategori produk bumbu dapur');

-- Insert Sample Products with Category
INSERT INTO products (nama, harga, stok, category_id) VALUES
('Indomie Godog', 3500, 10, 1),
('Vit 1000ml', 3000, 40, 2),
('Kecap', 12000, 20, 3);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create Transactions Table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    total_amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Transaction Details Table
CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL,
    subtotal INT NOT NULL
);

-- Create Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_transaction_details_transaction_id
    ON transaction_details(transaction_id);
CREATE INDEX IF NOT EXISTS idx_transaction_details_product_id
    ON transaction_details(product_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at
    ON transactions(created_at);
